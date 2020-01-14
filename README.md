# Soundther 

<img src="https://github.com/ohVish/Soundther/blob/master/icon.png" width=200>  
Aplicación para subir y descargar ficheros de audio de la blockchain.  

## Prerrequisitos
Los prerrequisitos necesarios para utilizar la aplicación de subida y descarga de ficheros de audio, *Soundther*, son:
* **Compilador GCC en su versión 9.2.**
* **Python 3.**
* **Go 1.13.**
* **El paquete externo Fyne.io para Go.**
* **El paquete Web3 para Python 3.**
* **El paquete Numpy para Python 3.**
* **El paquete Scipy para Python 3.**
* **El paquete Tqdm para Python3.**
* **El software Geth.**
* **La biblioteca OpenGL.**
* **El software SoX (Sound eXchange).**

## Descarga
Para descargar la aplicación disponemos de un repositorio público donde se encuentran todos los ficheros para su correcto funcionamiento. El repositorio de github es el siguiente: https://github.com/ohVish/Soundther .

Para su descarga basta con utilizar el siguiente comando en la terminal:  
`git clone https://github.com/ohVish/Soundther`


## Instalación y uso
Una vez realizada la descarga de la aplicación, podemos realizar fácilmente la instalación en un sistema Linux Ubuntu, empleando el makefile que proporcionamos junto con la aplicación el cual se compone de estas cuatro instrucciones:  
* `make common`: Instala todas las dependencias básicas y el compilador gcc.
* `make go`: Instala Go en su versión 1.13.
* `make python3`: Instala Python en su versión 3.7. Si no existe un comando llamado python3, ejecuta la siguiente línea:   
   `sudo echo 'alias python3=python' $>>$~/.bashrc`  
   Sustituya el segundo *python* por el binario de Python 3 que quiera utilizar.
* `make install`: Empieza a instalar todos los paquetes de Python y Go necesarios para utilizar la aplicación. Además, también instala Geth.
* `make start`: Ejecuta un script para poner en marcha la aplicación y empezar a utilizarla.

Para una sencilla instalación simplemente con la orden `make install` del makefile instalará todo lo necesario para ejecutar la aplicación. Una vez ejecutada la orden `make start`, se mostrará interfaz de usuario inicial. La interfaz se compone principalmente de cuatro vistas principales, que se encuentran en la barra vertical que se encuentra a la izquierda de la aplicación en todo momento:

* La vista inicial, que es la que veremos al iniciar la aplicación, es meramente introductoria y tiene una leve descripción de la funcionalidad de la aplicación.
![Vista inicial](https://github.com/ohVish/Soundther/blob/master/pictures/inicio.png)

* La vista con la funcionalidad para subir archivos de audio. Al principio, nos pedirá el nombre de usuario que tenemos en la aplicación registrado y a su vez, este tiene asociado una cartera en la blockchain.
Una vez entremos, lo siguiente que vemos se nos muestra es el saldo disponible en nuestra cuenta junto con el nombre de nuestro archivo de audio, la duración que va a tener en el caso que queramos grabarlo y el precio al que queremos vender nuestro archivo de audio. Se nos permite grabar nuestro audio o subir uno que ya tengamos. La siguiente pantalla es para confirmar que queramos subir ese audio, permitiéndonos la reproducción del mismo. Una vez confirmamos, se realizará el contrato y se subirá el archivo de audio a la blockchain.
![Vista de subida](https://github.com/ohVish/Soundther/blob/master/pictures/subir.png)

* La vista con la funcionalidad de descargar archivos de audio. De forma similar a la vista anterior, al principio nos pedirá el nombre de usuario que tenemos registrado en la aplicación. Posteriormente, se muestra nuestro saldo y una lista con todas los archivos de audio que podemos comprar. Una vez seleccionemos el que queramos,se cargará el cobro a nuestra cartera y se descargará de la blockchain.
![Vista de compra](https://github.com/ohVish/Soundther/blob/master/pictures/descargar.png)

* La última vista es simplemente un reproductor de audio, para que todos los archivos de audio comprados, se puedan reproducir desde la misma aplicación.
![Vista de reproducción](https://github.com/ohVish/Soundther/blob/master/pictures/descargar.png)
