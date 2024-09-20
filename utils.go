package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var Logo = `
 █████               ██████             █████████                              
░░███               ███░░███           ███░░░░░███                             
 ░███  ████████    ░███ ░░░   ██████  ███     ░░░  ████████   ██████  ████████ 
 ░███ ░░███░░███  ███████    ███░░███░███         ░░███░░███ ███░░███░░███░░███
 ░███  ░███ ░███ ░░░███░    ░███ ░███░███    █████ ░███ ░░░ ░███████  ░███ ░███
 ░███  ░███ ░███   ░███     ░███ ░███░░███  ░░███  ░███     ░███░░░   ░███ ░███
 █████ ████ █████  █████    ░░██████  ░░█████████  █████    ░░██████  ░███████ 
░░░░░ ░░░░ ░░░░░  ░░░░░      ░░░░░░    ░░░░░░░░░  ░░░░░      ░░░░░░   ░███░░░  
                                                                      ░███     
                                                                      █████    
                                                                     ░░░░░     
`

type Pattern struct {
	Name       string `yaml:"name"`
	Regex      string `yaml:"regex"`
	Confidence string `yaml:"confidence"`
}

type Config struct {
	Patterns []struct {
		Pattern Pattern `yaml:"pattern"`
	} `yaml:"patterns"`
}

func GetPatterns() ([]Pattern, error) {
	patternsFile := filepath.Join("default-patterns", "patterns-stable.yml")
	return readFile(patternsFile)
}

func readFile(filename string) ([]Pattern, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	patterns := make([]Pattern, len(config.Patterns))
	for i, p := range config.Patterns {
		patterns[i] = p.Pattern
	}

	return patterns, nil
}

func ReadInput(inputFlag string) (string, error) {
	if inputFlag != "" {
		content, err := os.ReadFile(inputFlag)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		return "", fmt.Errorf("no input provided")
	}

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	return string(output), nil
}
