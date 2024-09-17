from setuptools import setup, find_packages
import os
import json
import sys

# Import version from _version.py
sys.path.append('src/infogrep')
from _version import VERSION

def create_default_config():
    home_path = os.path.expanduser("~")
    config_dir = os.path.join(home_path, ".config")
    config_file = os.path.join(config_dir, "infogrep.patterns.json")

    os.makedirs(config_dir, exist_ok=True)
    
    # Always create a fresh config with default patterns
    default_patterns = {
        "secrets": "default-patterns/rules-stable.yml",
        "pii": "default-patterns/pii-stable.yml"
    }
    
    with open(config_file, 'w') as f:
        json.dump(default_patterns, f, indent=2)

    print(f"Config file created/updated at: {config_file}")

setup(
    name="infogrep",
    version=VERSION,
    author="Giardi",
    description="Grep for sensitive info",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    url="https://github.com/Giardi77/InfoGrep",
    package_dir={"": "src"},
    packages=find_packages(where="src"),
    include_package_data=True,
    package_data={
        "infogrep": ["default-patterns/*.yml"],
    },
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

# Always run create_default_config() during setup
create_default_config()
