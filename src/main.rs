use clap::Parser;
mod utils;
use regex::Regex;
use std::collections::HashMap;
use std::fs::File;
use std::io::{Read, Seek, SeekFrom};
use std::path::Path;
use utils::{get_files_to_scan, get_pattern_file, load_patterns, truncate_string, Pattern};

#[derive(Parser, Debug)]
#[command(name = "InfoGrep", about = "Grep for sensitive info", long_about = None)]
struct Args {
    /// Input file or directory
    #[arg(short, long, value_name = "INPUT")]
    input: String,

    /// Pattern to use
    #[arg(short, long, value_name = "PATTERN", default_value = "secrets")]
    pattern: String,

    /// Truncate output to this many characters
    #[arg(short, long, value_name = "TRUNCATE", default_value = "400")]
    truncate: usize,
}

struct CompiledPattern {
    regex: Regex,
    name: String,
    confidence: String,
}

fn compile_patterns(patterns: &[Pattern]) -> Result<Vec<CompiledPattern>, regex::Error> {
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

fn scan_file(
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
                println!(
                    "Match found in {} (around line {})",
                    file_path.display(),
                    approx_line
                );
                println!("Pattern: {} ({})", pattern.name, pattern.confidence);
                println!("Matched: {}", truncated_match);
                println!("---");
            }
        }

        // Move file pointer back by the length of leftover to ensure we don't miss anything
        file.seek(SeekFrom::Current(-(leftover.len() as i64)))?;
    }

    Ok(())
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Args::parse();
    let pattern_file = get_pattern_file(&args.pattern)?;
    let patterns = load_patterns(&pattern_file)?;
    let compiled_patterns = compile_patterns(&patterns.patterns)?;
    println!("Compiled {} patterns", compiled_patterns.len());

    let files_to_scan = get_files_to_scan(&args.input)?;

    for file in files_to_scan {
        scan_file(&file, &compiled_patterns, args.truncate)?;
    }

    Ok(())
}
