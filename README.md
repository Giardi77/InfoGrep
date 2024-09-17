![Logo](freeze.png)

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


# Installation

    pip install git+https://github.com/giardino-io/infogrep.git

### Uninstall

    pip uninstall infogrep

# Usage

The default pattern is 'secrets' wich points to default-patterns/rules-stable.yml, it contains a lot of regex for sensitive 
info such as **Api Keys** (aws, github and a lot more), **Asymmetric Private Keys** etc ...
Another pre-installed patterns yaml is the **'pii'**, containing a lot of regex for **emails, phone numbers and more**.

## Examples

Scan a file:

    infogrep -i file1.txt

Scan a directory:

    infogrep -i my_dir

Scan content from a pipe:

    cat file.js | infogrep

Add a custom pattern

    infogrep -a mypattern:/path/to/my_patterns.yaml

Scan with a custom pattern:

    infogrep -f file.js -p mypattern

# Contribute

if you find this tool helpfull and want to give a better/new regex or anything that can improve performace pull request will be welcomed!
