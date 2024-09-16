from setuptools import setup, find_packages
import os
import json
import utils

home_path = os.path.expanduser("~")
if os.path.isfile(f"{home_path}/.config/infogrep.patterns.json"):
    pass
else:
    Patterns = {}
    Patterns['secrets'] = os.path.abspath('default-patterns/rules-stable.yml')
    Patterns['pii'] = os.path.abspath('default-patterns/pii-stable.yml')

    
    with open(f"{home_path}/.config/infogrep.patterns.json", 'w') as json_file:
        json.dump(Patterns, json_file, indent=2)

setup(
    name = "InfoGrep",                          
    version = utils.Version,                              
    author = "Giardi",                         
    description = "Grep for sensitive info",  
    packages = find_packages(),
    py_modules = ["infogrep","utils"],              
    entry_points = {                              
        'console_scripts': [
            'infogrep=infogrep:main',        
        ],
    },
    install_requires = [                          
        'pyyaml',
        'requests',
    ],
)
