package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
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

	dataPoints := map[string][]interfaces.SknDataSeries{} // legend, points
	points := []interfaces.SknDataSeries{}
	morePoints := []interfaces.SknDataSeries{}
	manyPoints := []interfaces.SknDataSeries{}
	rand.NewSource(25.0)
	for x := 1; x < 50; x++ {
		points = append(points, entities.NewSknDataSeries(
			rand.Float32()*25.0,
			theme.ColorOrange,
			time.Now().Format(time.RFC3339)))
	}
	rand.NewSource(50.0)
	for x := 1; x < 125; x++ {
		morePoints = append(morePoints, entities.NewSknDataSeries(
			rand.Float32()*50.0,
			theme.ColorRed,
			time.Now().Format(time.RFC3339)))
	}
	rand.NewSource(75.0)
	for x := 1; x < 120; x++ {
		manyPoints = append(manyPoints, entities.NewSknDataSeries(
			rand.Float32()*75.0,
			theme.ColorPurple,
			time.Now().Format(time.RFC3339)))
	}

	dataPoints["first"] = points
	dataPoints["second"] = morePoints

	mw, err := components.NewSknLineChart("ggApcMon", "Time Series", "Temperature", &dataPoints)
	if err != nil {
		fmt.Println(err.Error())
	}

	go (func() {
		for i := 0; i < 60; i++ {
			mw.ApplySingleDataPoint("steady", entities.NewSknDataSeries(
				89.0,
				theme.ColorYellow,
				time.Now().Format(time.RFC3339)))
			time.Sleep(time.Second)
		}
	})()

	err = mw.ApplyNewDataSeries("many", manyPoints)
	if err != nil {
		fmt.Println("ApplyNewDataSeries", err.Error())
	}

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
