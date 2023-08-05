package sknlinechart

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"log"
	"os"
	"sync"
)

// ChartOption alternate methodof sett chart properties
type ChartOption func(lc *LineChartSkn) error

// ChartOptions container of chart options
type ChartOptions struct {
	opts []ChartOption
}

// NewChartOptions() returns a container for chart options
func NewChartOptions(opts ...ChartOption) *ChartOptions {
	if opts != nil {
		return &ChartOptions{
			opts: opts,
		}
	} else {
		return &ChartOptions{
			opts: []ChartOption{},
		}
	}
}

// Add adds a ChartOption to the container
func (o *ChartOptions) Add(opt ChartOption) {
	o.opts = append(o.opts, opt)
}

// Apply applies the ChartOption to the provided linechart : Internal Use Only
func (o *ChartOptions) Apply(lc *LineChartSkn) error {
	err := errors.New("")
	for _, opt := range o.opts {
		errOpt := opt(lc)
		if errOpt != nil {
			err = fmt.Errorf("%s\n%s", err.Error(), errOpt.Error())
		}
	}
	if len(err.Error()) < 10 {
		err = nil
	}
	return err
}

// NewLineChartViaOptions Create the Line Chart using ChartOptions model
// be careful not to exceed the series data point limit, which defaults to 150
//
// can return a valid chart object and an error object; errors really should be handled
// and are caused by data points exceeding the container limit of 150; they will be truncated
func NewLineChartViaOptions(options *ChartOptions) (LineChart, error) {

	w := &LineChartSkn{ // Create this widget with an initial text value
		dataPoints:              make(map[string][]*ChartDatapoint),
		dataPointStrokeSize:     2.0,
		dataSeriesAdded:         true,
		dataPointXLimit:         150,
		dataPointYLimit:         float32(10 * 13),
		chartScaleMultiplier:    10,
		enableDataPointMarkers:  true,
		enableHorizGridLines:    true,
		enableVertGridLines:     true,
		enableMousePointDisplay: true,
		enableColorLegend:       true,
		mouseDisplayStr:         "",
		mouseDisplayPosition:    &fyne.Position{},
		mouseDisplayFrameColor:  string(theme.ColorNameForeground),
		topLeftLabel:            "",
		topCenteredLabel:        "",
		topRightLabel:           "",
		leftMiddleLabel:         "",
		rightMiddleLabel:        "",
		bottomLeftLabel:         "",
		bottomCenteredLabel:     "",
		bottomRightLabel:        "",
		minSize:                 fyne.NewSize(320+theme.Padding()*4, 240+theme.Padding()*4),
		objectsCache:            []fyne.CanvasObject{}, // everything except datapoints, markers, and mousebox
		mapsLock:                sync.RWMutex{},
		logger:                  log.New(os.Stdout, "[DEBUG] ", log.Lmicroseconds|log.Lshortfile),
	}

	err := options.Apply(w)

	w.ExtendBaseWidget(w) // Initialize the BaseWidget
	return w, err
}

// Options

// WithTitle sets the top center label
func WithTitle(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.topCenteredLabel = label
		return nil
	}
}

// WithFooter set the bottom center label
func WithFooter(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.bottomCenteredLabel = label
		return nil
	}
}

// WithTopLeftLabel sets the chart label
func WithTopLeftLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.topLeftLabel = label
		return nil
	}
}

// WithTopRightLabel  sets the chart label
func WithTopRightLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.topRightLabel = label
		return nil
	}
}

// WithBottomLeftLabel  sets the chart label
func WithBottomLeftLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.bottomLeftLabel = label
		return nil
	}
}

// WithBottomRightLabel  sets the chart label
func WithBottomRightLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.bottomRightLabel = label
		return nil
	}
}

// WithLeftScaleLabel  sets the chart label
func WithLeftScaleLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.leftMiddleLabel = label
		return nil
	}
}

// WithRightScaleLabel  sets the chart label
func WithRightScaleLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.rightMiddleLabel = label
		return nil
	}
}

// WithYScaleFactor controls the yScale value y time 13 equals max y scale
func WithYScaleFactor(maxYScaleLabel int) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.dataPointYLimit = float32(maxYScaleLabel * 13)
		lc.chartScaleMultiplier = maxYScaleLabel
		return nil
	}
}

// WithMinSize sets the minimum x/y size of chart
func WithMinSize(width, height float32) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.minSize = fyne.NewSize(width+theme.Padding()*4, height+theme.Padding()*4)
		return nil
	}
}

// WithDataPointMarkers enables the use of markers on chart lines; mouse button 2 also toggles
func WithDataPointMarkers(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.enableDataPointMarkers = enable
		return nil
	}
}

// WithHorizGridLines enables horizontal grid lines display
func WithHorizGridLines(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.enableHorizGridLines = enable
		return nil
	}
}

// WithVertGridLines enables vertical grid lines display
func WithVertGridLines(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.enableVertGridLines = enable
		return nil
	}
}

// WithMousePointDisplay enables OnHover display over any line point
func WithMousePointDisplay(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.enableMousePointDisplay = enable
		return nil
	}
}

// WithColorLegend shows colored series legend in bottom right of chart
func WithColorLegend(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.enableColorLegend = enable
		return nil
	}
}

// WithDebugLogging activate logger to record method entry/exits
func WithDebugLogging(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.debugLoggingEnabled = enable
		return nil
	}
}

// WithOnHoverPointCallback set callback function for datapoint under mouse postion
func WithOnHoverPointCallback(callBack func(series string, dataPoint ChartDatapoint)) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.OnHoverPointCallback = callBack
		return nil
	}
}

// WithDataPoints Primary series data to initialize chart with
func WithDataPoints(seriesData map[string][]*ChartDatapoint) ChartOption {
	return func(lc *LineChartSkn) error {
		if seriesData == nil {
			return errors.New("dataPoint Params cannot be nil")
		}
		err := errors.New("")
		dpl := 150 // max xScale
		for key, points := range seriesData {
			cnt := len(points)
			if cnt > dpl {
				for len(points) > dpl {
					points = RemoveIndexFromSlice(0, points)
				}
				seriesData[key] = points
				err = fmt.Errorf("%s\n::NewLineChart() dataPoint contents exceeds the point count limit[Action: truncated leading]. Series: %s, points: %d, Limit: %d", err.Error(), key, cnt, dpl)
			}
		}
		for key, points := range seriesData {
			lc.dataPoints[key] = points
		}

		if len(err.Error()) < 10 {
			err = nil
		}

		return err
	}
}
