use clap::Parser;
mod scanner;
mod utils;
use rayon::prelude::*; // Import ParallelIterator trait
use std::time::Instant;
use utils::{create_default_config, get_files_to_scan, get_pattern_file, load_patterns};

#[derive(Parser, Debug)]
#[command(name = "InfoGrep", about = "Grep for sensitive info", long_about = None)]
struct Args {
    /// Input file or directory (Required)
    #[arg(short, long, value_name = "INPUT")]
    input: String,

    /// Pattern to use (Default: secrets)
    #[arg(short, long, value_name = "PATTERN", default_value = "secrets")]
    pattern: String,

    /// Truncate output to this many characters (Default: 400)
    #[arg(short, long, value_name = "TRUNCATE", default_value = "400")]
    truncate: usize,

    /// Number of worker threads to use (Deafult: 2)
    #[arg(short, long, value_name = "WORKERS", default_value = "2")]
    workers: usize,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Args::parse();

    // Create default config file if it doesn't exist
    create_default_config()?;

    let start = Instant::now();

    let pattern_file = get_pattern_file(&args.pattern)?;
    let patterns = load_patterns(&pattern_file)?;
    let compiled_patterns = scanner::compile_patterns(&patterns.patterns)?;
    println!("Compiled {} patterns", compiled_patterns.len());

    let files_to_scan = get_files_to_scan(&args.input)?;

    rayon::ThreadPoolBuilder::new()
        .num_threads(args.workers)
        .build_global()?;

    files_to_scan.par_iter().for_each(|file| {
        if let Err(e) = scanner::scan_file(file, &compiled_patterns, args.truncate) {
            eprintln!("Error scanning file {}: {}", file.display(), e);
        }
    });

    let duration = start.elapsed();
    println!("Time taken: {:?}", duration);

    Ok(())
}
