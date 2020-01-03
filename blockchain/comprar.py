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
    direccion_contrato=argumentos[1]
    comprador=argumentos[2]
    ruta=argumentos[3]
    tipo=argumentos[4]

    saldoInicialComprador=w3.fromWei(w3.eth.getBalance(comprador),'ether')

    truffleFile = json.load(open('./build/contracts/Sound.json'))
    abi = truffleFile['abi']
    bytecode = truffleFile['bytecode']
    contrato= w3.eth.contract(bytecode=bytecode, abi=abi,address=direccion_contrato)

    precio=contrato.functions.getPrecio().call()
    print("Precio del audio:",precio)
    
    w3.geth.personal.unlockAccount(comprador,'space treat blame exhibit tissue book decide fury exhaust hazard library effort',1500)

    if comprador not in contrato.functions.compras().call():
        compra=contrato.functions.buySound().transact({ 'from':comprador,
                                                        'value':precio,
                                                        'gas': 1728712,
                                                        'gasPrice': w3.toWei('21', 'gwei')})

        tx_receipt = w3.eth.waitForTransactionReceipt(compra)
    else:
        print("Ya lo ha comprado")

    n=100
    audio=[]
    lenAudio=int(contrato.functions.getLen().call())
    fs = contrato.functions.getRate().call()
    print("Longitud del audio:",lenAudio,"\nFrecuencia de muestreo:",fs)

    for i in tqdm(range(0,lenAudio,n)):
        try:
            audio+=contrato.functions.download(i,i+n).call({"from":comprador})
        except:
            audio.append(0)

    switcher = {
        "numpy.int8":np.int8,
        "numpy.int16":np.int16,
        "numpy.int32":np.int32
    }
    tipo = tipo.split("'")[1]
    data=np.array(audio, dtype=switcher.get(tipo))

    wavfile.write(ruta,fs,data)

    saldoFinalComprador=w3.fromWei(w3.eth.getBalance(comprador),'ether')

    print("FACTURA DE LA OPERACION:",saldoInicialComprador,"-",w3.fromWei(precio,'ether'),'-',w3.fromWei(1728712*w3.toWei('21', 'gwei'),'ether'),"=",saldoFinalComprador)
