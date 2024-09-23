use crate::utils::truncate_string;
use colored::*;
use regex::Regex;
use std::fs::File;
use std::io::{Read, Seek, SeekFrom};
use std::path::Path;

pub struct CompiledPattern {
    pub regex: Regex,
    pub name: String,
    pub confidence: String,
}

pub fn compile_patterns(
    patterns: &[crate::utils::Pattern],
) -> Result<Vec<CompiledPattern>, regex::Error> {
    patterns
        .iter()
        .map(|p| {
            Ok(CompiledPattern {
                regex: Regex::new(&p.pattern.regex)?,
                name: p.pattern.name.clone(),
                confidence: p.pattern.confidence.clone(),
            })
        })
        .collect()
}

pub fn scan_file(
    file_path: &Path,
    compiled_patterns: &[CompiledPattern],
    truncate: usize,
) -> Result<(), Box<dyn std::error::Error>> {
    let mut file = File::open(file_path)?;
    let file_size = file.metadata()?.len();
    let chunk_size: u64 = 1024 * 1024; // 1MB chunks
    let mut buffer = vec![0; chunk_size as usize];
    let mut leftover = String::new();

    for chunk_start in (0..file_size).step_by(chunk_size as usize) {
        let bytes_read = file.read(&mut buffer)?;
        if bytes_read == 0 {
            break;
        }

        let mut chunk = leftover.clone() + &String::from_utf8_lossy(&buffer[..bytes_read]);

        // If this isn't the last chunk, find a suitable break point
        if chunk_start + chunk_size < file_size {
            if let Some(last_newline) = chunk.rfind('\n') {
                leftover = chunk.split_off(last_newline + 1);
            } else {
                leftover = chunk.split_off(chunk.len() - 100); // Arbitrary split if no newline
            }
        } else {
            leftover.clear();
        }

        for pattern in compiled_patterns {
            for mat in pattern.regex.find_iter(&chunk) {
                let matched = &chunk[mat.start()..mat.end()];
                let truncated_match = truncate_string(matched, truncate);
                let approx_line = (chunk_start / 80) + (mat.start() as u64 / 80) + 1; // Rough line number estimate

                // Determine the color based on the confidence level
                let confidence_colored = match pattern.confidence.as_str() {
                    "low" => pattern.confidence.blue(),
                    "medium" => pattern.confidence.yellow(),
                    "high" => pattern.confidence.red(),
                    _ => pattern.confidence.normal(),
                };

                println!("[{}] (position: {})", file_path.display(), approx_line);
                println!("[{}] [{}]", pattern.name, confidence_colored);
                println!("\n{}\n", truncated_match);
            }
        }

        // Move file pointer back by the length of leftover to ensure we don't miss anything
        file.seek(SeekFrom::Current(-(leftover.len() as i64)))?;
    }

    Ok(())
}
