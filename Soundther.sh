#!/bin/bash
go build
cmd="./soundther"
$cmd &
cd ./blockchain/
python3 launch.py
