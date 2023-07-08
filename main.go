package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"ggApcMon/internal/entities"
	"ggApcMon/internal/interfaces"
	"ggApcMon/internal/views/components"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	gui := app.NewWithID("net.skoona.mq2influx")
	w := gui.NewWindow("Custom Widget Development")

	points := []interfaces.SknDataSeries{}
	rand.NewSource(100.0)

	for x := 1; x < 120; x++ {
		points = append(points, entities.NewSknDataSeries(float32(x),
			rand.Float32()*100.0,
			time.Now().Format(time.RFC3339)))
	}

	mw := components.NewSknLineChart("ggApcMon", "Time Series", "Temperature", points)
	w.Resize(fyne.NewSize(1024, 512))
	w.SetContent(mw)

	go func(a fyne.App) {
		systemSignalChannel := make(chan os.Signal, 1)
		signal.Notify(systemSignalChannel, syscall.SIGINT, syscall.SIGTERM)
		sig := <-systemSignalChannel // wait on ctrl-c
		//cancelService()              // provider
		fmt.Println(sig.String())
		a.Quit()
	}(gui)

	w.ShowAndRun()
	time.Sleep(3 * time.Second)
}
