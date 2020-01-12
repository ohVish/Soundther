#!/bin/bash
go build
cmd="./Soundther"
$cmd &
cd ./blockchain/
python3 launch.py
