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
	"math"
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
	"fyne.io/fyne/v2/container"
	"image/color"
	"strconv"
	"strings"
)

// LineChartSkn widget implements the LineChart interface
// to display multiple series of data points
// which will roll off older point beyond the  point limit.
type LineChartSkn struct {
	widget.BaseWidget       // Inherit from BaseWidget
	datapointOrSeriesAdded  bool
	dataPointXLimit         int
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
	dataPointScale          fyne.Size
	minSize                 fyne.Size
	propertiesLock          sync.RWMutex
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
func NewLineChart(topTitle, bottomTitle string, dataPoints *map[string][]*ChartDatapoint) (LineChart, error) {
	if dataPoints == nil {
		return nil, errors.New("dataPoint Params cannot be nil")
	}
	err := errors.New("")
	dpl := 150
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
		datapointOrSeriesAdded:  true,
		dataPointXLimit:         dpl,
		dataPointScale:          fyne.NewSize(float32(dpl), 130.0), // max x/y scales, and x data points
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
		minSize:                 fyne.NewSize(420+theme.Padding()*4, 315+theme.Padding()*4),
		objectsCache:            []fyne.CanvasObject{}, // everything except datapoints, markers, and mousebox
		propertiesLock:          sync.RWMutex{},
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
		w.propertiesLock.Lock()
		w.dataPoints[seriesName] = newSeries
		w.datapointOrSeriesAdded = true
		w.propertiesLock.Unlock()
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

	w.propertiesLock.Lock()

	if len(w.dataPoints[seriesName]) <= w.dataPointXLimit {
		w.dataPoints[seriesName] = append(w.dataPoints[seriesName], newDataPoint)
	} else {
		w.dataPoints[seriesName] = ShiftSlice(newDataPoint, w.dataPoints[seriesName])
	}
	w.datapointOrSeriesAdded = true
	w.propertiesLock.Unlock()
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
	w.propertiesLock.Lock()
	matched := false

found:
	for key, points := range w.dataPoints {
		for idx, point := range points {
			top, bottom := (*point).MarkerPosition()
			if !me.Position.IsZero() && !top.IsZero() {
				if me.Position.X > top.X && me.Position.X < bottom.X &&
					me.Position.Y > top.Y-1 && me.Position.Y < bottom.Y {
					w.debugLog("MouseMoved() matched Mouse: ", me.Position, ", Top: ", top, ", Bottom: ", bottom)
					value := fmt.Sprint(key, ", Index: ", idx, ", Value: ", (*point).Value(), "    \n[", (*point).Timestamp(), "]")
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
	w.propertiesLock.Unlock()
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
	parts := strings.Split(value, "\n[")
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
		//w.logger.Println(a...)
		_ = w.logger.Output(2, fmt.Sprint(a...))
	}
}

// Widget Renderer code starts here
type lineChartRenderer struct {
	widget                *LineChartSkn     // Reference to the widget holding the current state
	chartFrame            *canvas.Rectangle // A chartFrame rectangle
	xInc                  float32
	yInc                  float32
	dataPoints            map[string][]*canvas.Line
	dataPointMarkers      map[string][]*canvas.Circle
	mouseDisplayContainer *fyne.Container
	xLines                []*canvas.Line
	yLines                []*canvas.Line
	xLabels               []*canvas.Text
	yLabels               []*canvas.Text
	topLeftDesc           *canvas.Text
	topCenteredDesc       *canvas.Text
	topRightDesc          *canvas.Text
	bottomLeftDesc        *canvas.Text
	bottomCenteredDesc    *canvas.Text
	bottomRightDesc       *canvas.Text
	leftMiddleBox         *fyne.Container
	rightMiddleBox        *fyne.Container
	colorLegend           *fyne.Container
}

var _ fyne.WidgetRenderer = (*lineChartRenderer)(nil)

// Create the renderer with a reference to the widget
// and all the objectsCache to be displayed for this custom widget
//
// Note: Do not size or move canvas objects here.
func newLineChartRenderer(lineChart *LineChartSkn) fyne.WidgetRenderer {
	lineChart.debugLog("::newLineChartRenderer() ENTER")
	startTime := time.Now()
	lineChart.propertiesLock.Lock()
	defer lineChart.propertiesLock.Unlock()

	var (
		dataPoints       = map[string][]*canvas.Line{}
		dpMaker          = map[string][]*canvas.Circle{}
		objs             []fyne.CanvasObject
		xlines, ylines   []*canvas.Line
		xLabels, yLabels []*canvas.Text
	)

	background := canvas.NewRectangle(color.Transparent)
	background.StrokeWidth = 0.75
	background.StrokeColor = theme.PrimaryColorNamed(theme.ColorBlue)
	objs = append(objs, background)

	border := canvas.NewRectangle(theme.OverlayBackgroundColor())
	border.StrokeColor = theme.PrimaryColorNamed(lineChart.mouseDisplayFrameColor)
	border.StrokeWidth = 2.0

	legend := widget.NewLabel("")
	legend.Alignment = fyne.TextAlignCenter
	legend.Wrapping = fyne.TextWrapWord
	legend.TextStyle = fyne.TextStyle{
		Bold:   true,
		Italic: true,
	}
	mouseDisplay := container.NewPadded(
		border,
		legend,
	)
	mouseDisplay.Hide()

	for i := 0; i < 13; i++ {
		x := canvas.NewLine(theme.PrimaryColorNamed(theme.ColorGreen))
		x.StrokeWidth = 0.25
		y := canvas.NewLine(theme.PrimaryColorNamed(theme.ColorGreen))
		y.StrokeWidth = 0.25
		xlines = append(xlines, x)
		ylines = append(ylines, y)
		objs = append(objs, x, y)
	}
	x := canvas.NewLine(theme.PrimaryColorNamed(theme.ColorGreen))
	x.StrokeWidth = 0.25
	xlines = append(xlines, x)
	objs = append(objs, x)

	for i := 0; i < 14; i++ {
		yt := strconv.Itoa((13 - i) * 10)
		yl := canvas.NewText(yt, theme.ForegroundColor())
		yl.Alignment = fyne.TextAlignTrailing
		yLabels = append(yLabels, yl)
		objs = append(objs, yl)
	}
	for i := 0; i < 16; i++ {
		xt := strconv.Itoa(i * 10)
		xl := canvas.NewText(xt, theme.ForegroundColor())
		xl.Alignment = fyne.TextAlignTrailing
		xLabels = append(xLabels, xl)
		objs = append(objs, xl)
	}

	colorLegend := container.NewHBox()
	for key, points := range lineChart.dataPoints {
		for _, point := range points {
			x := canvas.NewLine(theme.PrimaryColorNamed((*point).ColorName()))
			x.StrokeWidth = 2.0
			dataPoints[key] = append(dataPoints[key], x)
			z := canvas.NewCircle(theme.PrimaryColorNamed((*point).ColorName()))
			z.StrokeWidth = 4.0
			z.Resize(fyne.NewSize(5, 5))
			dpMaker[key] = append(dpMaker[key], z)
		}
		z := canvas.NewText(key, theme.PrimaryColorNamed((*points[0]).ColorName()))
		colorLegend.Add(z)
	}

	topCenteredDesc := canvas.NewText(lineChart.topCenteredLabel, theme.ForegroundColor())
	topCenteredDesc.TextSize = 24
	topCenteredDesc.TextStyle = fyne.TextStyle{
		Bold:   true,
		Italic: false,
	}
	objs = append(objs, topCenteredDesc)

	bottomCenteredDesc := canvas.NewText(lineChart.bottomCenteredLabel, theme.ForegroundColor())
	bottomCenteredDesc.TextSize = 16
	bottomCenteredDesc.TextStyle = fyne.TextStyle{
		Bold:   false,
		Italic: true,
	}
	objs = append(objs, bottomCenteredDesc)

	// vertical text for X/Y legends since no text rotation is available
	lBox := container.NewVBox()
	for _, c := range lineChart.leftMiddleLabel {
		z := canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground)))
		z.TextStyle = fyne.TextStyle{Monospace: true}
		z.TextSize = 14
		z.Alignment = fyne.TextAlignCenter
		lBox.Add(z)
	}
	objs = append(objs, lBox)

	rBox := container.NewVBox()
	for _, c := range lineChart.rightMiddleLabel {
		z := canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground)))
		z.TextStyle = fyne.TextStyle{Monospace: true}
		z.TextSize = 14
		z.Alignment = fyne.TextAlignCenter
		rBox.Add(z)
	}
	objs = append(objs, rBox)

	tl := canvas.NewText(lineChart.topLeftLabel, theme.ForegroundColor())
	tr := canvas.NewText(lineChart.topRightLabel, theme.ForegroundColor())
	bl := canvas.NewText(lineChart.bottomLeftLabel, theme.ForegroundColor())
	br := canvas.NewText(lineChart.bottomRightLabel, theme.ForegroundColor())
	objs = append(objs, tl, tr, bl, br)

	// save all except data points, markers, and mouse box
	lineChart.objectsCache = append(lineChart.objectsCache, objs...)

	lineChart.debugLog("::newLineChartRenderer() EXIT. Elapsed.microseconds: ", time.Until(startTime).Microseconds())

	return &lineChartRenderer{
		widget:                lineChart,
		chartFrame:            background,
		xLines:                xlines,
		yLines:                ylines,
		xLabels:               xLabels,
		yLabels:               yLabels,
		dataPoints:            dataPoints,
		topLeftDesc:           tl,
		topCenteredDesc:       topCenteredDesc,
		topRightDesc:          tr,
		bottomLeftDesc:        bl,
		bottomCenteredDesc:    bottomCenteredDesc,
		bottomRightDesc:       br,
		leftMiddleBox:         lBox,
		rightMiddleBox:        rBox,
		dataPointMarkers:      dpMaker,
		mouseDisplayContainer: mouseDisplay,
		colorLegend:           colorLegend,
	}
}

// manageLabelVisibility called by refresh to show/hide as needed
func (r *lineChartRenderer) manageLabelVisibility() {
	startTime := time.Now()
	r.widget.debugLog("lineChartRenderer::manageLabelVisibility() ENTER")
	if r.topLeftDesc.Text != "" {
		if !r.topLeftDesc.Visible() {
			r.topLeftDesc.Show()
		}
	} else {
		r.topLeftDesc.Hide()
	}
	if r.topCenteredDesc.Text != "" {
		if !r.topCenteredDesc.Visible() {
			r.topCenteredDesc.Show()
		}
	} else {
		r.topCenteredDesc.Hide()
	}
	if r.topRightDesc.Text != "" {
		if !r.topRightDesc.Visible() {
			r.topRightDesc.Show()
		}
	} else {
		r.topRightDesc.Hide()
	}
	if r.widget.leftMiddleLabel != "" {
		if !r.leftMiddleBox.Visible() {
			r.leftMiddleBox.Show()
		}
	} else {
		r.leftMiddleBox.Hide()
	}
	if r.widget.rightMiddleLabel != "" {
		if !r.rightMiddleBox.Visible() {
			r.rightMiddleBox.Show()
		}
	} else {
		r.rightMiddleBox.Hide()
	}
	if r.bottomLeftDesc.Text != "" {
		if !r.bottomLeftDesc.Visible() {
			r.bottomLeftDesc.Show()
		}
	} else {
		r.bottomLeftDesc.Hide()
	}
	if r.bottomCenteredDesc.Text != "" {
		if !r.bottomCenteredDesc.Visible() {
			r.bottomCenteredDesc.Show()
		}
	} else {
		r.bottomCenteredDesc.Hide()
	}
	if r.bottomRightDesc.Text != "" {
		if !r.bottomRightDesc.Visible() {
			r.bottomRightDesc.Show()
		}
	} else {
		r.bottomRightDesc.Hide()
	}
	if r.widget.enableColorLegend {
		if r.colorLegend.Hidden {
			r.colorLegend.Show()
		}
	} else {
		if !r.colorLegend.Hidden {
			r.colorLegend.Hide()
		}
	}

	for _, line := range r.xLines {
		if r.widget.enableHorizGridLines {
			if !line.Visible() {
				line.Show()
			}
		} else {
			line.Hide()
		}
	}
	for _, line := range r.yLines {
		if r.widget.enableVertGridLines {
			if !line.Visible() {
				line.Show()
			}
		} else {
			line.Hide()
		}
	}
	r.widget.debugLog("lineChartRenderer::manageLabelVisibility() EXIT. Elapsed.microseconds: ", time.Until(startTime).Microseconds())
}

// Refresh method is called if the state of the widget changes or the
// theme is changed
func (r *lineChartRenderer) Refresh() {
	r.widget.debugLog("lineChartRenderer::Refresh() ENTER")
	startTime := time.Now()

	if r.widget.datapointOrSeriesAdded {
		r.verifyDataPoints()
	}

	r.leftMiddleBox.RemoveAll()
	for _, c := range r.widget.leftMiddleLabel {
		z := canvas.NewText(
			strings.ToUpper(string(c)),
			theme.PrimaryColorNamed(string(theme.ColorNameForeground)))
		z.TextSize = 14
		z.TextStyle = fyne.TextStyle{Monospace: true}
		z.Alignment = fyne.TextAlignCenter
		r.leftMiddleBox.Add(z)
	}
	r.leftMiddleBox.Refresh()

	r.rightMiddleBox.RemoveAll()
	for _, c := range r.widget.rightMiddleLabel {
		z := canvas.NewText(
			strings.ToUpper(string(c)),
			theme.PrimaryColorNamed(string(theme.ColorNameForeground)))
		z.TextSize = 14
		z.TextStyle = fyne.TextStyle{Monospace: true}
		z.Alignment = fyne.TextAlignCenter

		r.rightMiddleBox.Add(z)
	}
	r.rightMiddleBox.Refresh()

	r.widget.propertiesLock.RLock()
	r.topLeftDesc.Text = r.widget.topLeftLabel
	r.topCenteredDesc.Text = r.widget.topCenteredLabel
	r.topRightDesc.Text = r.widget.topRightLabel
	r.bottomLeftDesc.Text = r.widget.bottomLeftLabel
	r.bottomCenteredDesc.Text = r.widget.bottomCenteredLabel
	r.bottomRightDesc.Text = r.widget.bottomRightLabel
	for _, v := range r.widget.objectsCache {
		v.Refresh()
	}
	r.widget.datapointOrSeriesAdded = false

	r.manageLabelVisibility()

	r.widget.propertiesLock.RUnlock()

	r.mouseDisplayContainer.Hide()
	r.widget.propertiesLock.Lock()
	r.mouseDisplayContainer.Objects[0].(*canvas.Rectangle).StrokeColor = theme.PrimaryColorNamed(r.widget.mouseDisplayFrameColor)
	r.mouseDisplayContainer.Objects[1].(*widget.Label).SetText(r.widget.mouseDisplayStr)
	r.widget.propertiesLock.Unlock()
	if r.widget.enableMousePointDisplay {
		if r.widget.mouseDisplayStr != "" {
			if !r.mouseDisplayContainer.Visible() {
				r.mouseDisplayContainer.Show()
			}
		} else {
			r.mouseDisplayContainer.Hide()
		}
	} else {
		r.mouseDisplayContainer.Hide()
	}

	r.widget.debugLog("lineChartRenderer::Refresh() EXIT. Elapsed.microseconds: ", time.Until(startTime).Microseconds())
}

// layoutSeries layout one series to position new elements
func (r *lineChartRenderer) layoutSeries(series string) {
	startTime := time.Now()

	r.widget.debugLog("lineChartRenderer::layoutSeries() ENTER. Series: ", series)
	// data points
	xp := r.xInc
	yp := r.yInc * 14
	yScale := (r.yInc * 10) / 100
	xScale := (r.xInc * 10) / 100
	var dp float32
	data := r.widget.dataPoints[series] // datasource
	lastPoint := fyne.NewPos(xp, yp)

	for idx, point := range data { // one set of lines
		if (*point).Value() > r.widget.dataPointScale.Height {
			dp = r.widget.dataPointScale.Height
		} else if (*point).Value() < 0.0 {
			dp = 0.0
		} else {
			dp = (*point).Value()
		}
		yy := yp - (dp * yScale) // using same datasource value
		xx := xp + (float32(idx) * xScale)

		xx = float32(math.Trunc(float64(xx)))
		yy = float32(math.Trunc(float64(yy)))

		thisPoint := fyne.NewPos(xx, yy)
		if idx == 0 {
			lastPoint.Y = yy
		}

		dpv := r.dataPoints[series][idx]
		dpv.Position1 = thisPoint
		dpv.Position2 = lastPoint
		lastPoint = thisPoint

		zt := fyne.NewPos(thisPoint.X-2, thisPoint.Y-2)
		dpm := r.dataPointMarkers[series][idx]
		dpm.Position1 = zt
		zb := fyne.NewPos(thisPoint.X+2, thisPoint.Y+2)
		dpm.Position2 = zb
		(*point).SetMarkerPosition(&zt, &zb)
		if r.widget.enableDataPointMarkers {
			if !dpm.Visible() {
				dpm.Show()
			}
		} else {
			dpm.Hide()
		}
	}
	var found bool
correct:
	for _, o := range r.colorLegend.Objects {
		if o.(*canvas.Text).Text == series {
			found = true
			break correct
		}
	}
	if !found {
		z := canvas.NewText(series, theme.PrimaryColorNamed((*data[0]).ColorName()))
		r.colorLegend.Add(z)
	}

	r.widget.debugLog("lineChartRenderer::layoutSeries() EXIT. Elapsed.microseconds: ", time.Until(startTime).Microseconds())
}

// Layout Given the size required by the fyne application
// move and re-size all custom widget canvas objects here
func (r *lineChartRenderer) Layout(s fyne.Size) {
	r.widget.debugLog("lineChartRenderer::Layout() ENTER: ", s)
	startTime := time.Now()

	r.widget.propertiesLock.Lock()
	defer r.widget.propertiesLock.Unlock()

	r.xInc = (s.Width - (theme.Padding() * 4)) / 16.0
	r.yInc = (s.Height - (theme.Padding() * 3)) / 16.0

	r.xInc = float32(math.Trunc(float64(r.xInc)))
	r.yInc = float32(math.Trunc(float64(r.yInc)))

	// grid Vert lines
	yp := 14.0 * r.yInc
	for idx, line := range r.xLines {
		xp := float32(idx+1) * r.xInc
		line.Position1 = fyne.NewPos(xp+r.xInc, r.yInc) //top
		line.Position2 = fyne.NewPos(xp+r.xInc, yp+8)
	}

	// grid Horiz lines
	xp := r.xInc
	for idx, line := range r.yLines {
		yp := float32(idx+1) * r.yInc
		line.Position1 = fyne.NewPos(xp-8, yp+r.yInc) // left
		line.Position2 = fyne.NewPos(xp*16, yp+r.yInc)
	}

	// grid scale labels
	xp = r.xInc
	yp = 14.0 * r.yInc
	for idx, label := range r.xLabels {
		xxp := float32(idx+1) * r.xInc
		label.Move(fyne.NewPos(xxp+8, yp+10))
	}
	for idx, label := range r.yLabels {
		yyp := float32(idx+1) * r.yInc
		label.Move(fyne.NewPos(xp*0.80, yyp-8))
	}

	// data points
	for key := range r.widget.dataPoints { // datasource
		r.layoutSeries(key)
	}

	r.chartFrame.Resize(fyne.NewSize(r.xInc*15, r.yInc*13))
	r.chartFrame.Move(fyne.NewPos(r.xInc, r.yInc))

	ts := fyne.MeasureText(
		r.topCenteredDesc.Text,
		r.topCenteredDesc.TextSize,
		r.topCenteredDesc.TextStyle)
	r.topCenteredDesc.Move(fyne.Position{X: (s.Width - ts.Width) / 2, Y: theme.Padding() / 4})

	ts = fyne.MeasureText(
		r.topRightDesc.Text,
		r.topRightDesc.TextSize,
		r.topRightDesc.TextStyle)
	r.topRightDesc.Move(fyne.Position{X: (s.Width - ts.Width) - theme.Padding(), Y: ts.Height / 4})
	r.topLeftDesc.Move(fyne.NewPos(theme.Padding(), ts.Height/4))

	msg := strings.Split(r.mouseDisplayContainer.Objects[1].(*widget.Label).Text, "\n")
	ts = fyne.MeasureText(msg[0], 14, r.mouseDisplayContainer.Objects[1].(*widget.Label).TextStyle)
	r.mouseDisplayContainer.Objects[1].(*widget.Label).Resize(fyne.NewSize(ts.Width-theme.Padding(), (2*ts.Height)+(theme.Padding()/2))) // allow room for wrap
	r.mouseDisplayContainer.Objects[0].(*canvas.Rectangle).Resize(fyne.NewSize(ts.Width+theme.Padding(), (2*ts.Height)+theme.Padding()))
	// top edge
	if r.widget.mouseDisplayPosition.Y < theme.Padding()/6 {
		r.widget.mouseDisplayPosition.Y = theme.Padding() / 6
	}
	// left edge
	if r.widget.mouseDisplayPosition.X < theme.Padding()/8 {
		r.widget.mouseDisplayPosition.X = theme.Padding() / 8
	}
	// right edge
	if (r.widget.mouseDisplayPosition.X + ts.Width) > s.Width-(theme.Padding()/4) {
		r.widget.mouseDisplayPosition.X = s.Width - ts.Width - theme.Padding() - (theme.Padding() / 4)
	}
	r.mouseDisplayContainer.Move(*r.widget.mouseDisplayPosition)

	ts = fyne.MeasureText("A", 14, fyne.TextStyle{Bold: true, Monospace: true})
	r.leftMiddleBox.Resize(fyne.NewSize(ts.Width+2, s.Height*0.70))
	r.leftMiddleBox.Move(fyne.NewPos(theme.Padding()/2, s.Height*0.15))

	r.rightMiddleBox.Resize(fyne.NewSize(ts.Width+2, s.Height*0.70))
	r.rightMiddleBox.Move(fyne.NewPos(s.Width-(ts.Width+2), s.Height*0.15))

	ts = fyne.MeasureText(
		r.bottomCenteredDesc.Text,
		r.bottomCenteredDesc.TextSize,
		r.bottomCenteredDesc.TextStyle)
	r.bottomCenteredDesc.Move(fyne.Position{X: (s.Width - ts.Width) / 2, Y: s.Height - ts.Height - theme.Padding()})

	ts = fyne.MeasureText(
		r.bottomRightDesc.Text,
		r.bottomRightDesc.TextSize,
		r.bottomRightDesc.TextStyle)
	r.bottomRightDesc.Move(fyne.NewPos((s.Width-ts.Width)-theme.Padding(), s.Height-ts.Height-theme.Padding()))
	r.bottomLeftDesc.Move(fyne.NewPos(theme.Padding()+2.0, s.Height-ts.Height-theme.Padding()))

	z := r.colorLegend.MinSize()
	r.colorLegend.Move(fyne.NewPos(s.Width-(z.Width+theme.Padding()), (r.yInc*15)+theme.Padding()))

	r.widget.debugLog("lineChartRenderer::Layout() EXIT. Elapsed.microseconds: ", time.Until(startTime).Microseconds())
}

// MinSize Create a minimum size for the widget.
// The smallest size is can be overridden by user
func (r *lineChartRenderer) MinSize() fyne.Size {
	startTime := time.Now()
	r.widget.debugLog("lineChartRenderer::MinSize() ENTER")
	rVal := fyne.NewSize(r.widget.minSize.Width, r.widget.minSize.Height)
	r.widget.debugLog("lineChartRenderer::MinSize() EXIT: renderer: ", rVal, ", Elapsed.microseconds: ", time.Until(startTime).Microseconds())
	return rVal
}

// Objects Return a list of each canvas object.
// but only the objects that have been enabled or are not at default value; i.e. ""
func (r *lineChartRenderer) Objects() []fyne.CanvasObject {
	r.widget.debugLog("lineChartRenderer::Objects() ENTER cnt: ", len(r.widget.objectsCache))
	startTime := time.Now()

	r.widget.propertiesLock.RLock()
	defer r.widget.propertiesLock.RUnlock()

	var objs []fyne.CanvasObject
	objs = append(objs, r.widget.objectsCache...)

	for key, lines := range r.dataPoints {
		for idx, line := range lines {
			objs = append(objs, line)
			marker := r.dataPointMarkers[key][idx]
			objs = append(objs, marker)
		}
	}

	objs = append(objs, r.colorLegend, r.mouseDisplayContainer)

	r.widget.debugLog("lineChartRenderer::Objects() EXIT cnt: ", len(objs), ", Elapsed.microseconds: ", time.Until(startTime).Microseconds())
	return objs
}

// Destroy Cleanup if resources have been allocated
func (r *lineChartRenderer) Destroy() {
	r.widget.debugLog("lineChartRenderer::Destroy() ENTER cnt: ", len(r.widget.objectsCache))
	r.widget.objectsCache = r.widget.objectsCache[:0]
	for key := range r.widget.dataPoints {
		r.widget.dataPoints[key] = r.widget.dataPoints[key][:0]
		r.dataPoints[key] = r.dataPoints[key][:0]
		r.dataPointMarkers[key] = r.dataPointMarkers[key][:0]
	}
	r.widget.debugLog("lineChartRenderer::Destroy() EXIT cnt: ", len(r.widget.objectsCache))
}

// verifyDataPoints Renderer method to inject newly add data series or points
// called by Refresh() to ensure new data is recognized
func (r *lineChartRenderer) verifyDataPoints() {
	startTime := time.Now()

	r.widget.debugLog("lineChartRenderer::VerifyDataPoints() ENTER")
	r.widget.propertiesLock.Lock()
	defer r.widget.propertiesLock.Unlock()

	var changedKeys []string
	var changed bool
	for key, points := range r.widget.dataPoints {
		if nil == r.dataPoints[key] {
			r.dataPoints[key] = []*canvas.Line{}
			r.dataPointMarkers[key] = []*canvas.Circle{}
		}
		changed = false
		for idx, point := range points {
			if idx > (len(r.dataPoints[key]) - 1) { // add added points
				changed = true
				x := canvas.NewLine(theme.PrimaryColorNamed((*point).ColorName()))
				x.StrokeWidth = 2.0
				r.dataPoints[key] = append(r.dataPoints[key], x)
				z := canvas.NewCircle(theme.PrimaryColorNamed((*point).ColorName()))
				z.StrokeWidth = 4.0
				z.Resize(fyne.NewSize(5, 5))
				r.dataPointMarkers[key] = append(r.dataPointMarkers[key], z)
			}
		}
		if changed {
			changedKeys = append(changedKeys, key)
		}
	}
	if len(changedKeys) > 0 {
		for _, series := range changedKeys {
			r.layoutSeries(series)
		}
	}
	r.widget.debugLog("lineChartRenderer::VerifyDataPoints() EXIT. Elapsed.microseconds: ", time.Until(startTime).Microseconds())
}
