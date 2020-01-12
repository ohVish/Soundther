# Makefile que permite la instalación de las dependencias de la aplicación.

start: 
	./Soundther.sh
dependencies:
	go get fyne.io/fyne/cmd/fyne
	python3 -m pip install numpy
	python3 -m pip install scipy
	python3 -m pip install web3
	python3 -m pip install tqdm
go:
	sudo add-apt-repository ppa:longsleep/golang-backports
	sudo apt-get update
	sudo apt-get install golang-go
python3:
	sudo add-apt-repository ppa:deadsnakes/ppa
	sudo apt update
	sudo apt install python3.7
common:
	sudo apt-get
	sudo apt-get install software-properties-common
	sudo add-apt-repository -y ppa:ethereum/ethereum
	sudo apt-get update
	sudo apt-get install ethereum
