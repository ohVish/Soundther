import json
from web3 import Web3, IPCProvider
import sys
from scipy.io import wavfile
import numpy as np
from tqdm import tqdm

w3=None

def deployContract(contract,creador,precio,rate):
    construct_txn = contract.constructor(w3.toWei(precio,"ether"),rate).transact({
        'from': creador,
        'nonce': w3.eth.getTransactionCount(creador),
        'gas': 1728712,
        'gasPrice': w3.toWei('21', 'gwei')})

    tx_receipt = w3.eth.waitForTransactionReceipt(construct_txn)

    return tx_receipt['contractAddress']

if __name__ == "__main__":
    w3 = Web3(IPCProvider("./testNet/geth.ipc"))

    argumentos=sys.argv

    carteraCreador=argumentos[1]
    w3.geth.personal.unlockAccount(carteraCreador,'space treat blame exhibit tissue book decide fury exhaust hazard library effort',1500)

    precio=int(argumentos[2])
    fichero=argumentos[3]
    fs, data = wavfile.read(fichero)
    tipo = type(data[0])

    data = list(map(int,data.tolist()))


    truffleFile = json.load(open('./build/contracts/Sound.json'))
    abi = truffleFile['abi']
    bytecode = truffleFile['bytecode']
    contrato= w3.eth.contract(bytecode=bytecode, abi=abi)

    contract_address=deployContract(contrato,carteraCreador,precio,fs)

    contrato = w3.eth.contract(bytecode=bytecode, abi=abi, address=contract_address)

    n=100
    
    for i in range(0,len(data),n):
        contrato.functions.upload(data[i:i+n]).transact({
        'from': carteraCreador,
        'gas': 2277744,
        'gasPrice':w3.toWei(21,"gwei")})
    
    '''

    #Para hacer pruebas de gas
    print("Subiendo archivo: ")
    valoresGas=[]
    for i in tqdm(range(0,len(data),n)):
        valoresGas.append(float(contrato.functions.upload(data[i:i+n]).estimateGas({
        'from': carteraCreador,
        'gas': '0x15AD008',
        'gasPrice': '0x4E3B29200'})))

    print(valoresGas)
    print("\n Cantidad de gas:",max(valoresGas))
    print("Precio del gas:",w3.eth.gasPrice())'''
    

    w3.geth.personal.lockAccount(carteraCreador)

    print(contract_address)
    print(tipo)
    
    