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
    direccion=argumentos[1]

    print(w3.fromWei(w3.eth.getBalance(direccion),'ether'))