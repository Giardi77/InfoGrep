import argparse
import re
import sys
import json
import os

from . import utils

parser = argparse.ArgumentParser(description="Grep for sensitive info")
parser.add_argument('-i', '--input', help="file or directory to scan")
parser.add_argument('-p', '--pattern', default='secrets', help="pick a pattern from .config/infogrep.patterns.json")
parser.add_argument('-a', '--add-pattern', help="add a pattern file to .config/infogrep.patterns.json (provide name:path [-a name:/path/to/pattern.yml])")

args = parser.parse_args()

def add_custom_pattern(name, path):
    config_dir = os.path.expanduser("~/.config")
    config_file = os.path.join(config_dir, "infogrep.patterns.json")

    # Ensure the config directory exists
    os.makedirs(config_dir, exist_ok=True)

    # Load existing patterns or create an empty dict if the file doesn't exist
    if os.path.exists(config_file):
        with open(config_file, 'r') as f:
            patterns = json.load(f)
    else:
        patterns = {}

    # Add or update the new pattern
    patterns[name] = path

    # Save the updated patterns back to the file
    with open(config_file, 'w') as f:
        json.dump(patterns, f, indent=2)

    print(f"Pattern '{name}' added successfully.")

def greppin(content: str, pattern: dict) -> list:
    matches = re.finditer(pattern['pattern']['regex'], content, re.MULTILINE)
    return [(match.group(), match.start()) for match in matches]  # Return match and its position

def truncate_match(match: str, max_chars: int = 300) -> str:
    if len(match) > max_chars:
        truncated = match[:max_chars]
        return f"{truncated}... (truncated, {len(match) - max_chars} more characters)"
    return match

PatternName = args.pattern

def main():
    if args.add_pattern:
        try:
            name, path = args.add_pattern.split(':')
            add_custom_pattern(name.strip(), path.strip())
            return
        except ValueError:
            print("Error: Invalid format for add-pattern. Use: name:/path/to/pattern.yml")
            return

    print(utils.logo)
    Patterns = utils.getPatterns(PatternName)

    if args.input:
        files2Scan = utils.get_all_abs_paths(args.input)
        for path in files2Scan:
            print(f"\rScanning [ {path} ]", end="", flush=True)
            try:
                with open(path, 'r') as file:
                    content = file.read()
                for pattern in Patterns['patterns']:
                    results = greppin(content, pattern)
                    if results:
                        print()  # Move to the next line before printing results
                        for res, pos in results:
                            truncated_res = truncate_match(res)
                            utils.print_result(pattern, truncated_res, path, pos)
            except Exception as e:
                print(f"\nAn error occurred while processing {path}: {e}")
    else:
        # Read from stdin if no input file/directory is specified
        content = sys.stdin.read()
        for pattern in Patterns['patterns']:
            results = greppin(content, pattern)
            for res, pos in results:
                truncated_res = truncate_match(res)
                utils.print_result(pattern, truncated_res, "stdin", pos)

    print()  # Print a newline at the end to move the cursor to the next line

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print('\rAborting ...')
        sys.exit(0)