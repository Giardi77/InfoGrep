from setuptools import setup, find_packages
import os
import json
import sys

def create_default_config():
    home_path = os.path.expanduser("~")
    config_dir = os.path.join(home_path, ".config")
    config_file = os.path.join(config_dir, "infogrep.patterns.json")

    if not os.path.exists(config_file):
        os.makedirs(config_dir, exist_ok=True)
        default_patterns = {
            "secrets": os.path.join(os.path.dirname(__file__), "default-patterns", "rules-stable.yml"),
            "pii": os.path.join(os.path.dirname(__file__), "default-patterns", "pii-stable.yml")
        }
        with open(config_file, 'w') as f:
            json.dump(default_patterns, f, indent=2)
    else:
        # If the file exists, let's try to load it and re-save it to ensure proper formatting
        try:
            with open(config_file, 'r') as f:
                existing_patterns = json.load(f)
            with open(config_file, 'w') as f:
                json.dump(existing_patterns, f, indent=2)
        except json.JSONDecodeError:
            print(f"Error: The existing {config_file} is not a valid JSON file. Please delete it and run the installation again.")
            sys.exit(1)

    print(f"Config file created/updated at: {config_file}")

setup(
    name="infogrep",
    version="1.0.1",
    author="Giardi",
    description="Grep for sensitive info",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    url="https://github.com/Giardi77/InfoGrep",
    package_dir={"": "src"},
    packages=find_packages(where="src"),
    include_package_data=True,
    install_requires=[
        "pyyaml",
    ],
    entry_points={
        "console_scripts": [
            "infogrep=infogrep.infogrep:main",
        ],
    },
    classifiers=[
        "Programming Language :: Python :: 3",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)

create_default_config()
