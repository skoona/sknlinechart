package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"ggApcMon/internal/views/components"
)

func main() {
	app := app.New()
	w := app.NewWindow("Custom Widget Development")
	mw := components.NewMyWidget("ggApcMon")
	w.Resize(fyne.NewSize(400, 300))
	w.SetContent(mw)
	w.ShowAndRun()
}
