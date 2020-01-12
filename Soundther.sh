#!/bin/bash
cd /Users/josemiguelaragonjurado/Soundther
go build
cmd="./soundther"
$cmd &
cd ./blockchain/
python3 launch.py
