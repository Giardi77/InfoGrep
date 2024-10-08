![Logo](Logo.png)

<h4 align="center">🏎️💨 Grep for sensitive info FAST! 🏎️💨</h4>

<p align="center">
  <a href="#Features">Features</a> •
  <a href="#Installation">Installation</a> •
  <a href="#Usage">Usage</a> •
  <a href="#Contribute">Contribute to this project</a>
</p>


# Features

- Grep ***files*** or ***directories*** for ****sensitive information**** using predefined patterns.
- Add custom patterns in YAML format.
- Parallel processing
- Get uniq results across multiple files


# Installation

    cargo install --git https://github.com/Giardi77/InfoGrep

# Usage

The default pattern is 'secrets' wich points to default-patterns/rules-stable.yml, it contains a lot of regex for sensitive 
info such as **Api Keys** (aws, github and a lot more), **Asymmetric Private Keys** etc ...
Another pre-installed patterns yaml is the **'pii'**, containing a lot of regex for **emails, phone numbers and more**.

## Examples

Scan a file:

    infogrep -i file1.txt

Scan a directory:

    infogrep -i my_dir

Add a custom pattern in ~/.config/infogrep.patterns.json with "name" : "/path/to/yaml.yml"

Scan with a custom pattern:

    infogrep -f file.js -p mypattern

Some regex might suck and match a lot of shit, you can use -t flag to truncate the output and see more results at once (default is 400 chars, if you want to see the whole thing set -t 0):

    infogrep -i my_dir -t 1000

Be **ULTRA** Fast:

    infogrep -i directory -w 8 #(8 workers in parallel)

# Contribute

if you find this tool helpfull and want to give a better/new regex or anything that can improve performace pull request will be welcomed!
