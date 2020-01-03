import json
from web3 import Web3, IPCProvider
import os
import sys
from scipy.io import wavfile
from tqdm import tqdm
import numpy as np

w3=None

if __name__ == "__main__":
    w3 = Web3(IPCProvider("./testNet/geth.ipc"))

    argumentos=sys.argv
    direccionContrato=argumentos[1]

    truffleFile = json.load(open('./build/contracts/Sound.json'))
    abi = truffleFile['abi']
    bytecode = truffleFile['bytecode']
    contrato= w3.eth.contract(bytecode=bytecode, abi=abi,address=direccionContrato)

    precio=w3.fromWei(contrato.functions.getPrecio().call(),"ether")

    print(precio)