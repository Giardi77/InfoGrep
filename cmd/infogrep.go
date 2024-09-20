package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/Giardi77/infogrep/pkg/utils"
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

func greppin(content string, pattern utils.Pattern) [][]int {
	re := regexp.MustCompile(pattern.Regex)
	return re.FindAllStringIndex(content, -1)
}

func truncateMatch(match string, maxChars int) string {
	if maxChars == 0 || len(match) <= maxChars {
		return match
	}
	return fmt.Sprintf("%s... (truncated, %d more characters)", match[:maxChars], len(match)-maxChars)
}

func main() {
	flag.Parse()

	if addPatternFlag != "" {
		parts := strings.SplitN(addPatternFlag, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Error: Invalid format for add-pattern. Use: name:/path/to/pattern.yml")
			return
		}
		err := utils.AddCustomPattern(parts[0], parts[1])
		if err != nil {
			fmt.Printf("Error adding custom pattern: %v\n", err)
		}
		return
	}

	fmt.Println(utils.Logo)

	patterns, err := utils.GetPatterns(patternFlag)
	if err != nil {
		fmt.Printf("Error loading patterns: %v\n", err)
		return
	}

	if inputFlag != "" {
		files, err := utils.GetAllAbsPaths(inputFlag)
		if err != nil {
			fmt.Printf("Error getting file paths: %v\n", err)
			return
		}
		for _, path := range files {
			fmt.Printf("\rScanning [ %s ]", path)
			content, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("\nError reading file %s: %v\n", path, err)
				continue
			}
			for _, pattern := range patterns.Patterns {
				results := greppin(string(content), pattern)
				if len(results) > 0 {
					fmt.Println() // Move to the next line before printing results
					for _, match := range results {
						result := string(content[match[0]:match[1]])
						truncatedResult := truncateMatch(result, truncateFlag)
						utils.PrintResult(pattern, truncatedResult, path, match[0])
					}
				}
			}
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		content, _ := ioutil.ReadAll(reader)
		for _, pattern := range patterns.Patterns {
			results := greppin(string(content), pattern)
			for _, match := range results {
				result := string(content[match[0]:match[1]])
				truncatedResult := truncateMatch(result, truncateFlag)
				utils.PrintResult(pattern, truncatedResult, "stdin", match[0])
			}
		}
	}

	fmt.Println() // Print a newline at the end
}
