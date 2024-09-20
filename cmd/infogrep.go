package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/Giardi77/infogrep/utils"
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
	patterns, err := utils.GetPatterns()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading patterns: %v\n", err)
		os.Exit(1)
	}

	input, err := utils.ReadInput(inputFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern.Regex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error compiling regex for %s: %v\n", pattern.Name, err)
			continue
		}

		matches := re.FindAllString(input, -1)
		for _, match := range matches {
			output := match
			if truncateFlag > 0 && len(output) > truncateFlag {
				output = output[:truncateFlag] + "..."
			}
			fmt.Printf("Found %s (Confidence: %s): %s\n", pattern.Name, pattern.Confidence, output)
		}
	}
}
