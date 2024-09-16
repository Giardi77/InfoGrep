import argparse
import re
import utils

parser = argparse.ArgumentParser(description="Grep for sensitive info")
parser.add_argument('-f','--file', help="file or files (comma-separated)")
parser.add_argument('-d','--directory', help="files in the directory")
#parser.add_argument('-l','--link', help="link, grep the resulting http request")
parser.add_argument('-p','--pattern', default='secrets', help="pick a pattern from .config/infogrep.json")
parser.add_argument('-a','--add-pattern', help="add a pattern file to .config/infogrep.json (provide name,path [-a pid,/usr/share/wordlists/pattern.yml])")

args = parser.parse_args()

def greppin(path: str, pattern: dict) -> str | None :
    try:
        with open(path, 'r') as file:
            file_content = file.read()
        
        matches = re.match(pattern['pattern']['regex'], file_content, re.MULTILINE)
        
        if matches:
            return matches.string

    except FileNotFoundError as e:
        print(f"FileNotFoundError: {e}")
    except Exception as e:
        print(f"An error occurred: {e}")

files2Scan = []

if args.file:
    files2Scan += utils.get_all_abs_paths(args.file)

if args.directory:
    files2Scan += utils.get_all_abs_paths(args.directory)

PatternName = args.pattern

def main():
    print(utils.logo)
    Patterns = utils.getPatterns(PatternName)
    for path in files2Scan:
        print(f"Scanning [ {path} ]")
        for pattern in Patterns['patterns']:
            try:
                res = greppin(path, pattern)
                if res:
                    utils.print_result(pattern, res)
            except KeyboardInterrupt:
                print('\rAborting ...')
                exit(0)


if __name__ == "__main__":
    main()