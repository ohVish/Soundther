@ECHO OFF
TITLE Soundther
SET dir="%~dp0"
IF "%1"=="" GOTO Continue
SET dir="%1"
:Continue
cd %dir%
go build
START /B soundther
cd ./blockchain/
python3 launch.py
PAUSE

