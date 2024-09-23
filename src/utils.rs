use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::error::Error;
use std::fs;
use std::path::PathBuf;

pub fn get_pattern_file(pattern: &str) -> Result<String, Box<dyn std::error::Error>> {
    let config_dir = dirs::home_dir()
        .ok_or("Could not find home directory")?
        .join(".config")
        .join("infogrep.patterns.json");

    let json_content = fs::read_to_string(config_dir)?;
    let json: Value = serde_json::from_str(&json_content)?;

    json.get(pattern)
        .and_then(Value::as_str)
        .map(String::from)
        .ok_or_else(|| format!("Pattern '{}' not found in config file", pattern).into())
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

pub fn load_patterns(file_path: &str) -> Result<Patterns, Box<dyn Error>> {
    let yaml_content = fs::read_to_string(file_path)?;
    let patterns: Patterns = serde_yaml::from_str(&yaml_content)?;
    Ok(patterns)
}

pub fn get_files_to_scan(input: &str) -> Result<Vec<PathBuf>, Box<dyn Error>> {
    let path = fs::canonicalize(input)?;
    if path.is_file() {
        Ok(vec![path])
    } else if path.is_dir() {
        let mut files = Vec::new();
        for entry in fs::read_dir(path)? {
            let entry = entry?;
            let path = entry.path();
            if path.is_file() {
                files.push(path);
            }
        }
        Ok(files)
    } else {
        Err(format!("'{}' is neither a file nor a directory", input).into())
    }
}

pub fn truncate_string(s: &str, max_chars: usize) -> String {
    if s.chars().count() <= max_chars {
        s.to_string()
    } else {
        s.chars().take(max_chars).collect::<String>() + "..."
    }
}
