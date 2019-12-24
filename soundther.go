package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

//Audio : Estructura para definir el tipo audio, que es lo que vamos a almacenar en el JSON.
type Audio struct {
	Title string `json:"title"`
	Data  []byte `json:"data"`
}

var audios []Audio

// Variables globales de uso común en todo el programa para la GUI.
var config, icon *fyne.StaticResource
var a fyne.App
var win fyne.Window
var logo *canvas.Image
var tab *widget.TabContainer
var filename string
var search string

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
			dir, _ := os.Getwd()
			a := exec.Command("play", dir+"/sounds/"+filename+".wav")
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
				dir, _ := os.Getwd()
				bytes, _ := ioutil.ReadFile(dir + "/sounds/" + filename + ".wav")
				jsonfile, err := ioutil.ReadFile("audios.json")
				if err != nil {
					log.Print(err)
				}
				err = json.Unmarshal(jsonfile, &audios)
				if err != nil {
					log.Print(err)
				}
				audios = append(audios, Audio{Title: filename, Data: bytes})
				jsonfile, err = json.Marshal(audios)
				if err != nil {
					log.Print(err)
				}
				ioutil.WriteFile("audios.json", jsonfile, 0644)
				emerg := a.NewWindow("Transacción realizada.")
				emerg.SetContent(widget.NewLabel("La transacción se ha realizado correctamente."))
				emerg.Show()
				initUI(0)
			})),
	)
	win.SetContent(tab)

}
func uploadAudio() fyne.Widget {
	nombre := widget.NewEntry()
	nombre.SetText("audio")
	form := widget.Form{OnSubmit: func() {
		filename = nombre.Text
	}}
	form.Append("Nombre del fichero:", nombre)

	return widget.NewVBox(
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewLabelWithStyle("Suba su audio aquí:", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewHBox(layout.NewSpacer(), &form, layout.NewSpacer()),
		layout.NewSpacer(),
		widget.NewHBox(widget.NewButtonWithIcon("Volver", theme.CancelIcon(), func() {
			initUI(1)
		}), layout.NewSpacer(),
			widget.NewButtonWithIcon("Subir", theme.ContentAddIcon(), func() {
				dir, _ := os.Getwd()
				_, err := os.Stat(dir + "/sounds/" + filename + ".wav")
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
				dir, _ := os.Getwd()
				a := exec.Command("rec", dir+"/sounds/"+filename+".wav")
				a.Start()
				time.Sleep(time.Second * 2)
				a.Process.Kill()
				bar.Stop()
				emerg.Close()
				win.Show()
				confirmAudio()
			})),
	)
}

func makeTransaction() fyne.CanvasObject {

	nombre := widget.NewEntry()
	apellido := widget.NewEntry()
	form := widget.Form{OnSubmit: func() {
		tab.CurrentTab().Content = uploadAudio()
		win.SetContent(tab)
	}}
	form.Append("Nombre", nombre)
	form.Append("Apellidos", apellido)
	return widget.NewVBox(widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewLabelWithStyle("Introduzca sus datos:", fyne.TextAlignCenter, fyne.TextStyle{}),
		&form,
		layout.NewSpacer())
}
func listTransactions() fyne.CanvasObject {

	jsonfile, err := ioutil.ReadFile("audios.json")
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(jsonfile, &audios)
	if err != nil {
		log.Print(err)
	}

	var str []string
	m := make(map[string][]byte)

	for _, value := range audios {
		str = append(str, value.Title)
		m[value.Title] = value.Data

	}

	g := widget.NewGroupWithScroller("Selecciona", widget.NewSelect(str, func(name string) {
		search = name
	}))

	return widget.NewVBox(widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		g,
		layout.NewSpacer(),
		widget.NewHBox(widget.NewButtonWithIcon("Volver", theme.CancelIcon(), func() {
			initUI(1)
		}),
			layout.NewSpacer(),
			widget.NewButton("Confirmar", func() {
				data, ok := m[search]
				if ok {
					dir, _ := os.Getwd()
					ioutil.WriteFile(dir+"/sounds/download/"+search+".wav", data, 0644)
					initUI(0)
					emerg := a.NewWindow("Transacción realizada.")
					emerg.SetContent(widget.NewVBox(widget.NewLabel("La transacción se realizó correctamente."),
						widget.NewLabel("El fichero se encuentra en sounds/download")))
					emerg.Show()
				} else {
					emerg := a.NewWindow("Error.")
					emerg.SetContent(widget.NewVBox(widget.NewLabel("No has seleccionado un audio válido")))
					emerg.Show()

				}
			})),
	)
}

func acceptTransaction() fyne.CanvasObject {

	nombre := widget.NewEntry()
	apellido := widget.NewEntry()
	form := widget.Form{OnSubmit: func() {
		tab.CurrentTab().Content = listTransactions()
		win.SetContent(tab)
	}}
	form.Append("Nombre", nombre)
	form.Append("Apellidos", apellido)
	return widget.NewVBox(widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewLabelWithStyle("Introduzca sus datos:", fyne.TextAlignCenter, fyne.TextStyle{}),
		&form,
		layout.NewSpacer())
}

func firstView() fyne.CanvasObject {
	logo := canvas.NewImageFromResource(icon)
	logo.SetMinSize(fyne.NewSize(200, 200))
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

func initUI(index int) {
	tab = widget.NewTabContainer(widget.NewTabItemWithIcon("Principal", theme.HomeIcon(), firstView()),
		widget.NewTabItemWithIcon("Subir audio", theme.ContentAddIcon(), makeTransaction()),
		widget.NewTabItemWithIcon("Descargar audio", theme.MenuDropDownIcon(), acceptTransaction()))
	tab.SetTabLocation(widget.TabLocationLeading)
	tab.SelectTabIndex(index)
	win.SetContent(tab)
}

func main() {
	filename = "audio"
	os.Setenv("FYNE_THEME", "light")
	os.Setenv("FYNE_SCALE", "1.0")
	// Recursos de la  aplicación
	icon = setResource("./icon.png", "icon")
	config = setResource("./config.png", "config")
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
