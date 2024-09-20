package utils

import (
	"encoding/json"
	"fmt"
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

func GetPatterns(patternType string) ([]Pattern, error) {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".config", "infogrep.patterns.json")

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config map[string]string
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	patternFile, ok := config[patternType]
	if !ok {
		return nil, fmt.Errorf("pattern type '%s' not found in config", patternType)
	}

	return readPatternsFile(patternFile)
}

func readPatternsFile(filename string) ([]Pattern, error) {
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

func ReadInputFlag(inputFlag string) ([]string, error) {

	absPath, err := filepath.Abs(inputFlag)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %v", err)
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("error accessing input: %v", err)
	}

	if fileInfo.IsDir() {
		return walkDirectory(absPath)
	}

	return []string{absPath}, nil
}

func walkDirectory(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("error getting absolute path for %s: %v", path, err)
			}
			files = append(files, absPath)
		}
		return nil
	})
	return files, err
}
