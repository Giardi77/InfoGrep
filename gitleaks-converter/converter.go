package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Rule represents a single rule in the TOML file
type Rule struct {
	ID          string   `toml:"id"`
	Description string   `toml:"description"`
	Regex       string   `toml:"regex"`
	Keywords    []string `toml:"keywords"`
}

// GitleaksConfig represents the structure of the gitleaks TOML file
type GitleaksConfig struct {
	Title     string `toml:"title"`
	Allowlist struct {
		Description string   `toml:"description"`
		Paths       []string `toml:"paths"`
	} `toml:"allowlist"`
	Rules []Rule `toml:"rules"`
}

// PatternDetails represents the pattern details in the YAML file
type PatternDetails struct {
	Name       string `yaml:"name"`
	Regex      string `yaml:"regex"`
	Confidence string `yaml:"confidence"`
}

// Pattern represents a single pattern in the YAML file
type Pattern struct {
	Pattern PatternDetails `yaml:"pattern"`
}

// Patterns represents the structure of the YAML file
type Patterns struct {
	Patterns []Pattern `yaml:"patterns"`
}

func main() {
	tomlFile := "gitleaks.toml"
	yamlFile := "default-patterns/gitleaks.yml"

	// Read the TOML file
	content, err := ioutil.ReadFile(tomlFile)
	if err != nil {
		log.Fatalf("Error reading TOML file: %v", err)
	}

	// Decode the TOML content
	var config GitleaksConfig
	if _, err := toml.Decode(string(content), &config); err != nil {
		log.Fatalf("Error decoding TOML file: %v", err)
	}

	// Prepare the YAML structure
	var patterns Patterns
	for _, rule := range config.Rules {
		pattern := Pattern{
			Pattern: PatternDetails{
				Name:       rule.ID,
				Regex:      rule.Regex,
				Confidence: "high", // Default confidence level
			},
		}
		patterns.Patterns = append(patterns.Patterns, pattern)
	}

	// Encode the YAML content
	yamlContent, err := yaml.Marshal(&patterns)
	if err != nil {
		log.Fatalf("Error encoding YAML file: %v", err)
	}

	// Write the YAML file
	if err := ioutil.WriteFile(yamlFile, yamlContent, 0644); err != nil {
		log.Fatalf("Error writing YAML file: %v", err)
	}

	fmt.Println("Transformation complete. YAML file created:", yamlFile)
}
