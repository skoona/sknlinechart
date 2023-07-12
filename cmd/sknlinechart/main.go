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

func main() {
	windowClosed := false

	gui := app.NewWithID("net.skoona.sknLineChart")
	w := gui.NewWindow("Custom Widget Development")

	w.SetOnClosed(func() {
		windowClosed = true
		fmt.Println("::main() Window Closed")
		time.Sleep(2 * time.Second)
	})

	dataPoints := map[string][]sknlinechart.LineChartDatapoint{} // legend, points
	var first, second, many []sknlinechart.LineChartDatapoint

	rand.NewSource(50.0)
	for x := 1; x < 125; x++ {
		val := rand.Float32() * 100.0
		if val > 30.0 {
			val = 30.0
		} else if val < 5.0 {
			val = 5.0
		}
		first = append(first, sknlinechart.NewLineChartDatapoint(
			val,
			theme.ColorOrange,
			time.Now().Format(time.RFC3339)))
	}
	for x := 1; x < 75; x++ {
		val := rand.Float32() * 40.0
		if val > 60.0 {
			val = 60.0
		} else if val < 35.0 {
			val = 35.0
		}
		second = append(second, sknlinechart.NewLineChartDatapoint(
			val,
			theme.ColorRed,
			time.Now().Format(time.RFC3339)))
	}
	for x := 1; x < 120; x++ {
		val := rand.Float32() * 75.0
		if val > 90.0 {
			val = 90.0
		} else if val < 65.0 {
			val = 65.0
		}
		many = append(many, sknlinechart.NewLineChartDatapoint(
			val,
			theme.ColorPurple,
			time.Now().Format(time.RFC3339)))
	}

	dataPoints["first"] = first
	dataPoints["second"] = second

	lineChart, err := sknlinechart.NewLineChart("Skoona Line Chart", "Example Time Series", &dataPoints)
	if err != nil {
		fmt.Println(err.Error())
	}

	go (func(chart sknlinechart.LineChart) {
		time.Sleep(10 * time.Second)
		err = lineChart.ApplyDataSeries("many", many)
		if err != nil {
			fmt.Println("ApplyDataSeries", err.Error())
		}
		time.Sleep(time.Second)
		for i := 0; i < 150; i++ {
			if windowClosed {
				break
			}
			chart.ApplyDataPoint("steady", sknlinechart.NewLineChartDatapoint(
				rand.Float32()*110.0,
				theme.ColorYellow,
				time.Now().Format(time.RFC3339)))
			if windowClosed {
				break
			}
			time.Sleep(time.Second)
		}
	})(lineChart)

	w.Resize(fyne.NewSize(982, 452))
	//w.SetContent(lineChart)
	w.SetContent(container.NewPadded(lineChart))

	go func(w *fyne.Window) {
		systemSignalChannel := make(chan os.Signal, 1)
		signal.Notify(systemSignalChannel, syscall.SIGINT, syscall.SIGTERM)
		sig := <-systemSignalChannel // wait on ctrl-c
		windowClosed = true
		fmt.Println(sig.String())
		(*w).Close()
	}(&w)

	w.ShowAndRun()
	time.Sleep(1 * time.Second)
}
