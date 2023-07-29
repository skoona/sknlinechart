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
	for x := 1; x < 151; x++ {
		val := rand.Float32() * 75.0
		if val > 75.0 {
			val = 75.0
		} else if val < 30.0 {
			val = 30.0
		}
		point := lc.NewChartDatapoint(val, theme.ColorBlue, time.Now().Format(time.RFC1123))
		dataPoints["Humidity"] = append(dataPoints["Humidity"], &point)
	}
	for x := 1; x < 151; x++ {
		val := rand.Float32() * 75.0
		if val > 95.0 {
			val = 95.0
		} else if val < 55.0 {
			val = 55.0
		}
		point := lc.NewChartDatapoint(val, theme.ColorRed, time.Now().Format(time.RFC1123))
		dataPoints["Temperature"] = append(dataPoints["Temperature"], &point)
	}

	// yScalefactor is represents the topmost value on the yScale divided by 13
	// Ex: 650y = 650/13=50, also 130y is 130/13=10
	// 13 because there are 13 vertical divisions not including zero
	lineChart, err := lc.NewLineChart(title, footer, 55, &dataPoints)
	if err != nil {
		fmt.Println(err.Error())
		if lineChart == nil {
			panic(err.Error())
		}
	}
	lineChart.SetLineStrokeSize(2.0)
	lineChart.EnableDebugLogging(true)
	lineChart.SetTopLeftLabel("top left")

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
	logger := log.New(os.Stdout, "[DEBUG] ", log.Lmicroseconds|log.Lshortfile)

	gui := app.NewWithID("net.skoona.sknLineChart")
	w := gui.NewWindow("Custom Widget Development")

	lineChart, err := makeChart("Skoona Line Chart", "Example Time Series")

	go (func(chart lc.LineChart) {
		var many []*lc.ChartDatapoint
		for x := 1; x < 161; x++ {
			val := rand.Float32() * 25.0
			if val > 50.0 {
				val = 50.0
			} else if val < 5.0 {
				val = 5.0
			}
			point := lc.NewChartDatapoint(val, theme.ColorPurple, time.Now().Format(time.RFC1123))
			many = append(many, &point)
		}
		time.Sleep(10 * time.Second)
		err = lineChart.ApplyDataSeries("AllAtOnce", many)
		if err != nil {
			logger.Println("ApplyDataSeries", err.Error())
		}
		time.Sleep(time.Second)

		smoothed := lc.NewGraphAverage("SmoothStream", 32)
		for i := 0; i < 200; i++ {
			if windowClosed {
				break
			}
			dVal := float64(rand.Float32() * 512.0)
			smoother := smoothed.AddValue(dVal)
			point := lc.NewChartDatapoint(float32(smoother), theme.ColorYellow, time.Now().Format(time.RFC1123))
			chart.ApplyDataPoint("SmoothStream", &point)

			point2 := lc.NewChartDatapoint(float32(dVal), theme.ColorPurple, time.Now().Format(time.RFC1123))
			chart.ApplyDataPoint("SteadyStream", &point2)
			if windowClosed {
				break
			}
			time.Sleep(time.Second)
		}
	})(lineChart)

	lineChart.SetOnHoverPointCallback(func(series string, p lc.ChartDatapoint) {
		logger.Printf("Chart Datapoint Selected Callback: series:%s, point: %v\n", series, p)
	})

	w.SetContent(container.NewPadded(lineChart))
	w.Resize(fyne.NewSize(982, 452))
	w.CenterOnScreen()

	go func(w *fyne.Window, stopFlag chan os.Signal) {
		signal.Notify(stopFlag, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stopFlag // wait on ctrl-c
		windowClosed = true
		logger.Println("Signal Received: ", sig.String())
		exitCode = 1
		(*w).Close()
	}(&w, systemSignalChannel)

	w.ShowAndRun()

	os.Exit(exitCode)
}
