use clap::Parser;
mod scanner;
mod utils;
use anyhow::Result;
use rayon::prelude::*; // Import ParallelIterator trait
use std::time::Instant;
use utils::{
    create_default_config, get_files_to_scan, get_pattern_file, load_patterns, print_logo,
};

#[derive(Parser, Debug)]
#[command(name = "InfoGrep", about = "Grep for sensitive info", long_about = None, version = env!("CARGO_PKG_VERSION"))]
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

    /// Number of worker threads to use
    #[arg(short, long, value_name = "WORKERS", default_value = "2")]
    workers: usize,

    /// Confidence level to use (low, medium, high)
    #[arg(short, long, value_name = "CONFIDENCE")]
    confidence: Option<String>,
}

fn main() -> Result<()> {
    let args = Args::parse();

    print_logo();

    // Create default config file if it doesn't exist
    create_default_config()?;

    // Start the timer
    let start = Instant::now();

    let pattern_file = get_pattern_file(&args.pattern)?;
    let patterns = load_patterns(&pattern_file)?;

    // Filter patterns based on confidence level if provided
    let filtered_patterns: Vec<_> = if let Some(confidence) = &args.confidence {
        patterns
            .patterns
            .into_iter()
            .filter(|p| p.pattern.confidence.eq_ignore_ascii_case(confidence))
            .collect()
    } else {
        patterns.patterns
    };

    if filtered_patterns.is_empty() {
        eprintln!(
            "No patterns found with confidence level: {:?}",
            args.confidence
        );
        return Ok(());
    }

    let compiled_patterns = scanner::compile_patterns(&filtered_patterns)?;
    println!("Compiled {} patterns", compiled_patterns.len());

    let files_to_scan = get_files_to_scan(&args.input)?;

    // Set the number of worker threads
    rayon::ThreadPoolBuilder::new()
        .num_threads(args.workers)
        .build_global()?;

    files_to_scan.par_iter().for_each(|file| {
        if let Err(e) = scanner::scan_file(file, &compiled_patterns, args.truncate) {
            eprintln!("Error scanning file {}: {}", file.display(), e);
        }
    });

    // Stop the timer
    let duration = start.elapsed();
    println!("Time taken: {:?}", duration);

    Ok(())
}
