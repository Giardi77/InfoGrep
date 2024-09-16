from typing import List
import os
import yaml
import json

Version = '0.1'

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
                                                                     ░░░░░        by giardi ({Version})
'''

def file_to_list(filename: str) -> List[str] :
    try:
        with open(filename, 'r') as file:
            lines = [line.strip() for line in file.readlines()]
        return lines

    except Exception as e:
        print(f'\nAn error occurred while tryin\' to open {filename}\n\n{e}')

def get_all_abs_paths(file_or_directory_raw: str) -> List[str]:
    file_or_directory = file_or_directory_raw.split(',')
    file_paths = []
    try:
        # Check if it's a single file
        if os.path.isfile(file_or_directory):
            # Add the absolute path of the file
            file_paths.append(os.path.abspath(file_or_directory))

        # If it's a directory
        else:
            for root, dirs, files in os.walk(file_or_directory):
                for file in files:
                    # Get the absolute path of each file
                    absolute_path = os.path.abspath(os.path.join(root, file))
                    file_paths.append(absolute_path)
        
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
    regex = pattern['pattern']['regex']
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

def getPatterns(name: str) -> dict:
    home_path = os.path.expanduser("~")

    with open(f"{home_path}/.config/infogrep.patterns.json", 'r') as json_file:
        my_patterns = json.load(json_file)
    with open(my_patterns[name], 'r') as file:
        Patterns = yaml.safe_load(file)

    return Patterns