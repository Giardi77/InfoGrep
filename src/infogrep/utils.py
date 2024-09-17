from typing import List
import os
import yaml
import json
import pkg_resources

# Import the version from setup.py
try:
    from setuptools_scm import get_version
    VERSION = get_version(root='..', relative_to=__file__)
except ImportError:
    try:
        from infogrep._version import version as VERSION
    except ImportError:
        VERSION = "unknown"

logo = rf'''
 █████               ██████             █████████                              
░░███               ███░░███           ███░░░░░███                             
 ░███  ████████    ░███ ░░░   ██████  ███     ░░░  ████████   ██████  ████████ 
 ░███ ░░███░░███  ███████    ███░░███░███         ░░███░░███ ███░░███░░███░░███
 ░███  ░███ ░███ ░░░███░    ░███ ░███░███    █████ ░███ ░░░ ░███████  ░███ ░███
 ░███  ░███ ░███   ░███     ░███ ░███░░███  ░░███  ░███     ░███░░░   ░███ ░███
 █████ ████ █████  █████    ░░██████  ░░█████████  █████    ░░██████  ░███████ 
░░░░░ ░░░░ ░░░░░  ░░░░░      ░░░░░░    ░░░░░░░░░  ░░░░░      ░░░░░░   ░███░░░  
                                                                      ░███     
                                                                      █████    
                                                                     ░░░░░        by giardi ({VERSION})
'''

def file_to_list(filename: str) -> List[str] :
    try:
        with open(filename, 'r') as file:
            lines = [line.strip() for line in file.readlines()]
        return lines

    except Exception as e:
        print(f'\nAn error occurred while tryin\' to open {filename}\n\n{e}')

def get_all_abs_paths(file_or_directory: str) -> List[str]:
    file_paths = []
    try:
        if os.path.isfile(file_or_directory):
            file_paths.append(os.path.abspath(file_or_directory))
        elif os.path.isdir(file_or_directory):
            for root, dirs, files in os.walk(file_or_directory):
                for file in files:
                    absolute_path = os.path.abspath(os.path.join(root, file))
                    file_paths.append(absolute_path)
        else:
            print(f"Error: {file_or_directory} is not a valid file or directory")
        return file_paths
    except Exception as e:
        print(f"Error: {e}")
        return []

# Define ANSI color codes
RED = '\033[91m'
YELLOW = '\033[93m'
BLUE = '\033[94m'
RESET = '\033[0m'  # Reset color to default

def print_result(pattern: dict, result: str) -> tuple:
    name = pattern['pattern']['name']
    confidence = pattern['pattern']['confidence']

    # Color based on confidence level
    if confidence == "high":
        confidence_color = RED
    elif confidence == "medium":
        confidence_color = YELLOW
    elif confidence == "low":
        confidence_color = BLUE
    else:
        confidence_color = RESET

    # Print with colors
    print(f"[{name}] [{confidence_color}{confidence}{RESET}]\n{result}\n")

def getPatterns(pattern_name):
    config_file = os.path.expanduser("~/.config/infogrep.patterns.json")
    
    if not os.path.exists(config_file):
        raise FileNotFoundError(f"Config file not found: {config_file}")

    with open(config_file, 'r') as f:
        patterns_config = json.load(f)

    if pattern_name not in patterns_config:
        raise ValueError(f"Pattern '{pattern_name}' not found in config file")

    pattern_file = patterns_config[pattern_name]
    
    # Use pkg_resources to get the correct path to the pattern file
    pattern_file_path = pkg_resources.resource_filename('infogrep', pattern_file)
    
    if not os.path.exists(pattern_file_path):
        raise FileNotFoundError(f"Pattern file not found: {pattern_file_path}")

    with open(pattern_file_path, 'r') as f:
        patterns = yaml.safe_load(f)

    return patterns