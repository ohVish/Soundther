#!/bin/bash
cd /Users/josemiguelaragonjurado/Soundther/
cmd="./soundther"
$cmd &
cd ./blockchain/
python3 launch.py
