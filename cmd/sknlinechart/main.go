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
	// yScalefactor  represents the value of each of the 13 Y divisions
	// Ex: 13 * 50 = 650, also 13 * 10 = 130,  50  and 10 being the yScalefactor value
	// 650 and 130 would be the top of the Y scale as there are 13 vertical divisions not including zero
	opts := lc.NewChartOptions()
	opts.Add(lc.WithDebugLogging(true))
	opts.Add(lc.WithFooter("With Options"))
	opts.Add(lc.WithTitle(title))
	opts.Add(lc.WithLeftScaleLabel("Temperature"))
	opts.Add(lc.WithRightScaleLabel("Humidity"))
	opts.Add(lc.WithDataPoints(dataPoints))
	opts.Add(lc.WithYScaleFactor(55))
	opts.Add(lc.WithOnHoverPointCallback(func(series string, p lc.ChartDatapoint) {
		fmt.Printf("Chart Datapoint Selected Callback: series:%s, point: %v\n", series, p)
	}))

	lineChart, err := lc.NewWithOptions(opts)
	if err != nil {
		fmt.Println(err.Error())
	}

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
		for i := 0; i < 300; i++ {
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
			time.Sleep(100 * time.Millisecond)
		}
	})(lineChart)

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
