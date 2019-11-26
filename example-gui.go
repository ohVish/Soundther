package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var config, icon *fyne.StaticResource

func setResource(file string, name string) *fyne.StaticResource {
	i, err := os.Open(file)

	if err != nil {
		log.Println(err)
	}
	b := bufio.NewReader(i)

	byte, err := ioutil.ReadAll(b)

	return fyne.NewStaticResource(name, byte)
}

func main() {
	os.Setenv("FYNE_THEME", "light")
	os.Setenv("FYNE_SCALE", "1.5")
	// Recursos de la  aplicación
	icon = setResource("./icon.png", "icon")
	config = setResource("./config.png", "config")

	a := app.New()
	win := a.NewWindow("")

	button := widget.NewButtonWithIcon("", config, func() {
		a.Quit()
	})
	button.Resize(fyne.NewSize(1000, 100))

	win.SetContent(fyne.NewContainerWithLayout(layout.NewFormLayout(), button, layout.NewSpacer()))
	win.SetTitle("Soundther")
	win.Resize(fyne.NewSize(500, 400))
	win.SetFixedSize(true)
	win.ShowAndRun()
}
