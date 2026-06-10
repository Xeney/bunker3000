package main

import (
	"bunker3000/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(gui.NewTheme())

	w := a.NewWindow("Бункер-3000")
	w.Resize(fyne.NewSize(800, 600))
	w.SetMaster()
	w.CenterOnScreen()

	gui.ShowMainMenu(w)

	w.ShowAndRun()
}
