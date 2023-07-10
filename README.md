# SknLineChart
Line chart with 120 horizontal, xscale, divisions displayed. The Y scale is limited to 100 divisions.

```go
package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/skoona/sknlinechart/pkg/components"
	"github.com/skoona/sknlinechart/pkg/entities"
	"github.com/skoona/sknlinechart/pkg/interfaces"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	gui := app.NewWithID("net.skoona.mq2influx")
	w := gui.NewWindow("Custom Widget Development")

	dataPoints := map[string][]interfaces.SknChartDatapoint{} // legend, points
	var points []interfaces.SknChartDatapoint
	var morePoints []interfaces.SknChartDatapoint
	var manyPoints []interfaces.SknChartDatapoint
	rand.NewSource(25.0)
	for x := 1; x < 50; x++ {
		points = append(points, entities.NewSknDatapoint(
			rand.Float32()*25.0,
			theme.ColorOrange,
			time.Now().Format(time.RFC3339)))
	}
	rand.NewSource(50.0)
	for x := 1; x < 125; x++ {
		morePoints = append(morePoints, entities.NewSknDatapoint(
			rand.Float32()*50.0,
			theme.ColorRed,
			time.Now().Format(time.RFC3339)))
	}
	rand.NewSource(75.0)
	for x := 1; x < 120; x++ {
		manyPoints = append(manyPoints, entities.NewSknDatapoint(
			rand.Float32()*75.0,
			theme.ColorPurple,
			time.Now().Format(time.RFC3339)))
	}

	dataPoints["first"] = points
	dataPoints["second"] = morePoints

	mw, err := components.NewSknLineChart("ggApcMon", "Time Series", &dataPoints)
	if err != nil {
		fmt.Println(err.Error())
	}

	go (func(chart interfaces.SknLineChart) {
		for {
			chart.ApplySingleDataPoint("steady", entities.NewSknDatapoint(
				rand.Float32()*100.0,
				theme.ColorYellow,
				time.Now().Format(time.RFC3339)))
			time.Sleep(time.Second)
		}
	})(mw)

	err = mw.ApplyNewDataSeries("many", manyPoints)
	if err != nil {
		fmt.Println("ApplyNewDataSeries", err.Error())
	}
	mw.EnableDataPointMarkers(false)

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

```

