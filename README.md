![Logo](freeze.png)

<h4 align="center">🏎️💨 Grep for sensitive info FAST! 🏎️💨</h4>

<p align="center">
  <a href="#Features">Features</a> •
  <a href="#Installation">Installation</a> •
  <a href="#Usage">Usage</a> •
  <a href="#configuration">Config</a> •
</p>


# Features

- Grep ***files*** or ***directories*** for ****sensitive information**** using predefined patterns.
- Add custom patterns in YAML format.


# Installation


    git clone https://github.com/yourusername/infogrep.git
    cd infogrep
    sudo python3 setup.py install --user

### Uninstall

    pip uninstall infogrep

# Usage

The default pattern is 'secrets' wich points to default-patterns/rules-stable.yml, it contains a lot of regex for sensitive 
info such as **Api Keys** (aws, github and a lot more), **Asymmetric Private Keys** etc ...
Another pre-installed patterns yaml is the **'pii'**, containing a lot of regex for **emails, phone numbers and more**.

## Examples

Scan a files:

    infogrep -f file1.txt,file2.js

Scan a directory:

    infogrep -d my_dir

Add a custom pattern

    infogrep -a mypattern:/path/to/my_patterns.yaml

Scan with a custom pattern:

    infogrep -f file.js -p mypattern
