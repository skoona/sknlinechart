package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/skoona/sknlinechart"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func makeChart(title, footer string) (sknlinechart.SknLineChart, error) {
	dataPoints := map[string][]*sknlinechart.LineChartDatapoint{} // legend, points

	rand.NewSource(1000.0)
	for x := 1; x < 130; x++ {
		val := rand.Float32() * 75.0
		if val > 75.0 {
			val = 75.0
		} else if val < 30.0 {
			val = 30.0
		}
		point := sknlinechart.NewLineChartDatapoint(val, theme.ColorBlue, time.Now().Format(time.RFC3339))
		dataPoints["Humidity"] = append(dataPoints["Humidity"], &point)
	}
	for x := 1; x < 130; x++ {
		val := rand.Float32() * 75.0
		if val > 95.0 {
			val = 95.0
		} else if val < 55.0 {
			val = 55.0
		}
		point := sknlinechart.NewLineChartDatapoint(val, theme.ColorRed, time.Now().Format(time.RFC3339))
		dataPoints["Temperature"] = append(dataPoints["Temperature"], &point)
	}

	lineChart, err := sknlinechart.NewLineChart(title, footer, &dataPoints)
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
	windowClosed := false

	gui := app.NewWithID("net.skoona.sknLineChart")
	w := gui.NewWindow("Custom Widget Development")

	w.SetOnClosed(func() {
		windowClosed = true
		fmt.Println("::main() Window Closed")
		time.Sleep(2 * time.Second)
	})

	lineChart, err := makeChart("Skoona Line Chart", "Example Time Series")

	go (func(chart sknlinechart.SknLineChart) {
		var many []*sknlinechart.LineChartDatapoint
		for x := 1; x < 121; x++ {
			val := rand.Float32() * 25.0
			if val > 50.0 {
				val = 50.0
			} else if val < 5.0 {
				val = 5.0
			}
			point := sknlinechart.NewLineChartDatapoint(val, theme.ColorPurple, time.Now().Format(time.RFC3339))
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
			point := sknlinechart.NewLineChartDatapoint(rand.Float32()*110.0, theme.ColorYellow, time.Now().Format(time.RFC3339))
			chart.ApplyDataPoint("SteadyStream", &point)
			if windowClosed {
				break
			}
			time.Sleep(time.Second)
		}
	})(lineChart)

	w.Resize(fyne.NewSize(982, 452))
	w.SetContent(container.NewPadded(lineChart))

	go func(w *fyne.Window) {
		systemSignalChannel := make(chan os.Signal, 1)
		signal.Notify(systemSignalChannel, syscall.SIGINT, syscall.SIGTERM)
		sig := <-systemSignalChannel // wait on ctrl-c
		windowClosed = true
		fmt.Println("Signal Received: ", sig.String())
		(*w).Close()
	}(&w)

	w.ShowAndRun()
	time.Sleep(1 * time.Second)
}
