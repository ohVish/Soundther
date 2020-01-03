package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

//User :
type User struct {
	Name   string `json:"name"`
	Wallet string `json:"wallet"`
}

//Contract :
type Contract struct {
	Title  string `json:"title"`
	Wallet string `json:"wallet"`
	Price  int    `json:"price"`
	Type   string `json:"type"`
}

var usuarios []User
var audios []Contract
var walletLogin string
var walletBuy string
var prize string

// Variables globales de uso común en todo el programa para la GUI.
var a fyne.App
var win fyne.Window
var logo *canvas.Image
var tab *widget.TabContainer
var filename string
var search string
var duration int
var play string
var directory string

func setResource(file string, name string) *fyne.StaticResource {
	i, err := os.Open(file)

	if err != nil {
		log.Println(err)
	}
	b := bufio.NewReader(i)

	byte, err := ioutil.ReadAll(b)

	return fyne.NewStaticResource(name, byte)
}
func confirmAudio() {
	tab.CurrentTab().Content = widget.NewVBox(
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewLabelWithStyle("¿Estás seguro de que es el audio que quieres subir?", fyne.TextAlignCenter, fyne.TextStyle{}),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Reproducir", theme.MailSendIcon(), func() {
			emerg := a.NewWindow("Reproduciendo...")
			bar := widget.NewProgressBarInfinite()
			emerg.SetContent(bar)
			emerg.CenterOnScreen()
			emerg.FixedSize()
			emerg.Resize(fyne.NewSize(500, 100))
			win.Hide()
			emerg.Show()
			a := exec.Command("play", directory+"/sounds/"+filename+".wav")
			a.Run()
			bar.Stop()
			emerg.Close()
			win.Show()
		}),
		widget.NewHBox(
			widget.NewButtonWithIcon("Volver", theme.CancelIcon(), func() {
				initUI(1)
			}),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Confirmar", theme.ConfirmIcon(), func() {
				emerg := a.NewWindow("Realizando la transacción...")
				bar := widget.NewProgressBarInfinite()
				emerg.SetContent(bar)
				emerg.CenterOnScreen()
				emerg.FixedSize()
				emerg.Resize(fyne.NewSize(500, 100))
				win.Hide()
				emerg.Show()
				file := directory + "/sounds/" + filename + ".wav"
				worker := exec.Command("python3", "deploy.py", walletLogin, prize, file)
				var stderr bytes.Buffer
				worker.Stderr = &stderr
				worker.Dir = directory + "/blockchain/"
				bytes, err := worker.Output()
				if err != nil {
					emerg.Close()
					log.Println(stderr.String())
					emerg = a.NewWindow("Error")
					emerg.SetContent(widget.NewLabel("La transacción no se ha podido realizar."))
					emerg.Show()

				} else {
					//Leemos json con los contratos ya hechos
					jsonfile, err := ioutil.ReadFile("contracts.json")
					if err != nil {
						log.Print(err)
					}
					err = json.Unmarshal(jsonfile, &audios)
					if err != nil {
						log.Print(err)
					}
					wallet := string(bytes)
					args := strings.Split(wallet, "\n")
					precio, _ := strconv.Atoi(prize)
					audios = append(audios, Contract{Title: filename, Wallet: args[0], Price: precio, Type: args[1]})
					//Guardamos los cambios en el json.
					jsonfile, err = json.Marshal(audios)
					if err != nil {
						log.Print(err)
					}
					ioutil.WriteFile("contracts.json", jsonfile, 0644)
					emerg.Close()
					emerg := a.NewWindow("Transacción realizada.")
					emerg.SetContent(widget.NewLabel("La transacción se ha realizado correctamente."))
					emerg.Show()

				}
				win.Show()
				initUI(0)

			})),
	)
	win.SetContent(tab)

}
func uploadAudio() fyne.Widget {
	nombre := widget.NewEntry()
	nombre.SetText("audio")
	duracion := widget.NewEntry()
	duracion.SetText("2")
	dinero := widget.NewEntry()
	dinero.SetText("20")
	form := widget.Form{OnSubmit: func() {
		filename = nombre.Text
		duration, _ = strconv.Atoi(duracion.Text)
		prize = dinero.Text
	}}
	form.Append("Nombre del fichero:", nombre)
	form.Append("Duracion(Solo para grabar)", duracion)
	form.Append("Precio:", dinero)

	//Sacar el balance de la cuenta
	worker := exec.Command("python3", "balanceCuenta.py", walletLogin)
	worker.Dir = directory + "/blockchain/"
	balance, err := worker.Output()
	if err != nil {
		log.Println(err)
	}
	return widget.NewVBox(
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(),
			widget.NewLabelWithStyle("Saldo:"+string(balance[0:len(balance)-1])+" Ether", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true, Italic: false, Monospace: false}),
			layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), &form, layout.NewSpacer()),
		layout.NewSpacer(),
		widget.NewHBox(widget.NewButtonWithIcon("Volver", theme.CancelIcon(), func() {
			initUI(1)
		}), layout.NewSpacer(),
			widget.NewButtonWithIcon("Subir", theme.ContentAddIcon(), func() {
				_, err := os.Stat(directory + "/sounds/" + filename + ".wav")
				if err == nil {
					confirmAudio()
				} else {
					emerg := a.NewWindow("Error")
					emerg.SetContent(widget.NewLabel("No existe el fichero."))
					emerg.Show()
				}
			}),
			widget.NewButtonWithIcon("Grabar", theme.MailSendIcon(), func() {
				emerg := a.NewWindow("Grabando...")
				bar := widget.NewProgressBarInfinite()
				emerg.SetContent(bar)
				emerg.CenterOnScreen()
				emerg.FixedSize()
				emerg.Resize(fyne.NewSize(500, 100))
				win.Hide()
				emerg.Show()
				a := exec.Command("rec", directory+"/sounds/"+filename+".wav")
				a.Start()
				time.Sleep(time.Second * time.Duration(duration))
				a.Process.Kill()
				bar.Stop()
				emerg.Close()
				win.Show()
				confirmAudio()
			})),
	)
}

func makeTransaction() fyne.CanvasObject {
	//Cargamos los usuarios
	jsonfile, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(jsonfile, &usuarios)
	if err != nil {
		log.Print(err)
	}
	nombre := widget.NewEntry()
	form := widget.Form{OnSubmit: func() {
		emerg := a.NewWindow("Cargando información...")
		bar := widget.NewProgressBarInfinite()
		emerg.SetContent(bar)
		emerg.CenterOnScreen()
		emerg.FixedSize()
		emerg.Resize(fyne.NewSize(500, 100))
		win.Hide()
		emerg.Show()
		find := false
		for _, v := range usuarios {
			if v.Name == nombre.Text {
				walletLogin = v.Wallet
				find = true
				break
			}
		}
		if !find {
			var stderr bytes.Buffer
			worker := exec.Command("python3", directory+"/blockchain/nuevaCuenta.py")
			worker.Stderr = &stderr
			wallet, err := worker.Output()
			if err != nil {
				b := a.NewWindow("a")
				b.SetContent(widget.NewLabel(stderr.String()))
				b.Show()
				log.Println(err)
			}
			walletLogin = string(wallet[0 : len(wallet)-1])
			usuarios = append(usuarios, User{Name: nombre.Text, Wallet: walletLogin})
			jsonfile, err = json.Marshal(usuarios)
			if err != nil {
				log.Print(err)
			}
			ioutil.WriteFile("users.json", jsonfile, 0644)

		}
		tab.CurrentTab().Content = uploadAudio()
		win.SetContent(tab)
		emerg.Close()
		win.Show()
	}}
	form.Append("Usuario", nombre)
	return widget.NewVBox(widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewLabelWithStyle("Introduzca sus datos:", fyne.TextAlignCenter, fyne.TextStyle{}),
		&form,
		layout.NewSpacer())
}
func listTransactions() fyne.CanvasObject {

	//Leemos json con los contratos ya hechos
	jsonfile, err := ioutil.ReadFile("contracts.json")
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(jsonfile, &audios)
	if err != nil {
		log.Print(err)
	}

	var str []string

	for i, v := range audios {
		str = append(str, strconv.Itoa(i+1)+"."+v.Title+"| Precio: "+strconv.Itoa(v.Price)+" Ether")
	}
	group := widget.NewGroupWithScroller("Selecciona", widget.NewSelect(str, func(name string) {
		search = name
	}))

	//Sacar el balance de la cuenta
	worker := exec.Command("python3", "balanceCuenta.py", walletLogin)
	worker.Dir = directory + "/blockchain/"
	balance, err := worker.Output()
	if err != nil {
		log.Println(err)
	}

	return widget.NewVBox(widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(),
			widget.NewLabelWithStyle("Saldo:"+string(balance[0:len(balance)-1])+" Ether", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true, Italic: false, Monospace: false}),
			layout.NewSpacer()),
		group,
		layout.NewSpacer(),
		widget.NewHBox(widget.NewButtonWithIcon("Volver", theme.CancelIcon(), func() {
			initUI(1)
		}),
			layout.NewSpacer(),
			widget.NewButton("Confirmar", func() {
				emerg := a.NewWindow("Comprando...")
				bar := widget.NewProgressBarInfinite()
				emerg.SetContent(bar)
				emerg.CenterOnScreen()
				emerg.FixedSize()
				emerg.Resize(fyne.NewSize(500, 100))
				win.Hide()
				emerg.Show()
				//El primer caracter es el índice de la lista
				index := int(search[0] - '0')
				index = index - 1
				if search != "" {
					search = strings.Split(search, ".")[1]
					search = strings.Split(search, "|")[0]
					worker := exec.Command("python3", "comprar.py", audios[index].Wallet, walletLogin, directory+"/sounds/download/"+search+".wav", audios[index].Type)
					worker.Dir = directory + "/blockchain/"
					err := worker.Run()
					if err != nil {
						initUI(0)
						emerg.Close()
						win.Show()
						emerg = a.NewWindow("Error.")
						emerg.SetContent(widget.NewVBox(widget.NewLabel("No has seleccionado un audio válido")))
						emerg.Show()
					} else {
						initUI(0)
						emerg.Close()
						win.Show()
						emerg = a.NewWindow("Transacción realizada.")
						emerg.SetContent(widget.NewVBox(widget.NewLabel("La transacción se realizó correctamente."),
							widget.NewLabel("El fichero se encuentra en sounds/download")))
						emerg.Show()
					}
				} else {
					emerg := a.NewWindow("Error.")
					emerg.SetContent(widget.NewVBox(widget.NewLabel("No has seleccionado un audio válido")))
					emerg.Show()

				}
			})),
	)
}

func acceptTransaction() fyne.CanvasObject {
	//Cargamos los usuarios
	jsonfile, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(jsonfile, &usuarios)
	if err != nil {
		log.Print(err)
	}
	nombre := widget.NewEntry()
	form := widget.Form{OnSubmit: func() {
		emerg := a.NewWindow("Cargando información...")
		bar := widget.NewProgressBarInfinite()
		emerg.SetContent(bar)
		emerg.CenterOnScreen()
		emerg.FixedSize()
		emerg.Resize(fyne.NewSize(500, 100))
		win.Hide()
		emerg.Show()
		find := false
		for _, v := range usuarios {
			if v.Name == nombre.Text {
				walletLogin = v.Wallet
				find = true
				break
			}
		}
		if !find {
			worker := exec.Command("python3", directory+"/blockchain/nuevaCuenta.py")
			wallet, err := worker.Output()
			if err != nil {
				log.Println(err)
			}
			walletLogin = string(wallet[0 : len(wallet)-1])
			usuarios = append(usuarios, User{Name: nombre.Text, Wallet: walletLogin})
			jsonfile, err = json.Marshal(usuarios)
			if err != nil {
				log.Print(err)
			}
			ioutil.WriteFile("users.json", jsonfile, 0644)

		}
		tab.CurrentTab().Content = listTransactions()
		win.SetContent(tab)
		emerg.Close()
		win.Show()
	}}
	form.Append("Usuario", nombre)
	return widget.NewVBox(widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewLabelWithStyle("Introduzca sus datos:", fyne.TextAlignCenter, fyne.TextStyle{}),
		&form,
		layout.NewSpacer())
}

func firstView() fyne.CanvasObject {
	text := widget.NewLabelWithStyle("¡Bienvenidos a Soundther!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	credits := widget.NewLabelWithStyle("Aplicación de prueba realizada por Javier y José Miguel", fyne.TextAlignCenter, fyne.TextStyle{Bold: false})
	return widget.NewVBox(text,
		layout.NewSpacer(),
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Aplicación para subir tu música y audios de forma\n totalmente anónima utilizando la blockchain.", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		widget.NewVBox(widget.NewLabelWithStyle("_", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}), layout.NewSpacer(), credits, layout.NewSpacer()),
	)
}
func playView() fyne.CanvasObject {
	worker := exec.Command("ls")
	worker.Dir = directory + "/sounds/download/"
	files, _ := worker.Output()
	list := strings.Split(string(files), "\n")
	group := widget.NewGroupWithScroller("Selecciona el audio:", widget.NewSelect(list, func(name string) {
		play = name
	}))
	return widget.NewVBox(
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		layout.NewSpacer(),
		group,
		layout.NewSpacer(),
		widget.NewHBox(layout.NewSpacer(), widget.NewButtonWithIcon("Reproducir", theme.MailSendIcon(), func() {
			emerg := a.NewWindow("Reproduciendo...")
			bar := widget.NewProgressBarInfinite()
			emerg.SetContent(bar)
			emerg.CenterOnScreen()
			emerg.FixedSize()
			emerg.Resize(fyne.NewSize(500, 100))
			win.Hide()
			emerg.Show()
			a := exec.Command("play", directory+"/sounds/download/"+play)
			a.Run()
			bar.Stop()
			emerg.Close()
			win.Show()
		}), layout.NewSpacer()),
	)
}
func initUI(index int) {
	tab = widget.NewTabContainer(widget.NewTabItemWithIcon("Principal", theme.HomeIcon(), firstView()),
		widget.NewTabItemWithIcon("Subir audio", theme.ContentAddIcon(), makeTransaction()),
		widget.NewTabItemWithIcon("Descargar audio", theme.MenuDropDownIcon(), acceptTransaction()),
		widget.NewTabItemWithIcon("Reproducir audio", theme.MailSendIcon(), playView()))
	tab.SetTabLocation(widget.TabLocationLeading)
	tab.SelectTabIndex(index)
	win.SetContent(tab)
}

func main() {

	ex, _ := os.Executable()
	dirToFile := filepath.Dir(ex)
	os.Chdir(dirToFile)
	os.Chdir("../Resources/")
	directory, _ = os.Getwd()
	filename = "audio"
	duration = 2
	prize = "20"
	os.Setenv("FYNE_THEME", "light")
	os.Setenv("FYNE_SCALE", "1.0")
	// Recursos de la  aplicación
	icon := setResource("./icon.png", "icon")
	logo = canvas.NewImageFromResource(icon)
	logo.SetMinSize(fyne.NewSize(200, 200))
	a = app.New()
	win = a.NewWindow("Principal")
	win.SetTitle("Soundther")
	win.Resize(fyne.NewSize(300, 200))
	initUI(0)
	win.CenterOnScreen()
	win.ShowAndRun()
}
