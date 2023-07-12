package sknlinechart

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sync"
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

// LineChartSkn widget to display multiple series of data points
// which will roll off older point beyond the 120 point limit.
type LineChartSkn struct {
	widget.BaseWidget       // Inherit from BaseWidget
	desktop.Hoverable       // support mouse tracking
	desktop.Mouseable       // Mouse Clicks
	datapointOrSeriesAdded  bool
	DataPointXLimit         int
	EnableDataPointMarkers  bool
	EnableHorizGridLines    bool
	EnableVertGridLines     bool
	EnableMousePointDisplay bool
	TopLeftLabel            string // The text to display in the widget
	TopCenteredLabel        string
	TopRightLabel           string
	LeftMiddleLabel         string
	RightMiddleLabel        string
	BottomLeftLabel         string
	BottomCenteredLabel     string
	BottomRightLabel        string
	mouseDisplayStr         string
	mouseDisplayPosition    *fyne.Position
	mouseDisplayFrameColor  string
	dataPoints              *map[string][]LineChartDatapoint
	dataPointScale          fyne.Size
	minSize                 fyne.Size
	objects                 []fyne.CanvasObject
	propertyLock            sync.RWMutex
}

var _ LineChart = (*LineChartSkn)(nil)

// NewLineChart Create the Line Chart
// be careful not to exceed the series data point limit, which defaults to 120
//
// can return a valid chart object and an error object; errors really should be handled
// and are caused by data points exceeding the container limit of 120; they will be truncated
func NewLineChart(topTitle, bottomTitle string, dataPoints *map[string][]LineChartDatapoint) (*LineChartSkn, error) {
	fmt.Println("::NewLineChart()")
	var err error
	dpl := 120
	for key, points := range *dataPoints {
		if len(points) > dpl {
			err = fmt.Errorf("%s\nNewLineChart() dataPoint contents exceeds the point count limit[Action: truncated leading]. Series: %s, points: %d, Limit: %d", err, key, len(points), dpl)
			for len(points) > dpl {
				points = RemoveIndexFromSlice(0, points)
			}
			(*dataPoints)[key] = points
			err = fmt.Errorf("%s\n NewLineChart() dataPoint contents exceeds the point count limit[Action: truncated leading]. Series: %s, points: %d, Limit: %d", err.Error(), key, len(points), dpl)
		}
	}

	w := &LineChartSkn{ // Create this widget with an initial text value
		dataPoints:              dataPoints,
		datapointOrSeriesAdded:  true,
		DataPointXLimit:         dpl,
		dataPointScale:          fyne.NewSize(float32(dpl), 110.0),
		EnableDataPointMarkers:  true,
		EnableHorizGridLines:    true,
		EnableVertGridLines:     true,
		EnableMousePointDisplay: true,
		mouseDisplayStr:         "",
		mouseDisplayPosition:    &fyne.Position{},
		mouseDisplayFrameColor:  string(theme.ColorNameForeground),
		TopLeftLabel:            "top left desc",
		TopCenteredLabel:        topTitle,
		TopRightLabel:           "",
		LeftMiddleLabel:         "left middle desc",
		RightMiddleLabel:        "right middle desc",
		BottomLeftLabel:         "bottom left desc",
		BottomCenteredLabel:     bottomTitle,
		BottomRightLabel:        "bottom right desc",
		minSize:                 fyne.NewSize(420+theme.Padding()*4, 315+theme.Padding()*4),
		objects:                 []fyne.CanvasObject{}, // everything except datapoints, markers, and mousebox
		propertyLock:            sync.RWMutex{},
	}
	w.ExtendBaseWidget(w) // Initialize the BaseWidget
	return w, err
}

// CreateRenderer Create the renderer. This is called by the fyne application
func (w *LineChartSkn) CreateRenderer() fyne.WidgetRenderer {
	fmt.Println("LineChartSkn::CreateRenderer()")
	return newLineChartRenderer(w)
}

// SetMinSize override the default min size of chart
func (w *LineChartSkn) SetMinSize(s fyne.Size) {
	w.minSize = s
}

// GetTopLeftLabel return text from top left label
func (w *LineChartSkn) GetTopLeftLabel() string {
	return w.TopLeftLabel
}

// GetTitle return text of the chart's title from top center
func (w *LineChartSkn) GetTitle() string {
	return w.TopCenteredLabel
}

// IsDataPointMarkersEnabled returns state of chart's use of data point markers on series data
func (w *LineChartSkn) IsDataPointMarkersEnabled() bool {
	return w.EnableDataPointMarkers
}

// IsHorizGridLinesEnabled returns state of chart's display of horizontal grid line
func (w *LineChartSkn) IsHorizGridLinesEnabled() bool {
	return w.EnableHorizGridLines
}

// IsVertGridLinesEnabled returns state of chart's display of vertical grid line
func (w *LineChartSkn) IsVertGridLinesEnabled() bool {
	return w.EnableVertGridLines
}

// IsMousePointDisplayEnabled return state of mouse popups when hovered over a chart datapoint
func (w *LineChartSkn) IsMousePointDisplayEnabled() bool {
	return w.EnableMousePointDisplay
}

// GetTopRightLabel returns text of top right label
func (w *LineChartSkn) GetTopRightLabel() string {
	return w.TopRightLabel
}

// GetMiddleLeftLabel returns text of middle left label
func (w *LineChartSkn) GetMiddleLeftLabel() string {
	return w.LeftMiddleLabel
}

// GetMiddleRightLabel returns text of middle right label
func (w *LineChartSkn) GetMiddleRightLabel() string {
	return w.RightMiddleLabel
}

// GetBottomLeftLabel returns text of bottom left label
func (w *LineChartSkn) GetBottomLeftLabel() string {
	return w.BottomLeftLabel
}

// GetBottomCenteredLabel returns text of bottom center label
func (w *LineChartSkn) GetBottomCenteredLabel() string {
	return w.BottomCenteredLabel
}

// GetBottomRightLabel returns text of bottom right label
func (w *LineChartSkn) GetBottomRightLabel() string {
	return w.BottomRightLabel
}

// SetTopLeftLabel sets text to be display on chart at top left
func (w *LineChartSkn) SetTopLeftLabel(newValue string) {
	w.TopLeftLabel = newValue
}

// SetTitle sets text to be display on chart at top center
func (w *LineChartSkn) SetTitle(newValue string) {
	w.TopCenteredLabel = newValue
}

// SetTopRightLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetTopRightLabel(newValue string) {
	w.TopRightLabel = newValue
}

// SetMiddleLeftLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetMiddleLeftLabel(newValue string) {
	w.LeftMiddleLabel = newValue
}

// SetMiddleRightLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetMiddleRightLabel(newValue string) {
	w.RightMiddleLabel = newValue
}

// SetBottomLeftLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetBottomLeftLabel(newValue string) {
	w.BottomLeftLabel = newValue
}

// SetBottomRightLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetBottomRightLabel(newValue string) {
	w.BottomRightLabel = newValue
}

// SetBottomCenteredLabel changes displayed text, empty disables display
func (w *LineChartSkn) SetBottomCenteredLabel(newValue string) {
	w.BottomCenteredLabel = newValue
}

// SetDataPointMarkers enables data point markers on display series points
func (w *LineChartSkn) SetDataPointMarkers(enable bool) {
	w.EnableDataPointMarkers = enable
}

// SetHorizGridLines enables chart horizontal grid lines
func (w *LineChartSkn) SetHorizGridLines(enable bool) {
	w.EnableHorizGridLines = enable
}

// SetVertGridLines enables chart vertical grid lines
func (w *LineChartSkn) SetVertGridLines(enable bool) {
	w.EnableVertGridLines = enable
}

// SetMousePointDisplay true/false, enables data point display under mouse pointer
func (w *LineChartSkn) SetMousePointDisplay(enable bool) {
	w.EnableMousePointDisplay = enable
}

// ApplyDataSeries adds a new series of data to existing chart set.
// throws error if new series exceeds containers point limit
func (w *LineChartSkn) ApplyDataSeries(seriesName string, newSeries []LineChartDatapoint) error {
	fmt.Println("LineChartSkn::ApplyDataSeries() ENTER")
	if w == nil {
		fmt.Println("LineChartSkn::ApplyDataSeries() ERROR EXIT")
		return fmt.Errorf("ApplyDataSeries() no active widget")
	}

	if len(newSeries) < w.DataPointXLimit {
		w.propertyLock.Lock()
		(*w.dataPoints)[seriesName] = newSeries
		w.datapointOrSeriesAdded = true
		w.propertyLock.Unlock()
		w.Refresh()
	} else {
		fmt.Println("LineChartSkn::ApplyDataSeries() ERROR EXIT")
		return fmt.Errorf("[%s] data series datapoints limit exceeded. limit:%d, count:%d", seriesName, w.DataPointXLimit, len(newSeries))
	}
	fmt.Println("LineChartSkn::ApplyDataSeries() EXIT")
	return nil
}

// ApplyDataPoint adds a new datapoint to an existing series
// will shift out the oldest point if containers limit is exceeded
func (w *LineChartSkn) ApplyDataPoint(seriesName string, newDataPoint LineChartDatapoint) {
	fmt.Println("LineChartSkn::ApplyDataPoint() ENTER")
	if w == nil {
		return
	}

	w.propertyLock.Lock()

	if len((*w.dataPoints)[seriesName]) < w.DataPointXLimit {
		(*w.dataPoints)[seriesName] = append((*w.dataPoints)[seriesName], newDataPoint)
	} else {
		(*w.dataPoints)[seriesName] = ShiftSlice(newDataPoint, (*w.dataPoints)[seriesName])
	}
	w.datapointOrSeriesAdded = true
	w.propertyLock.Unlock()
	w.Refresh()
	fmt.Println("LineChartSkn::ApplyDataPoint() EXIT")
}

// MinSize Create a minimum size for the widget.
// The smallest size is can be overridden by user
// also in renderer
func (w *LineChartSkn) MinSize() fyne.Size {
	fmt.Println("LineChartSkn::MinSize() ENTER")
	w.ExtendBaseWidget(w)
	val := w.BaseWidget.MinSize()
	fmt.Println("LineChartSkn::MinSize() EXIT: ", val, "Current: ", w.Size())
	return val
}

// Refresh triggers a redraw
func (w *LineChartSkn) Refresh() {
	fmt.Println("LineChartSkn::Refresh() ENTER")
	w.ExtendBaseWidget(w)
	w.BaseWidget.Refresh()
	fmt.Println("LineChartSkn::Refresh() EXIT")
}

// Resize sets a new size for the label.
// This should only be called if it is not in a container with a layout manager.
func (w *LineChartSkn) Resize(s fyne.Size) {
	fmt.Println("LineChartSkn::Resize() ENTER")
	w.BaseWidget.Resize(s)
	fmt.Println("LineChartSkn::Resize() EXIT")
}

// MouseDown btn.2 toggles markers, btn.1 toggles mouse point display
func (w *LineChartSkn) MouseDown(me *desktop.MouseEvent) {
	fmt.Println("LineChartSkn::MouseDown() ENTER")
	if me.Button == desktop.MouseButtonSecondary {
		w.EnableDataPointMarkers = !w.EnableDataPointMarkers
		w.Refresh()
	} else if me.Button == desktop.MouseButtonPrimary {
		w.EnableMousePointDisplay = !w.EnableMousePointDisplay
		w.Refresh()
	}
	fmt.Println("LineChartSkn::MouseDown() EXIT")
}

// MouseUp unused interface method
func (w *LineChartSkn) MouseUp(*desktop.MouseEvent) {
	fmt.Println("LineChartSkn::MouseUP()")
}

// MouseIn unused interface method
func (w *LineChartSkn) MouseIn(*desktop.MouseEvent) {
	fmt.Println("LineChartSkn::MouseIn()")
}

// MouseMoved interface method to discover which data point is under mouse
func (w *LineChartSkn) MouseMoved(me *desktop.MouseEvent) {
	fmt.Println("LineChartSkn::MouseMoved()")
	if !w.EnableMousePointDisplay {
		return
	}
	for key, points := range *w.dataPoints {
		for idx, point := range points {
			top, bottom := point.MarkerPosition()
			if !me.Position.IsZero() && !top.IsZero() {
				if me.Position.X >= top.X && me.Position.X <= bottom.X &&
					me.Position.Y >= top.Y && me.Position.Y <= bottom.Y {
					value := fmt.Sprint(" Series: ", key, ", Index: ", idx, ", Value: ", point.Value(), " [ ", point.Timestamp(), " ]")
					w.enableMouseContainer(value, point.ColorName(), &me.Position).Refresh()
				}
			}
		}
	}
}

// MouseOut disable display of mouse data point display
func (w *LineChartSkn) MouseOut() {
	fmt.Println("LineChartSkn::MouseOut()")
	w.disableMouseContainer()
}

// enableMouseContainer private method to prepare values need by renderer to create pop display
// composes display text, captures position and colorName for use by renderer
func (w *LineChartSkn) enableMouseContainer(value, frameColor string, mousePosition *fyne.Position) *LineChartSkn {
	fmt.Println("LineChartSkn::enableMouseContainer() ENTER")
	w.mouseDisplayStr = value
	w.mouseDisplayFrameColor = frameColor
	ct := canvas.NewText(value, theme.PrimaryColorNamed(frameColor))
	parts := strings.Split(value, "[")
	ts := fyne.MeasureText(parts[0], ct.TextSize, ct.TextStyle)
	mp := &fyne.Position{X: mousePosition.X - (ts.Width / 2), Y: mousePosition.Y - (3 * ts.Height) - theme.Padding()}
	w.mouseDisplayPosition = mp
	fmt.Println("LineChartSkn::enableMouseContainer() EXIT")
	return w
}

// disableMouseContainer private method to manage mouse leaving window
// blank string will prevent display
func (w *LineChartSkn) disableMouseContainer() {
	fmt.Println("LineChartSkn::disableMouseContainer()")
	w.mouseDisplayStr = ""
	w.Refresh()
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
}

var _ fyne.WidgetRenderer = (*lineChartRenderer)(nil)

// Create the renderer with a reference to the widget
// and all the objects to be displayed for this custom widget
//
// Note: Do not size or move canvas objects here.
func newLineChartRenderer(lineChart *LineChartSkn) *lineChartRenderer {
	fmt.Println("::newLineChartRenderer() ENTER")
	lineChart.ExtendBaseWidget(lineChart)
	lineChart.propertyLock.Lock()
	defer lineChart.propertyLock.Unlock()

	objs := []fyne.CanvasObject{}

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

	dataPoints := map[string][]*canvas.Line{}
	dpMaker := map[string][]*canvas.Circle{}
	var xlines, ylines []*canvas.Line
	var xLabels, yLabels []*canvas.Text

	for i := 0; i < 11; i++ {
		x := canvas.NewLine(theme.PrimaryColorNamed(theme.ColorGreen))
		x.StrokeWidth = 0.25
		y := canvas.NewLine(theme.PrimaryColorNamed(theme.ColorGreen))
		y.StrokeWidth = 0.25
		xlines = append(xlines, x)
		ylines = append(ylines, y)
		objs = append(objs, x, y)
	}

	for i := 0; i < 12; i++ {
		yt := strconv.Itoa((11 - i) * 10)
		yl := canvas.NewText(yt, theme.ForegroundColor())
		yl.Alignment = fyne.TextAlignTrailing
		yLabels = append(yLabels, yl)
		objs = append(objs, yl)
	}
	for i := 0; i < 13; i++ {
		xt := strconv.Itoa(i * 10)
		xl := canvas.NewText(xt, theme.ForegroundColor())
		xl.Alignment = fyne.TextAlignTrailing
		xLabels = append(xLabels, xl)
		objs = append(objs, xl)
	}

	for key, points := range *lineChart.dataPoints {
		for _, point := range points {
			x := canvas.NewLine(theme.PrimaryColorNamed(point.ColorName()))
			x.StrokeWidth = 2.0
			dataPoints[key] = append(dataPoints[key], x)
			z := canvas.NewCircle(theme.PrimaryColorNamed(point.ColorName()))
			z.StrokeWidth = 4.0
			dpMaker[key] = append(dpMaker[key], z)
		}
	}

	topCenteredDesc := canvas.NewText(lineChart.TopCenteredLabel, theme.ForegroundColor())
	topCenteredDesc.TextSize = 24
	topCenteredDesc.TextStyle = fyne.TextStyle{
		Bold:   true,
		Italic: false,
	}
	objs = append(objs, topCenteredDesc)

	bottomCenteredDesc := canvas.NewText(lineChart.BottomCenteredLabel, theme.ForegroundColor())
	bottomCenteredDesc.TextSize = 16
	bottomCenteredDesc.TextStyle = fyne.TextStyle{
		Bold:   false,
		Italic: true,
	}
	objs = append(objs, bottomCenteredDesc)

	// vertical text for X/Y legends since no text rotation is available
	lBox := container.NewVBox()
	for _, c := range lineChart.LeftMiddleLabel {
		lBox.Add(
			canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground))),
		)
	}
	objs = append(objs, lBox)

	rBox := container.NewVBox()
	for _, c := range lineChart.RightMiddleLabel {
		rBox.Add(
			canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground))),
		)
	}
	objs = append(objs, rBox)

	tl := canvas.NewText(lineChart.TopLeftLabel, theme.ForegroundColor())
	tr := canvas.NewText(lineChart.TopRightLabel, theme.ForegroundColor())
	bl := canvas.NewText(lineChart.BottomLeftLabel, theme.ForegroundColor())
	br := canvas.NewText(lineChart.BottomRightLabel, theme.ForegroundColor())
	objs = append(objs, tl, tr, bl, br)

	// save all except data points, markers, and mouse box
	lineChart.objects = append(lineChart.objects, objs...)

	fmt.Println("::newLineChartRenderer() EXIT")

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
	}
}

// Refresh method is called if the state of the widget changes or the
// theme is changed
func (r *lineChartRenderer) Refresh() {
	fmt.Println("lineChartRenderer::Refresh() ENTER")
	if r.widget.datapointOrSeriesAdded {
		r.verifyDataPoints()
	}

	r.mouseDisplayContainer.Objects[0].(*canvas.Rectangle).StrokeColor = theme.PrimaryColorNamed(r.widget.mouseDisplayFrameColor)
	r.mouseDisplayContainer.Objects[1].(*widget.Label).SetText(r.widget.mouseDisplayStr)

	r.topLeftDesc.Text = r.widget.TopLeftLabel
	r.topCenteredDesc.Text = r.widget.TopCenteredLabel
	r.topRightDesc.Text = r.widget.TopRightLabel
	r.bottomLeftDesc.Text = r.widget.BottomLeftLabel
	r.bottomCenteredDesc.Text = r.widget.BottomCenteredLabel
	r.bottomRightDesc.Text = r.widget.BottomRightLabel

	r.leftMiddleBox.RemoveAll()
	for _, c := range r.widget.LeftMiddleLabel {
		z := canvas.NewText(
			strings.ToUpper(string(c)),
			theme.PrimaryColorNamed(string(theme.ColorNameForeground)))
		z.TextSize = 12
		r.leftMiddleBox.Add(z)
	}
	r.rightMiddleBox.RemoveAll()
	for _, c := range r.widget.RightMiddleLabel {
		z := canvas.NewText(
			strings.ToUpper(string(c)),
			theme.PrimaryColorNamed(string(theme.ColorNameForeground)))
		z.TextSize = 12
		r.rightMiddleBox.Add(z)
	}

	if r.widget.datapointOrSeriesAdded {
		r.Layout(r.widget.Size())
	}

	r.widget.propertyLock.RLock()
	defer r.widget.propertyLock.RUnlock()

	for _, v := range r.widget.objects {
		v.Refresh()
	}

	if r.widget.datapointOrSeriesAdded {
		for _, points := range r.dataPoints {
			for _, point := range points {
				point.Refresh()
			}
		}
		for _, markers := range r.dataPointMarkers {
			for _, mark := range markers {
				mark.Refresh()
			}
		}
		r.widget.datapointOrSeriesAdded = false
	}
	r.mouseDisplayContainer.Refresh()
	fmt.Println("lineChartRenderer::Refresh() EXIT")
}

// Layout Given the size required by the fyne application
// move and re-size the all custom widget canvas objects here
func (r *lineChartRenderer) Layout(s fyne.Size) {
	fmt.Println("lineChartRenderer::Layout() ENTER: ", s)
	r.widget.propertyLock.Lock()
	defer r.widget.propertyLock.Unlock()

	r.xInc = (s.Width - theme.Padding()) / 14.0
	r.yInc = (s.Height - theme.Padding()) / 14.0

	// grid Vert lines
	yp := 12.0 * r.yInc
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
		line.Position2 = fyne.NewPos(xp*13, yp+r.yInc)
	}

	// grid scale labels
	xp = r.xInc
	yp = 12.0 * r.yInc
	for idx, label := range r.xLabels {
		xxp := float32(idx+1) * r.xInc
		label.Move(fyne.NewPos(xxp+8, yp+10))
	}
	for idx, label := range r.yLabels {
		yyp := float32(idx+1) * r.yInc
		label.Move(fyne.NewPos(xp*0.80, yyp-8))
	}

	// data points
	xp = r.xInc
	yp = r.yInc * 12
	yScale := (r.yInc * 10) / 100
	xScale := (r.xInc * 10) / 100
	dp := float32(1.0)
	for key, data := range *r.widget.dataPoints { // datasource
		lastPoint := fyne.NewPos(xp, yp)

		for idx, point := range data { // one set of lines
			if point.Value() > r.widget.dataPointScale.Height {
				dp = r.widget.dataPointScale.Height
			} else if point.Value() < 0.0 {
				dp = 0.0
			} else {
				dp = point.Value()
			}
			yy := yp - (dp * yScale) // using same datasource value
			xx := xp + (float32(idx) * xScale)
			thisPoint := fyne.NewPos(xx, yy)
			if idx == 0 {
				lastPoint.Y = yy
			}

			dpv := r.dataPoints[key][idx]
			dpv.Position1 = thisPoint
			dpv.Position2 = lastPoint
			lastPoint = thisPoint

			zt := fyne.NewPos(thisPoint.X-3, thisPoint.Y-3)
			dpm := r.dataPointMarkers[key][idx]
			dpm.Position1 = zt
			zb := fyne.NewPos(thisPoint.X+3, thisPoint.Y+3)
			dpm.Position2 = zb
			point.SetMarkerPosition(&zt, &zb)
			dpm.Resize(fyne.NewSize(5, 5))
		}
	}

	r.chartFrame.Resize(fyne.NewSize(r.xInc*12, r.yInc*11))
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

	msg := strings.Split(r.mouseDisplayContainer.Objects[1].(*widget.Label).Text, " [ ")
	ts = fyne.MeasureText(msg[0], 14, r.mouseDisplayContainer.Objects[1].(*widget.Label).TextStyle)

	r.mouseDisplayContainer.Objects[1].(*widget.Label).Resize(fyne.NewSize(ts.Width-theme.Padding(), (2*ts.Height)+theme.Padding())) // allow room for wrap
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

	r.leftMiddleBox.Resize(fyne.NewSize(ts.Height, s.Height*0.75))
	r.leftMiddleBox.Move(fyne.NewPos(2*theme.Padding(), s.Height*0.1))

	r.rightMiddleBox.Resize(fyne.NewSize(ts.Height, s.Height*0.75))
	r.rightMiddleBox.Move(fyne.NewPos(s.Width-ts.Height-(2*theme.Padding()), s.Height*0.1))

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

	fmt.Println("lineChartRenderer::Layout() EXIT")
}

// MinSize Create a minimum size for the widget.
// The smallest size is can be overridden by user
func (r *lineChartRenderer) MinSize() fyne.Size {
	fmt.Println("lineChartRenderer::MinSize() ENTER")
	r.widget.ExtendBaseWidget(r.widget)
	val := fyne.NewSize(r.widget.minSize.Width, r.widget.minSize.Height)
	fmt.Println("lineChartRenderer::MinSize() EXIT: ", val)
	return val
}

// Objects Return a list of each canvas object.
// but only the objects that have been enabled or are not at default value; i.e. ""
func (r *lineChartRenderer) Objects() []fyne.CanvasObject {
	fmt.Println("lineChartRenderer::Objects() ENTER")
	r.widget.propertyLock.Lock()
	defer r.widget.propertyLock.Unlock()

	objs := []fyne.CanvasObject{}
	objs = append(objs, r.widget.objects...)

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
	if r.widget.LeftMiddleLabel != "" {
		if !r.leftMiddleBox.Visible() {
			r.leftMiddleBox.Show()
		}
	} else {
		r.leftMiddleBox.Hide()
	}
	if r.widget.RightMiddleLabel != "" {
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

	for _, line := range r.xLines {
		if r.widget.EnableHorizGridLines {
			if !line.Visible() {
				line.Show()
			}
		} else {
			line.Hide()
		}
	}
	for _, line := range r.yLines {
		if r.widget.EnableVertGridLines {
			if !line.Visible() {
				line.Show()
			}
		} else {
			line.Hide()
		}
	}

	for key, lines := range r.dataPoints {
		for idx, line := range lines {
			objs = append(objs, line)
			marker := r.dataPointMarkers[key][idx]
			objs = append(objs, marker)
			if r.widget.EnableDataPointMarkers {
				if !marker.Visible() {
					marker.Show()
				}
			} else {
				marker.Hide()
			}
		}
	}

	objs = append(objs, r.mouseDisplayContainer)
	if r.widget.EnableMousePointDisplay {
		if r.mouseDisplayContainer.Objects[1].(*widget.Label).Text != "" {
			if !r.mouseDisplayContainer.Visible() {
				r.mouseDisplayContainer.Show()
			}
		} else {
			r.mouseDisplayContainer.Hide()
		}
	} else {
		r.mouseDisplayContainer.Hide()
	}
	fmt.Println("lineChartRenderer::Objects() EXIT")
	return objs
}

// Destroy Cleanup if resources have been allocated
func (r *lineChartRenderer) Destroy() {
	fmt.Println("lineChartRenderer::Destroy()")
}

// verifyDataPoints Renderer method to inject newly add data series or points
// called by Refresh() to ensure new data is recognized
func (r *lineChartRenderer) verifyDataPoints() {
	fmt.Println("lineChartRenderer::VerifyDataPoints() ENTER")
	r.widget.propertyLock.Lock()
	defer r.widget.propertyLock.Unlock()

	for key, points := range *r.widget.dataPoints {
		if nil == r.dataPoints[key] {
			r.dataPoints[key] = []*canvas.Line{}
			r.dataPointMarkers[key] = []*canvas.Circle{}
		}
		for idx, point := range points {
			if idx > (len(r.dataPoints[key]) - 1) { // add added points
				x := canvas.NewLine(theme.PrimaryColorNamed(point.ColorName()))
				x.StrokeWidth = 2.0
				r.dataPoints[key] = append(r.dataPoints[key], x)
				z := canvas.NewCircle(theme.PrimaryColorNamed(point.ColorName()))
				z.StrokeWidth = 4.0
				r.dataPointMarkers[key] = append(r.dataPointMarkers[key], z)
			}
		}
	}
	fmt.Println("lineChartRenderer::VerifyDataPoints() EXIT")
}
