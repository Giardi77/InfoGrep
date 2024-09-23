use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::fs;
use std::path::PathBuf;

#[derive(Debug, Serialize, Deserialize)]
pub struct PatternConfig {
    secrets: String,
    pii: String,
    gitleaks: String,
}

pub fn create_default_config() -> Result<()> {
    let config_dir = dirs::home_dir()
        .context("Could not find home directory")?
        .join(".config/infogrep");
    let config_file = config_dir.join("infogrep.patterns.json");

    if !config_file.exists() {
        let default_config = PatternConfig {
            secrets: "default-patterns/rules-stable.yml".to_string(),
            pii: "default-patterns/pii-stable.yml".to_string(),
            gitleaks: "default-patterns/gitleaks.yml".to_string(),
        };

        fs::create_dir_all(&config_dir)
            .with_context(|| format!("Failed to create config directory: {:?}", config_dir))?;
        let config_content = serde_json::to_string_pretty(&default_config)
            .context("Failed to serialize default config")?;
        fs::write(&config_file, config_content)
            .with_context(|| format!("Failed to write config file: {:?}", config_file))?;
        println!("Default config file created at {:?}", config_file);
    }

    let default_patterns = vec!["rules-stable.yml", "pii-stable.yml", "gitleaks.yml"];
    let patterns_dir = config_dir.join("default-patterns");

    // Ensure the default-patterns directory exists
    fs::create_dir_all(&patterns_dir)
        .with_context(|| format!("Failed to create patterns directory: {:?}", patterns_dir))?;

    for pattern in default_patterns {
        let pattern_path = patterns_dir.join(pattern);
        if !pattern_path.exists() {
            println!("Downloading missing default pattern {} ...", pattern);
            let url = format!("https://raw.githubusercontent.com/Giardi77/InfoGrep/refs/heads/Version-3-Rust/default-patterns/{}", pattern);
            let response = reqwest::blocking::get(&url)
                .with_context(|| format!("Failed to fetch URL: {}", url))?;
            let content = response
                .text()
                .with_context(|| format!("Failed to read response text from URL: {}", url))?;
            fs::write(&pattern_path, content)
                .with_context(|| format!("Failed to write pattern file: {:?}", pattern_path))?;
            println!("Downloaded and saved pattern file: {}", pattern);
        }
    }

    Ok(())
}

pub fn get_pattern_file(pattern: &str) -> Result<String> {
    let config_dir = dirs::home_dir()
        .context("Could not find home directory")?
        .join(".config/infogrep");
    let config_file = config_dir.join("infogrep.patterns.json");

    let json_content = fs::read_to_string(&config_file)
        .with_context(|| format!("Failed to read config file: {:?}", config_file))?;
    let json: Value =
        serde_json::from_str(&json_content).context("Failed to parse config file as JSON")?;

    let relative_path = json
        .get(pattern)
        .and_then(Value::as_str)
        .map(String::from)
        .ok_or_else(|| anyhow::anyhow!("Pattern '{}' not found in config file", pattern))?;

    let full_path = config_dir.join(relative_path);
    Ok(full_path.to_string_lossy().to_string())
}

#[derive(Debug, Deserialize, Serialize)]
pub struct PatternDetails {
    pub name: String,
    pub regex: String,
    pub confidence: String,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Pattern {
    pub pattern: PatternDetails,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Patterns {
    pub patterns: Vec<Pattern>,
}

pub fn load_patterns(file_path: &str) -> Result<Patterns> {
    let yaml_content = fs::read_to_string(file_path)
        .with_context(|| format!("Failed to read pattern file: {}", file_path))?;
    let patterns: Patterns = serde_yaml::from_str(&yaml_content)
        .with_context(|| format!("Failed to parse pattern file as YAML: {}", file_path))?;
    Ok(patterns)
}

pub fn get_files_to_scan(input: &str) -> Result<Vec<PathBuf>> {
    let path = fs::canonicalize(input)
        .with_context(|| format!("Failed to canonicalize input path: {}", input))?;
    if path.is_file() {
        Ok(vec![path])
    } else if path.is_dir() {
        let mut files = Vec::new();
        for entry in
            fs::read_dir(path).with_context(|| format!("Failed to read directory: {}", input))?
        {
            let entry = entry?;
            let path = entry.path();
            if path.is_file() {
                files.push(path);
            }
        }
        Ok(files)
    } else {
        Err(anyhow::anyhow!(
            "'{}' is neither a file nor a directory",
            input
        ))
    }
}

pub fn truncate_string(s: &str, max_chars: usize) -> String {
    if s.chars().count() <= max_chars {
        s.to_string()
    } else {
        s.chars().take(max_chars).collect::<String>() + "..."
    }
}

pub fn print_logo() {
    println!(
        r#"
  _____        __         ___                
  \_   \_ __  / _| ___   / _ \_ __ ___ _ __  
   / /\/ '_ \| |_ / _ \ / /_\/ '__/ _ \ '_ \ 
/\/ /_ | | | |  _| (_) / /_\\| | |  __/ |_) |
\____/ |_| |_|_|  \___/\____/|_|  \___| .__/ 
                                      |_|    
"#
    );
}
