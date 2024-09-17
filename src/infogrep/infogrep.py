import argparse
import re
import utils
import sys
import json
import os

from . import utils

parser = argparse.ArgumentParser(description="Grep for sensitive info")
parser.add_argument('-i', '--input', help="file or directory to scan")
parser.add_argument('-p', '--pattern', default='secrets', help="pick a pattern from .config/infogrep.json")
parser.add_argument('-a', '--add-pattern', help="add a pattern file to .config/infogrep.json (provide name:path [-a name:/path/to/pattern.yml])")

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

def greppin(content: str, pattern: dict) -> str | None:
    matches = re.search(pattern['pattern']['regex'], content, re.MULTILINE)
    if matches:
        return matches.group()
    return None

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
                    res = greppin(content, pattern)
                    if res:
                        print()  # Move to the next line before printing results
                        utils.print_result(pattern, res)
            except Exception as e:
                print(f"\nAn error occurred while processing {path}: {e}")
    else:
        # Read from stdin if no input file/directory is specified
        content = sys.stdin.read()
        for pattern in Patterns['patterns']:
            res = greppin(content, pattern)
            if res:
                utils.print_result(pattern, res, "stdin")

    print()  # Print a newline at the end to move the cursor to the next line

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print('\rAborting ...')
        sys.exit(0)