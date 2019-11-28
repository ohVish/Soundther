package main

import (
	"bufio"
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

var config, icon *fyne.StaticResource
var a fyne.App
var win fyne.Window
var logo *canvas.Image
var tab *widget.TabContainer
var filename string

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
		widget.NewLabelWithStyle("Prueba", fyne.TextAlignCenter, fyne.TextStyle{}),
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
			a := exec.Command("play", dir+"/"+filename+".wav")
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
				_, err := os.Stat("./" + filename + ".wav")
				if os.IsNotExist(err) {
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
				a := exec.Command("rec", dir+"/"+filename+".wav")
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
		widget.NewTabItemWithIcon("Subir audio", theme.ContentAddIcon(), makeTransaction()))
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
