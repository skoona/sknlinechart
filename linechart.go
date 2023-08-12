package sknlinechart

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"sync"
	"time"
)

/*
 * LineChart
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

import (
	"strings"
)

// LineChartSkn widget implements the LineChart interface
// to display multiple series of data points
// which will roll off older point beyond the  point limit.
type LineChartSkn struct {
	widget.BaseWidget       // Inherit from BaseWidget
	dataSeriesAdded         bool
	datapointAdded          bool
	dataPointStrokeSize     float32
	dataPointXLimit         int
	dataPointYLimit         float32
	chartScaleMultiplier    int
	enableDataPointMarkers  bool
	enableHorizGridLines    bool
	enableVertGridLines     bool
	enableMousePointDisplay bool
	enableColorLegend       bool
	topLeftLabel            string // The text to display in the widget
	topCenteredLabel        string
	topRightLabel           string
	leftMiddleLabel         string
	rightMiddleLabel        string
	bottomLeftLabel         string
	bottomCenteredLabel     string
	bottomRightLabel        string
	mouseDisplayStr         string
	mouseDisplayPosition    *fyne.Position
	mouseDisplayFrameColor  string
	dataPoints              map[string][]*ChartDatapoint
	minSize                 fyne.Size
	mapsLock                sync.RWMutex
	debugLoggingEnabled     bool
	logger                  *log.Logger
	// Private: Exposed for Testing; DO NOT USE
	objectsCache         []fyne.CanvasObject
	OnHoverPointCallback func(series string, dataPoint ChartDatapoint)
}

var _ LineChart = (*LineChartSkn)(nil)
var _ fyne.Widget = (*LineChartSkn)(nil)
var _ fyne.CanvasObject = (*LineChartSkn)(nil)

// NewLineChart Create the Line Chart
// be careful not to exceed the series data point limit, which defaults to 150
//
// can return a valid chart object and an error object; errors really should be handled
// and are caused by data points exceeding the container limit of 150; they will be truncated
func NewLineChart(topTitle, bottomTitle string, yScaleFactor int, dataPoints *map[string][]*ChartDatapoint) (LineChart, error) {
	return New(topTitle, bottomTitle, yScaleFactor, dataPoints)
}
func New(topTitle, bottomTitle string, yScaleFactor int, dataPoints *map[string][]*ChartDatapoint) (LineChart, error) {
	if dataPoints == nil {
		return nil, errors.New("dataPoint Params cannot be nil")
	}
	err := errors.New("")
	dpl := 150 // max xScale
	for key, points := range *dataPoints {
		cnt := len(points)
		if cnt > dpl {
			for len(points) > dpl {
				points = RemoveIndexFromSlice(0, points)
			}
			(*dataPoints)[key] = points
			err = fmt.Errorf("%s\n::NewLineChart() dataPoint contents exceeds the point count limit[Action: truncated leading]. Series: %s, points: %d, Limit: %d", err.Error(), key, cnt, dpl)
		}
	}
	if len(err.Error()) < 10 {
		err = nil
	}
	w := &LineChartSkn{ // Create this widget with an initial text value
		dataPoints:              *dataPoints,
		dataPointStrokeSize:     2.0,
		dataSeriesAdded:         true,
		dataPointXLimit:         dpl,
		dataPointYLimit:         float32(yScaleFactor * 13),
		chartScaleMultiplier:    yScaleFactor,
		enableDataPointMarkers:  true,
		enableHorizGridLines:    true,
		enableVertGridLines:     true,
		enableMousePointDisplay: true,
		enableColorLegend:       true,
		mouseDisplayStr:         "",
		mouseDisplayPosition:    &fyne.Position{},
		mouseDisplayFrameColor:  string(theme.ColorNameForeground),
		topLeftLabel:            "",
		topCenteredLabel:        topTitle,
		topRightLabel:           "",
		leftMiddleLabel:         "",
		rightMiddleLabel:        "",
		bottomLeftLabel:         "",
		bottomCenteredLabel:     bottomTitle,
		bottomRightLabel:        "",
		minSize:                 fyne.NewSize(320+theme.Padding()*4, 240+theme.Padding()*4),
		objectsCache:            []fyne.CanvasObject{}, // everything except datapoints, markers, and mousebox
		mapsLock:                sync.RWMutex{},
		logger:                  log.New(os.Stdout, "[DEBUG] ", log.Lmicroseconds|log.Lshortfile),
	}
	w.ExtendBaseWidget(w) // Initialize the BaseWidget
	return w, err
}

// CreateRenderer Create the renderer. This is called by the fyne application
func (w *LineChartSkn) CreateRenderer() fyne.WidgetRenderer {
	startTime := time.Now()
	w.debugLog("LineChartSkn::CreateRenderer()")
	r := newLineChartRenderer(w)
	w.debugLog("LineChartSkn::CreateRenderer() EXIT. Elapsed.microseconds: ", time.Until(startTime).Microseconds())
	return r
}

func (w *LineChartSkn) SetOnHoverPointCallback(f func(series string, dataPoint ChartDatapoint)) {
	w.OnHoverPointCallback = f
}

// SetMinSize set the minimum size limit for the linechart
func (w *LineChartSkn) SetMinSize(s fyne.Size) {
	w.debugLog("LineChartSkn::SetMinSize()")
	w.minSize = s
}

// GetTopLeftLabel return text from top left label
func (w *LineChartSkn) GetTopLeftLabel() string {
	return w.topLeftLabel
}

// GetTitle return text of the chart's title from top center
func (w *LineChartSkn) GetTitle() string {
	return w.topCenteredLabel
}

// IsDataPointMarkersEnabled returns state of chart's use of data point markers on series data
func (w *LineChartSkn) IsDataPointMarkersEnabled() bool {
	return w.enableDataPointMarkers
}

// IsHorizGridLinesEnabled returns state of chart's display of horizontal grid line
func (w *LineChartSkn) IsHorizGridLinesEnabled() bool {
	return w.enableHorizGridLines
}

// IsVertGridLinesEnabled returns state of chart's display of vertical grid line
func (w *LineChartSkn) IsVertGridLinesEnabled() bool {
	return w.enableVertGridLines
}

// IsColorLegendEnabled returns state of color legend at bottom right of chart
func (w *LineChartSkn) IsColorLegendEnabled() bool {
	return w.enableColorLegend
}

// IsMousePointDisplayEnabled return state of mouse popups when hovered over a chart datapoint
func (w *LineChartSkn) IsMousePointDisplayEnabled() bool {
	return w.enableMousePointDisplay
}

// GetLineStrokeSize sets thickness of all lines drawn
func (w *LineChartSkn) GetLineStrokeSize() float32 {
	return w.dataPointStrokeSize
}

// GetTopRightLabel returns text of top right label
func (w *LineChartSkn) GetTopRightLabel() string {
	return w.topRightLabel
}

// GetMiddleLeftLabel returns text of middle left label
func (w *LineChartSkn) GetMiddleLeftLabel() string {
	return w.leftMiddleLabel
}

// GetMiddleRightLabel returns text of middle right label
func (w *LineChartSkn) GetMiddleRightLabel() string {
	return w.rightMiddleLabel
}

// GetBottomLeftLabel returns text of bottom left label
func (w *LineChartSkn) GetBottomLeftLabel() string {
	return w.bottomLeftLabel
}

// GetBottomCenteredLabel returns text of bottom center label
func (w *LineChartSkn) GetBottomCenteredLabel() string {
	return w.bottomCenteredLabel
}

// GetBottomRightLabel returns text of bottom right label
func (w *LineChartSkn) GetBottomRightLabel() string {
	return w.bottomRightLabel
}

// SetLineStrokeSize sets thickness of all lines drawn
func (w *LineChartSkn) SetLineStrokeSize(newSize float32) {
	w.dataPointStrokeSize = newSize
}

// SetTopLeftLabel sets text to be display on chart at top left
func (w *LineChartSkn) SetTopLeftLabel(newValue string) {
	w.topLeftLabel = newValue
}

// SetTitle sets text to be display on chart at top center
func (w *LineChartSkn) SetTitle(newValue string) {
	w.topCenteredLabel = newValue
}

// SetTopRightLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetTopRightLabel(newValue string) {
	w.topRightLabel = newValue
}

// SetMiddleLeftLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetMiddleLeftLabel(newValue string) {
	w.leftMiddleLabel = newValue
}

// SetMiddleRightLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetMiddleRightLabel(newValue string) {
	w.rightMiddleLabel = newValue
}

// SetBottomLeftLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetBottomLeftLabel(newValue string) {
	w.bottomLeftLabel = newValue
}

// SetBottomRightLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetBottomRightLabel(newValue string) {
	w.bottomRightLabel = newValue
}

// SetBottomCenteredLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetBottomCenteredLabel(newValue string) {
	w.bottomCenteredLabel = newValue
}

// SetDataPointMarkers enables data point markers on display series points
func (w *LineChartSkn) SetDataPointMarkers(enable bool) {
	w.enableDataPointMarkers = enable
}

// SetHorizGridLines enables chart horizontal grid lines
func (w *LineChartSkn) SetHorizGridLines(enable bool) {
	w.enableHorizGridLines = enable
}

// SetColorLegend enables the color legend at bottom right on chart
func (w *LineChartSkn) SetColorLegend(enable bool) {
	w.enableColorLegend = enable
}

// SetVertGridLines enables chart vertical grid lines
func (w *LineChartSkn) SetVertGridLines(enable bool) {
	w.enableVertGridLines = enable
}

// SetMousePointDisplay true/false, enables data point display under mouse pointer
func (w *LineChartSkn) SetMousePointDisplay(enable bool) {
	w.enableMousePointDisplay = enable
}

// ApplyDataSeries adds a new series of data to existing chart set.
// throws error if new series exceeds containers point limit
func (w *LineChartSkn) ApplyDataSeries(seriesName string, newSeries []*ChartDatapoint) error {
	startTime := time.Now()

	w.debugLog("LineChartSkn::ApplyDataSeries() ENTER")
	if w == nil {
		w.debugLog("LineChartSkn::ApplyDataSeries() ERROR EXIT")
		return fmt.Errorf("ApplyDataSeries() no active widget")
	}

	if len(newSeries) <= w.dataPointXLimit {
		w.mapsLock.Lock()
		w.dataPoints[seriesName] = newSeries
		w.dataSeriesAdded = true
		w.mapsLock.Unlock()
		w.Refresh()
	} else {
		w.debugLog("LineChartSkn::ApplyDataSeries() ERROR EXIT")
		return fmt.Errorf("[%s] data series datapoints limit exceeded. limit:%d, count:%d", seriesName, w.dataPointXLimit, len(newSeries))
	}
	w.debugLog("LineChartSkn::ApplyDataSeries() EXIT. Elapsed.microseconds: ", time.Until(startTime).Microseconds())
	return nil
}

// ApplyDataPoint adds a new datapoint to an existing series
// will shift out the oldest point if containers limit is exceeded
func (w *LineChartSkn) ApplyDataPoint(seriesName string, newDataPoint *ChartDatapoint) {
	startTime := time.Now()

	w.debugLog("LineChartSkn::ApplyDataPoint() ENTER")
	if w == nil {
		return
	}

	w.mapsLock.Lock()

	if len(w.dataPoints[seriesName]) <= w.dataPointXLimit {
		w.dataPoints[seriesName] = append(w.dataPoints[seriesName], newDataPoint)
	} else {
		w.dataPoints[seriesName] = ShiftSlice(newDataPoint, w.dataPoints[seriesName])
	}
	w.datapointAdded = true
	w.mapsLock.Unlock()
	w.Refresh()
	w.debugLog("LineChartSkn::ApplyDataPoint() EXIT. Elapsed.microseconds: ", time.Until(startTime).Microseconds())
}

// Tapped From the Tappable Interface
func (w *LineChartSkn) Tapped(*fyne.PointEvent) {
	w.debugLog("LineChartSkn::Tapped() ENTER")
	w.enableMousePointDisplay = !w.enableMousePointDisplay
	w.Refresh()
	w.debugLog("LineChartSkn::Tapped() EXIT")
}

// TappedSecondary From the SecondaryTappable Interface
func (w *LineChartSkn) TappedSecondary(*fyne.PointEvent) {
	w.debugLog("LineChartSkn::TappedSecondary() ENTER")
	w.enableDataPointMarkers = !w.enableDataPointMarkers
	w.Refresh()
	w.debugLog("LineChartSkn::TappedSecondary() EXIT")
}

// MouseIn unused interface method
func (w *LineChartSkn) MouseIn(*desktop.MouseEvent) {
	w.debugLog("LineChartSkn::MouseIn()")
}

// MouseMoved interface method to discover which data point is under mouse
func (w *LineChartSkn) MouseMoved(me *desktop.MouseEvent) {
	startTime := time.Now()

	w.debugLog("LineChartSkn::MouseMoved() ENTER")
	if !w.enableMousePointDisplay {
		w.debugLog("LineChartSkn::MouseMoved(disabled) EXIT")
		return
	}
	w.mapsLock.Lock()
	matched := false

found:
	for key, points := range w.dataPoints {
		for idx, point := range points {
			top, bottom := (*point).MarkerPosition()
			if !me.Position.IsZero() && !top.IsZero() {
				if me.Position.X > top.X && me.Position.X < bottom.X &&
					me.Position.Y > top.Y-1 && me.Position.Y < bottom.Y {
					w.debugLog("MouseMoved() matched Mouse: ", me.Position, ", Top: ", top, ", Bottom: ", bottom)
					value := fmt.Sprint(key, ", Index: ", idx, ", Value: ", (*point).Value(), "    [", (*point).Timestamp(), "]")
					w.enableMouseContainer(value, (*point).ColorName(), &me.Position)
					if w.OnHoverPointCallback != nil {
						w.OnHoverPointCallback(strings.Clone(key), (*point).Copy())
					}
					matched = true
					break found
				}
			}
		}
	}
	w.mapsLock.Unlock()
	if matched {
		w.Refresh()
	}
	w.debugLog("LineChartSkn::MouseMoved() EXIT. Elapsed.microseconds: ", time.Until(startTime).Microseconds())
}

// MouseOut disable display of mouse data point display
func (w *LineChartSkn) MouseOut() {
	w.debugLog("LineChartSkn::MouseOut()")
	w.disableMouseContainer()
}

// enableMouseContainer private method to prepare values need by renderer to create pop display
// composes display text, captures position and colorName for use by renderer
func (w *LineChartSkn) enableMouseContainer(value, frameColor string, mousePosition *fyne.Position) *LineChartSkn {
	startTime := time.Now()
	w.debugLog("LineChartSkn::enableMouseContainer() ENTER")

	w.mouseDisplayStr = value
	w.mouseDisplayFrameColor = frameColor
	ct := canvas.NewText(value, theme.PrimaryColorNamed(frameColor))
	parts := strings.Split(value, "[")
	ts := fyne.MeasureText(parts[0], ct.TextSize, ct.TextStyle)
	mp := &fyne.Position{X: mousePosition.X - (ts.Width / 2), Y: mousePosition.Y - (3 * ts.Height) - theme.Padding()}
	w.mouseDisplayPosition = mp

	w.debugLog("LineChartSkn::enableMouseContainer() EXIT. Elapsed.microseconds: ", time.Until(startTime).Microseconds())
	return w
}

// disableMouseContainer private method to manage mouse leaving window
// blank string will prevent display
func (w *LineChartSkn) disableMouseContainer() {
	w.debugLog("LineChartSkn::disableMouseContainer()")
	w.mouseDisplayStr = ""
	w.Refresh()
}

// ObjectCount testing method return static object count
func (w *LineChartSkn) ObjectCount() int {
	w.debugLog("LineChartSkn::ObjectCount()")
	return len(w.objectsCache)
}

// EnableDebugLogging turns method entry/exit logging on or off
func (w *LineChartSkn) EnableDebugLogging(enable bool) {
	w.debugLoggingEnabled = enable
}
func (w *LineChartSkn) debugLog(a ...any) {
	if w.debugLoggingEnabled {
		_ = w.logger.Output(2, fmt.Sprint(a...))
	}
}
