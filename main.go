package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/sebastianrau/money-clock/pkg/gui"
	"github.com/sebastianrau/money-clock/pkg/logo"
)

func main() {

	a := app.NewWithID("com.github.sebastianrau.money-clock")
	a.SetIcon(logo.ResourceLogoPng)

	w := a.NewWindow("Hello World")

	gui.NewMainGui(w, a)

	w.ShowAndRun()
}
