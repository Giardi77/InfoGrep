package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/Giardi77/infogrep/utils"
	"io"
	"os"
	"regexp"
	"runtime"
	"sync"
)

var (
	inputFlag      string
	patternFlag    string
	addPatternFlag string
	truncateFlag   int
)

func init() {
	flag.StringVar(&inputFlag, "i", "", "file or directory to scan")
	flag.StringVar(&inputFlag, "input", "", "file or directory to scan")
	flag.StringVar(&patternFlag, "p", "secrets", "pick a pattern from .config/infogrep.patterns.json")
	flag.StringVar(&patternFlag, "pattern", "secrets", "pick a pattern from .config/infogrep.patterns.json")
	flag.StringVar(&addPatternFlag, "a", "", "add a pattern file to .config/infogrep.patterns.json (provide name:path)")
	flag.StringVar(&addPatternFlag, "add-pattern", "", "add a pattern file to .config/infogrep.patterns.json (provide name:path)")
	flag.IntVar(&truncateFlag, "t", 400, "Truncation length for output (0 for no truncation)")
	flag.IntVar(&truncateFlag, "truncate", 400, "Truncation length for output (0 for no truncation)")
}

func main() {
	flag.Parse()

	patterns, err := utils.GetPatterns(patternFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading patterns: %v\n", err)
		os.Exit(1)
	}

	filesToGrep, err := utils.ReadInputFlag(inputFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	numWorkers := runtime.NumCPU()
	distributeAndProcess(filesToGrep, patterns, truncateFlag, numWorkers)
}

func distributeAndProcess(files []string, patterns []utils.Pattern, truncateFlag, numWorkers int) {
	var wg sync.WaitGroup
	filesChan := make(chan string, len(files))
	resultsChan := make(chan string, 1000)

	// Precompile regex patterns
	compiledPatterns := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		compiledPatterns[i] = regexp.MustCompile(pattern.Regex)
	}

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(filesChan, resultsChan, patterns, compiledPatterns, truncateFlag, &wg)
	}

	// Start result printer goroutine
	go printResults(resultsChan)

	// Distribute files to workers
	for _, file := range files {
		filesChan <- file
	}
	close(filesChan)

	wg.Wait()
	close(resultsChan)
}

func worker(files <-chan string, results chan<- string, patterns []utils.Pattern, compiledPatterns []*regexp.Regexp, truncateFlag int, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range files {
		Greppin(file, results, patterns, compiledPatterns, truncateFlag)
	}
}

func Greppin(filePath string, results chan<- string, patterns []utils.Pattern, compiledPatterns []*regexp.Regexp, truncateFlag int) {
	file, err := os.Open(filePath)
	if err != nil {
		results <- fmt.Sprintf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 4*1024*1024) // 4MB buffer
	buffer := make([]byte, 4*1024*1024)
	lineNum := 1
	var partialLine []byte

	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			results <- fmt.Sprintf("Error reading file %s: %v\n", filePath, err)
			break
		}

		chunk := buffer[:n]
		lines := bytes.Split(append(partialLine, chunk...), []byte("\n"))
		partialLine = nil

		if err != io.EOF {
			partialLine = lines[len(lines)-1]
			lines = lines[:len(lines)-1]
		}

		for _, line := range lines {
			for i, re := range compiledPatterns {
				matches := re.FindAll(line, -1)
				for _, match := range matches {
					output := string(match)
					if truncateFlag > 0 && len(output) > truncateFlag {
						output = output[:truncateFlag] + "..."
					}
					results <- fmt.Sprintf("%s:%d: Found %s (Confidence: %s): %s\n", filePath, lineNum, patterns[i].Name, patterns[i].Confidence, output)
				}
			}
			lineNum++
		}

		if err == io.EOF {
			break
		}
	}
}

func printResults(results <-chan string) {
	for result := range results {
		fmt.Print(result)
	}
}
