package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	lc "github.com/skoona/sknlinechart"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func makeChart(title, footer string) (lc.LineChart, error) {
	dataPoints := map[string][]*lc.ChartDatapoint{} // legend, points

	rand.NewSource(1000.0)
	for x := 1; x < 130; x++ {
		val := rand.Float32() * 75.0
		if val > 75.0 {
			val = 75.0
		} else if val < 30.0 {
			val = 30.0
		}
		point := lc.NewChartDatapoint(val, theme.ColorBlue, time.Now().Format(time.RFC3339))
		dataPoints["Humidity"] = append(dataPoints["Humidity"], &point)
	}
	for x := 1; x < 130; x++ {
		val := rand.Float32() * 75.0
		if val > 95.0 {
			val = 95.0
		} else if val < 55.0 {
			val = 55.0
		}
		point := lc.NewChartDatapoint(val, theme.ColorRed, time.Now().Format(time.RFC3339))
		dataPoints["Temperature"] = append(dataPoints["Temperature"], &point)
	}

	lineChart, err := lc.NewLineChart(title, footer, &dataPoints)
	if err != nil {
		fmt.Println(err.Error())
	}
	lineChart.EnableDebugLogging(true)
	lineChart.SetTopLeftLabel("top left")
	//lineChart.SetTopRightLabel("top right")

	lineChart.SetMiddleLeftLabel("Temperature")
	lineChart.SetMiddleRightLabel("Humidity")

	lineChart.SetBottomLeftLabel("bottom left")
	lineChart.SetBottomRightLabel("bottom right")

	return lineChart, err
}

func main() {
	systemSignalChannel := make(chan os.Signal, 1)
	exitCode := 0
	windowClosed := false

	gui := app.NewWithID("net.skoona.sknLineChart")
	w := gui.NewWindow("Custom Widget Development")

	lineChart, err := makeChart("Skoona Line Chart", "Example Time Series")

	go (func(chart lc.LineChart) {
		var many []*lc.ChartDatapoint
		for x := 1; x < 121; x++ {
			val := rand.Float32() * 25.0
			if val > 50.0 {
				val = 50.0
			} else if val < 5.0 {
				val = 5.0
			}
			point := lc.NewChartDatapoint(val, theme.ColorPurple, time.Now().Format(time.RFC3339))
			many = append(many, &point)
		}
		time.Sleep(10 * time.Second)
		err = lineChart.ApplyDataSeries("AllAtOnce", many)
		if err != nil {
			fmt.Println("ApplyDataSeries", err.Error())
		}
		time.Sleep(time.Second)
		for i := 0; i < 150; i++ {
			if windowClosed {
				break
			}
			point := lc.NewChartDatapoint(rand.Float32()*110.0, theme.ColorYellow, time.Now().Format(time.RFC3339))
			chart.ApplyDataPoint("SteadyStream", &point)
			if windowClosed {
				break
			}
			time.Sleep(time.Second)
		}
	})(lineChart)

	lineChart.SetOnHoverPointCallback(func(p lc.ChartDatapoint) {
		log.Printf("Chart Datapoint Selected Callback: %v\n", p)
	})

	w.SetContent(container.NewPadded(lineChart))
	w.Resize(fyne.NewSize(982, 452))

	go func(w *fyne.Window, stopFlag chan os.Signal) {
		signal.Notify(stopFlag, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stopFlag // wait on ctrl-c
		windowClosed = true
		fmt.Println("Signal Received: ", sig.String())
		exitCode = 1
		(*w).Close()
	}(&w, systemSignalChannel)

	w.ShowAndRun()

	os.Exit(exitCode)
}
