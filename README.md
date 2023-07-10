# SknLineChart
Line chart with 120 horizontal, xscale, divisions displayed. The Y scale is limited to 100 divisions.  Written in Go using the Fyne GUI framework.

## Features
* Multiple Series of datapoints rendered as a single line
* Series should be the same color. Each point in this chart accepts a themed color name
* 120 datapoint are displayed on the x scale of chart, with 100 as the default Y value.
* More than 120 causes the earliest points to be rolled off the screen; each series is independently scrolls
* Data points can be add at any time, cause the series to possible scroll automatically
* The 120 x limit will throw and error on creation of the chart, or on the replacement of its active dataset.
* Data point markers are toggled with mouse button 2
* Hovering over a data point will show a popup near the mouse pointer, showing series, value, index, and timestamp of data under mouse
* Mouse button 1 will turn off the sticky popup
* Labels are available for all four corners of window, include bottom and top centered titles
* left and right middle labels can be used as scale descriptions
* Any label left empty will not be displayed.
* Horizontal and Vertical chart grid lines can also be turned off/on

### SknLineChart Interface
```go
type SknLineChart interface {
	IsDataPointMarkersEnabled() bool // mouse button 2 toggles
	IsHorizGridLinesEnabled() bool
	IsVertGridLinesEnabled() bool
	IsMousePointDisplayEnabled() bool // hoverable and mouse button one

	SetDataPointMarkers(enable bool)
	SetHorizGridLines(enable bool)
	SetVertGridLines(enable bool)
	SetMousePointDisplay(enable bool)

	GetTopLeftLabel() string
	GetTitle() string
	GetTopRightLabel() string
	GetMiddleLeftLabel() string
	GetMiddleRightLabel() string
	GetBottomLeftLabel() string
	GetBottomCenteredLabel() string
	GetBottomRightLabel() string

	SetTopLeftLabel(newValue string)
	SetTitle(newValue string)
	SetTopRightLabel(newValue string)
	SetMiddleLeftLabel(newValue string)
	SetMiddleRightLabel(newValue string)
	SetBottomLeftLabel(newValue string)
	SetBottomCenteredLabel(newValue string)
	SetBottomRightLabel(newValue string)

	SetMinSize(s fyne.Size)

	ReplaceAllDataSeries(newData *map[string][]SknChartDatapoint) error
	ApplyNewDataSeries(seriesName string, newSeries []SknChartDatapoint) error
	ApplySingleDataPoint(seriesName string, newDataPoint SknChartDatapoint)
}
```

## Fyne Custom Widget Strategy
```go
/*
 * SknLineChart
 * Custom Fyne 2.0 Widget
 * Strategy
 * 1. Define Widget Named/Exported Struct
 *    1. export fields when possible
 * 2. Define Widget Renderer Named/unExported Struct
 *    1. un-exportable fields when possible
 * 3. Define NewWidget() *ExportedStruct method
 *    1. Define state variables for this widget
 *    2. Extend the BaseWidget
 *    3. Define Widget required methods
 *       1. CreateRenderer() fyne.WidgetRenderer, call newRenderer() below
 *    4. Define any methods required by additional interfaces, like
 *       desktop.Mouseable for mouse button support
 *       1. MouseDown(me MouseEvent)
 *       2. MouseUp(me MouseEvent)
 *       desktop.Hoverable for mouse movement support
 *       1. MouseIn(me MouseEvent)
 *       2. MouseMoved(me MouseEvent)  used to display data point under mouse
 *       3. MouseOut()
 *
 * 4. Define newRenderer() *notExportedStruct method
 *    1. Create canvas objects to be used in display
 *    2. Initialize there content if practical; not required
 *    3. Implement the required WidgetRenderer methods
 * 	  4. Refresh()               call refresh on each object
 * 	  5. Layout(s fyne.Size)     resize & move objects
 * 	  6. MinSize()  fyne.Size    return the minimum size needed
 * 	  7. Object() []fyne.Canvas  return the objects to be displayed
 * 	  8. Destroy()               cleanup if needed to prevent leaks
 * 5. In general state methods are the public api with or without getters/setters
 *    and the renderer creates the displayable objects, applies state/value to them, and
 *    manages their display.
 */
```


### Example
```go
package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/skoona/sknlinechart/pkg/components"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	windowClosed := false

	gui := app.NewWithID("net.skoona.mq2influx")
	w := gui.NewWindow("Custom Widget Development")

	w.SetOnClosed(func() {
		windowClosed = true
		fmt.Println("Window Closed")
		time.Sleep(2 * time.Second)
	})

	dataPoints := map[string][]components.SknChartDatapoint{} // legend, points
	var first, second, many []components.SknChartDatapoint

	rand.NewSource(50.0)
	for x := 1; x < 125; x++ {
		val := rand.Float32() * 100.0
		if val > 30.0 {
			val = 30.0
		} else if val < 5.0 {
			val = 5.0
		}
		first = append(first, components.NewSknDatapoint(
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
		second = append(second, components.NewSknDatapoint(
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
		many = append(many, components.NewSknDatapoint(
			val,
			theme.ColorPurple,
			time.Now().Format(time.RFC3339)))
	}

	dataPoints["first"] = first
	dataPoints["second"] = second

	lineChart, err := components.NewSknLineChart("Skoona Line Chart", "Example Time Series", &dataPoints)
	if err != nil {
		fmt.Println(err.Error())
	}

	go (func(chart components.SknLineChart) {
		time.Sleep(10 * time.Second)
		err = lineChart.ApplyNewDataSeries("many", many)
		if err != nil {
			fmt.Println("ApplyNewDataSeries", err.Error())
		}
		for {
			if windowClosed {
				break
			}
			chart.ApplySingleDataPoint("steady", components.NewSknDatapoint(
				rand.Float32()*110.0,
				theme.ColorYellow,
				time.Now().Format(time.RFC3339)))
			if windowClosed {
				break
			}
			time.Sleep(time.Second)
		}
	})(lineChart)

	w.Resize(fyne.NewSize(1024, 512))
	w.SetContent(lineChart)
	//w.SetContent(container.NewPadded(lineChart))

	go func(w *fyne.Window) {
		systemSignalChannel := make(chan os.Signal, 1)
		signal.Notify(systemSignalChannel, syscall.SIGINT, syscall.SIGTERM)
		sig := <-systemSignalChannel // wait on ctrl-c
		windowClosed = true
		fmt.Println(sig.String())
		(*w).Close()
	}(&w)

	w.ShowAndRun()
	time.Sleep(3 * time.Second)
}

```

### Project Layout
```
├── LICENSE
├── README.md
├── cmd
│   └── sknlinechart
│       └── main.go
├── go.mod
├── go.sum
└── pkg
    ├── commons
    │   └── mapsliceutils.go
    └── components
        ├── datapoint.go
        └── linechart.go

```


### Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request


### LICENSE
The application is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
