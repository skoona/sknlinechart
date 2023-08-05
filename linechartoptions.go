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

type ChartOption func(lc *LineChartSkn) error
type ChartOptions struct {
	opts []ChartOption
}

func NewChartOptions(opts ...ChartOption) *ChartOptions {
	return &ChartOptions{
		opts: opts,
	}
}
func (o *ChartOptions) Add(opt ChartOption) {
	o.opts = append(o.opts, opt)
}
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

// NewLineChartWithOptions Create the Line Chart
// be careful not to exceed the series data point limit, which defaults to 150
//
// can return a valid chart object and an error object; errors really should be handled
// and are caused by data points exceeding the container limit of 150; they will be truncated
func NewLineChartWithOptions(options *ChartOptions) (LineChart, error) {

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

func WithTitle(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.topCenteredLabel = label
		return nil
	}
}
func WithFooter(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.bottomCenteredLabel = label
		return nil
	}
}
func WithTopLeftLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.topLeftLabel = label
		return nil
	}
}
func WithTopRightLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.topRightLabel = label
		return nil
	}
}
func WithBottomLeftLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.bottomLeftLabel = label
		return nil
	}
}
func WithBottomRightLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.bottomRightLabel = label
		return nil
	}
}
func WithLeftScaleLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.leftMiddleLabel = label
		return nil
	}
}
func WithRightScaleLabel(label string) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.rightMiddleLabel = label
		return nil
	}
}
func WithYScaleFactor(maxYScaleLabel int) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.dataPointYLimit = float32(maxYScaleLabel * 13)
		lc.chartScaleMultiplier = maxYScaleLabel
		return nil
	}
}
func WithMinSize(width, height float32) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.minSize = fyne.NewSize(width+theme.Padding()*4, height+theme.Padding()*4)
		return nil
	}
}
func WithDataPointMarkers(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.enableDataPointMarkers = enable
		return nil
	}
}
func WithHorizGridLines(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.enableHorizGridLines = enable
		return nil
	}
}
func WithVertGridLines(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.enableVertGridLines = enable
		return nil
	}
}
func WithMousePointDisplay(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.enableMousePointDisplay = enable
		return nil
	}
}
func WithColorLegend(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.enableColorLegend = enable
		return nil
	}
}
func WithDebugLogging(enable bool) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.debugLoggingEnabled = enable
		return nil
	}
}

func WithOnHoverPointCallback(callBack func(series string, dataPoint ChartDatapoint)) ChartOption {
	return func(lc *LineChartSkn) error {
		lc.OnHoverPointCallback = callBack
		return nil
	}
}
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
