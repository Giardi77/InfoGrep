package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Pattern struct {
	Name       string `yaml:"name"`
	Regex      string `yaml:"regex"`
	Confidence string `yaml:"confidence"`
}

type PatternFile struct {
	Patterns []Pattern `yaml:"patterns"`
}

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

func AddCustomPattern(name, path string) error {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir = filepath.Join(configDir, ".config")
	configFile := filepath.Join(configDir, "infogrep.patterns.json")

	// Ensure the config directory exists
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Load existing patterns or create an empty map if the file doesn't exist
	patterns := make(map[string]string)
	if _, err := os.Stat(configFile); err == nil {
		data, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &patterns)
		if err != nil {
			return err
		}
	}

	// Add or update the new pattern
	patterns[name] = path

	// Save the updated patterns back to the file
	data, err := json.MarshalIndent(patterns, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Pattern '%s' added successfully.\n", name)
	return nil
}

func GetPatterns(patternName string) (PatternFile, error) {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return PatternFile{}, err
	}
	configFile := filepath.Join(configDir, ".config", "infogrep.patterns.json")

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return PatternFile{}, err
	}

	var patterns map[string]string
	err = json.Unmarshal(data, &patterns)
	if err != nil {
		return PatternFile{}, err
	}

	patternFile, ok := patterns[patternName]
	if !ok {
		return PatternFile{}, fmt.Errorf("pattern '%s' not found", patternName)
	}

	data, err = ioutil.ReadFile(patternFile)
	if err != nil {
		return PatternFile{}, err
	}

	var patternData PatternFile
	err = yaml.Unmarshal(data, &patternData)
	if err != nil {
		return PatternFile{}, err
	}

	return patternData, nil
}

func GetAllAbsPaths(fileOrDirectory string) ([]string, error) {
	var filePaths []string

	fileInfo, err := os.Stat(fileOrDirectory)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		err := filepath.Walk(fileOrDirectory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				filePaths = append(filePaths, absPath)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		absPath, err := filepath.Abs(fileOrDirectory)
		if err != nil {
			return nil, err
		}
		filePaths = append(filePaths, absPath)
	}

	return filePaths, nil
}

func PrintResult(pattern Pattern, result string, filePath string, position int) {
	var confidenceColor string
	switch pattern.Confidence {
	case "high":
		confidenceColor = "\033[31m" // Red
	case "medium":
		confidenceColor = "\033[33m" // Yellow
	case "low":
		confidenceColor = "\033[34m" // Blue
	default:
		confidenceColor = "\033[0m" // Reset
	}

	fmt.Printf("\n[%s] [%s%s\033[0m]\n\n%s\n", pattern.Name, confidenceColor, pattern.Confidence, result)
	fmt.Printf("File: %s\n", filePath)
	fmt.Printf("Position: %d\n", position)
}
