package components

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/sknlinechart/pkg/commons"
	"image/color"
	"strconv"
	"strings"
)

// SknLineChart feature list
type SknLineChart interface {
	IsDataPointMarkersEnabled() bool
	IsHorizGridLinesEnabled() bool
	IsVertGridLinesEnabled() bool
	IsMousePointDisplayEnabled() bool

	EnableDataPointMarkers(newValue bool)
	EnabledHorizGridLines(newValue bool)
	EnableVertGridLine(newValue bool)
	EnableMousePointDisplay(newValue bool)

	GetTopLeftDescription() string
	GetTitle() string
	GetTopRightDescription() string
	GetMiddleLeftDescription() string
	GetMiddleRightDescription() string
	GetBottomLeftDescription() string
	GetBottomCenteredDescription() string
	GetBottomRightDescription() string
	SetTopLeftDescription(newValue string)

	SetTitle(newValue string)
	SetTopRightDescription(newValue string)
	SetMiddleLeftDescription(newValue string)
	SetMiddleRightDescription(newValue string)
	SetBottomLeftDescription(newValue string)
	SetBottomRightDescription(newValue string)
	SetBottomCenteredDescription(newValue string)

	SetMinSize(s fyne.Size)
	ReplaceDataSeries(newData *map[string][]SknChartDatapoint) error

	ApplyNewDataSeries(seriesName string, newSeries []SknChartDatapoint) error
	ApplySingleDataPoint(seriesName string, newDataPoint SknChartDatapoint)
}

// LineChartSkn widget to display multiple series of data points
// which will roll off older point beyond the 120 point limit.
type LineChartSkn struct {
	widget.BaseWidget       // Inherit from BaseWidget
	desktop.Hoverable       // support mouse tracking
	desktop.Mouseable       // Mouse Clicks
	dataPoints              *map[string][]SknChartDatapoint
	dataPointScale          fyne.Size
	dataPointLimit          int
	enableDataPointMarkers  bool
	enableHorizGridLines    bool
	enableVertGridLines     bool
	enableMousePointDisplay bool
	mouseDisplayStr         string
	mouseDisplayPosition    *fyne.Position
	mouseDisplayFrameColor  string
	topLeftDesc             string // The text to display in the widget
	topCenteredDesc         string
	topRightDesc            string
	leftMiddleDesc          string
	rightMiddleDesc         string
	bottomCenteredDesc      string
	bottomLeftDesc          string
	bottomRightDesc         string
	minSize                 fyne.Size
}

var _ SknLineChart = (*LineChartSkn)(nil)

// NewSknLineChart Create the Line Chart
// be careful not to exceed the series data point limit, which defaults to 120
//
// can return a valid chart object and an error object; errors really should be handled
// and are caused by data points exceeding the container limit of 120; they will be truncated
func NewSknLineChart(tTitle, bTitle string, dataPoints *map[string][]SknChartDatapoint) (*LineChartSkn, error) {
	var err error
	dpl := 120
	for key, points := range *dataPoints {
		if len(points) > dpl {
			err = fmt.Errorf("%s\nNewSknLineChart() dataPoint contents exceeds the point count limit[Action: truncated leading]. Series: %s, points: %d, Limit: %d", err, key, len(points), dpl)
			for len(points) > dpl {
				points = commons.RemoveIndexFromSlice(0, points)
			}
			(*dataPoints)[key] = points
			err = fmt.Errorf("%s\n NewSknLineChart() dataPoint contents exceeds the point count limit[Action: truncated leading]. Series: %s, points: %d, Limit: %d", err.Error(), key, len(points), dpl)
		}
	}

	w := &LineChartSkn{ // Create this widget with an initial text value
		dataPoints:              dataPoints,
		dataPointLimit:          dpl,
		dataPointScale:          fyne.NewSize(120.0, 110.0),
		enableDataPointMarkers:  true,
		enableHorizGridLines:    true,
		enableVertGridLines:     true,
		enableMousePointDisplay: true,
		mouseDisplayStr:         "",
		mouseDisplayPosition:    &fyne.Position{},
		mouseDisplayFrameColor:  string(theme.ColorNameForeground),
		topLeftDesc:             "top left desc",
		topCenteredDesc:         tTitle,
		topRightDesc:            "",
		leftMiddleDesc:          "left middle desc",
		rightMiddleDesc:         "right middle desc",
		bottomLeftDesc:          "bottom left desc",
		bottomCenteredDesc:      bTitle,
		bottomRightDesc:         "bottom right desc",
		minSize:                 fyne.NewSize(400+theme.Padding()*4, 300+theme.Padding()*4),
	}
	w.ExtendBaseWidget(w) // Initialize the BaseWidget
	return w, err
}

// CreateRenderer Create the renderer. This is called by the fyne application
func (w *LineChartSkn) CreateRenderer() fyne.WidgetRenderer {
	return newSknLineChartRenderer(w)
}

// SetMinSize override the default min size of chart
func (w *LineChartSkn) SetMinSize(s fyne.Size) {
	w.minSize = s
}

// GetTopLeftDescription return text from top left label
func (w *LineChartSkn) GetTopLeftDescription() string {
	return w.topLeftDesc
}

// GetTitle return text of the chart's title from top center
func (w *LineChartSkn) GetTitle() string {
	return w.topCenteredDesc
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

// IsMousePointDisplayEnabled return state of mouse popups when hovered over a chart datapoint
func (w *LineChartSkn) IsMousePointDisplayEnabled() bool {
	return w.enableMousePointDisplay
}

// GetTopRightDescription returns text of top right label
func (w *LineChartSkn) GetTopRightDescription() string {
	return w.topRightDesc
}

// GetMiddleLeftDescription returns text of middle left label
func (w *LineChartSkn) GetMiddleLeftDescription() string {
	return w.leftMiddleDesc
}

// GetMiddleRightDescription returns text of middle right label
func (w *LineChartSkn) GetMiddleRightDescription() string {
	return w.rightMiddleDesc
}

// GetBottomLeftDescription returns text of bottom left label
func (w *LineChartSkn) GetBottomLeftDescription() string {
	return w.bottomLeftDesc
}

// GetBottomCenteredDescription returns text of bottom center label
func (w *LineChartSkn) GetBottomCenteredDescription() string {
	return w.bottomCenteredDesc
}

// GetBottomRightDescription returns text of bottom right label
func (w *LineChartSkn) GetBottomRightDescription() string {
	return w.bottomRightDesc
}

// SetTopLeftDescription sets text to be display on chart at top left
func (w *LineChartSkn) SetTopLeftDescription(newValue string) {
	w.topLeftDesc = newValue
}

// SetTitle sets text to be display on chart at top center
func (w *LineChartSkn) SetTitle(newValue string) {
	w.topCenteredDesc = newValue
}

// SetTopRightDescription changes displayed text, empty disables display
func (w *LineChartSkn) SetTopRightDescription(newValue string) {
	w.topRightDesc = newValue
}

// SetMiddleLeftDescription changes displayed text, empty disables display
func (w *LineChartSkn) SetMiddleLeftDescription(newValue string) {
	w.leftMiddleDesc = newValue
}

// SetMiddleRightDescription changes displayed text, empty disables display
func (w *LineChartSkn) SetMiddleRightDescription(newValue string) {
	w.rightMiddleDesc = newValue
}

// SetBottomLeftDescription changes displayed text, empty disables display
func (w *LineChartSkn) SetBottomLeftDescription(newValue string) {
	w.bottomLeftDesc = newValue
}

// SetBottomRightDescription changes displayed text, empty disables display
func (w *LineChartSkn) SetBottomRightDescription(newValue string) {
	w.bottomRightDesc = newValue
}

// SetBottomCenteredDescription changes displayed text, empty disables display
func (w *LineChartSkn) SetBottomCenteredDescription(newValue string) {
	w.bottomCenteredDesc = newValue
}

// EnableDataPointMarkers enables data point markers on display series points
func (w *LineChartSkn) EnableDataPointMarkers(newValue bool) {
	w.enableDataPointMarkers = newValue
}

// EnabledHorizGridLines enables chart horizontal grid lines
func (w *LineChartSkn) EnabledHorizGridLines(newValue bool) {
	w.enableHorizGridLines = newValue
}

// EnableVertGridLine enables chart vertical grid lines
func (w *LineChartSkn) EnableVertGridLine(newValue bool) {
	w.enableVertGridLines = newValue
}

// EnableMousePointDisplay true/false, enables data point display under mouse pointer
func (w *LineChartSkn) EnableMousePointDisplay(newValue bool) {
	w.enableMousePointDisplay = newValue
}

// ReplaceDataSeries replaces all the chart's series data
// throws an error if any series exceeds the container limit
func (w *LineChartSkn) ReplaceDataSeries(newData *map[string][]SknChartDatapoint) error {
	if w.dataPoints == nil {
		return fmt.Errorf("ReplaceDataSeries() no active widget")
	}
	if len(*w.dataPoints) <= len(*newData) {
		for key, points := range *newData {
			if len(points) > w.dataPointLimit {
				return fmt.Errorf("[%s] data series datapoints limit exceeded. limit:%d, count:%d", key, w.dataPointLimit, len(points))
			}
		}
		w.dataPoints = newData
		w.Refresh()
	} else {
		return fmt.Errorf("newData must be larger[%d] than or equal to existing[%d]", len(*newData), len(*w.dataPoints))
	}
	return nil
}

// ApplyNewDataSeries adds a new series of data to existing chart set.
// throws error if new series exceeds containers point limit
func (w *LineChartSkn) ApplyNewDataSeries(seriesName string, newSeries []SknChartDatapoint) error {
	if w == nil {
		return fmt.Errorf("ApplyNewDataSeries() no active widget")
	}

	if len(newSeries) < w.dataPointLimit {
		(*w.dataPoints)[seriesName] = newSeries
		w.Refresh()
	} else {
		return fmt.Errorf("[%s] data series datapoints limit exceeded. limit:%d, count:%d", seriesName, w.dataPointLimit, len(newSeries))
	}
	return nil
}

// ApplySingleDataPoint adds a new datapoint to an existing series
// will shift out the oldest point if containers limit is exceeded
func (w *LineChartSkn) ApplySingleDataPoint(seriesName string, newDataPoint SknChartDatapoint) {
	if w == nil {
		return
	}
	if len((*w.dataPoints)[seriesName]) < w.dataPointLimit {
		(*w.dataPoints)[seriesName] = append((*w.dataPoints)[seriesName], newDataPoint)
	} else {
		(*w.dataPoints)[seriesName] = commons.ShiftSlice(newDataPoint, (*w.dataPoints)[seriesName])
	}
	w.Refresh()
}

func (w *LineChartSkn) MouseDown(me *desktop.MouseEvent) {
	if me.Button == desktop.MouseButtonSecondary {
		w.enableDataPointMarkers = !w.enableDataPointMarkers
		w.Refresh()
	} else if me.Button == desktop.MouseButtonPrimary {
		w.disableMouseContainer()
	}
}
func (w *LineChartSkn) MouseUp(*desktop.MouseEvent) {}

// MouseIn unused interface method
func (w *LineChartSkn) MouseIn(*desktop.MouseEvent) {}

// MouseMoved interface method to discover which data point is under mouse
func (w *LineChartSkn) MouseMoved(me *desktop.MouseEvent) {
	if !w.enableMousePointDisplay {
		return
	}
	for key, points := range *w.dataPoints {
		for idx, point := range points {
			top, bottom := point.MarkerPosition()
			if !me.Position.IsZero() && !top.IsZero() {
				if me.Position.X >= top.X && me.Position.X <= bottom.X &&
					me.Position.Y >= top.Y && me.Position.Y <= bottom.Y {
					value := fmt.Sprint(" Series: ", key, ", Index: ", idx, ", Value: ", point.Value(), " [ ", point.Timestamp(), " ]")
					w.enableMouseContainer(value, point.ColorName(), &me.Position)
				}
			}
		}
	}
}

// MouseOut disable display of mouse data point display
func (w *LineChartSkn) MouseOut() {
	w.disableMouseContainer()
}

// enableMouseContainer private method to prepare values need by renderer to create pop display
// composes display text, captures position and colorName for use by renderer
func (w *LineChartSkn) enableMouseContainer(value, frameColor string, mousePosition *fyne.Position) *LineChartSkn {
	w.mouseDisplayStr = value
	w.mouseDisplayFrameColor = frameColor
	ct := canvas.NewText(value, theme.PrimaryColorNamed(frameColor))
	parts := strings.Split(value, "[")
	ts := fyne.MeasureText(parts[0], ct.TextSize, ct.TextStyle)
	mp := &fyne.Position{X: mousePosition.X - (ts.Width / 2), Y: mousePosition.Y - (3 * ts.Height) - theme.Padding()}
	w.mouseDisplayPosition = mp
	return w
}

// disableMouseContainer private method to manage mouse leaving window
// blank string will prevent display
func (w *LineChartSkn) disableMouseContainer() {
	w.mouseDisplayStr = ""
	w.Refresh()
}

// Widget Renderer code starts here
type sknLineChartRenderer struct {
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

// Create the renderer with a reference to the widget
// and all the objects to be displayed for this custom widget
//
// Note: Do not size or move canvas objects here.
func newSknLineChartRenderer(lineChart *LineChartSkn) *sknLineChartRenderer {
	fmt.Println("sknLineChartRenderer::newSknLineChartRenderer() called")
	background := canvas.NewRectangle(color.Transparent)
	background.StrokeWidth = 0.75
	background.StrokeColor = theme.PrimaryColorNamed(theme.ColorBlue)

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

	}
	for i := 0; i < 12; i++ {
		yt := strconv.Itoa((11 - i) * 10)
		yl := canvas.NewText(yt, theme.ForegroundColor())
		yl.Alignment = fyne.TextAlignTrailing
		yLabels = append(yLabels, yl)
	}
	for i := 0; i < 13; i++ {
		xt := strconv.Itoa(i * 10)
		xl := canvas.NewText(xt, theme.ForegroundColor())
		xl.Alignment = fyne.TextAlignTrailing
		xLabels = append(xLabels, xl)
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

	topCenteredDesc := canvas.NewText(lineChart.topCenteredDesc, theme.ForegroundColor())
	topCenteredDesc.TextSize = 24
	topCenteredDesc.TextStyle = fyne.TextStyle{
		Bold:   true,
		Italic: false,
	}
	bottomCenteredDesc := canvas.NewText(lineChart.bottomCenteredDesc, theme.ForegroundColor())
	bottomCenteredDesc.TextSize = 16
	bottomCenteredDesc.TextStyle = fyne.TextStyle{
		Bold:   false,
		Italic: true,
	}

	// vertical text for X/Y legends since no text rotation is available
	lBox := container.NewVBox()
	for _, c := range lineChart.leftMiddleDesc {
		lBox.Add(
			canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground))),
		)
	}
	rBox := container.NewVBox()
	for _, c := range lineChart.rightMiddleDesc {
		rBox.Add(
			canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground))),
		)
	}

	return &sknLineChartRenderer{
		widget:                lineChart,
		chartFrame:            background,
		xLines:                xlines,
		yLines:                ylines,
		xLabels:               xLabels,
		yLabels:               yLabels,
		dataPoints:            dataPoints,
		topLeftDesc:           canvas.NewText(lineChart.topLeftDesc, theme.ForegroundColor()),
		topCenteredDesc:       topCenteredDesc,
		topRightDesc:          canvas.NewText(lineChart.topRightDesc, theme.ForegroundColor()),
		bottomLeftDesc:        canvas.NewText(lineChart.bottomLeftDesc, theme.ForegroundColor()),
		bottomCenteredDesc:    bottomCenteredDesc,
		bottomRightDesc:       canvas.NewText(lineChart.bottomRightDesc, theme.ForegroundColor()),
		leftMiddleBox:         lBox,
		rightMiddleBox:        rBox,
		dataPointMarkers:      dpMaker,
		mouseDisplayContainer: mouseDisplay,
	}
}

// Refresh method is called if the state of the widget changes or the
// theme is changed
func (r *sknLineChartRenderer) Refresh() {
	fmt.Println(" ====> sknLineChartRenderer::Refresh() called")
	r.VerifyDataPoints()

	r.chartFrame.Refresh()            // Redraw the chartFrame first
	for idx, line := range r.xLines { // grid
		line.Refresh()
		r.yLines[idx].Refresh()
	}
	for _, xlbl := range r.xLabels { // labels
		xlbl.Refresh()
	}
	for _, ylbl := range r.yLabels { // labels
		ylbl.Refresh()
	}
	for key, lines := range r.dataPoints { // data points, and markers
		for idx, point := range lines {
			point.Refresh()
			r.dataPointMarkers[key][idx].Refresh()
		}
	}

	r.mouseDisplayContainer.Objects[0].(*canvas.Rectangle).StrokeColor = theme.PrimaryColorNamed(r.widget.mouseDisplayFrameColor)
	r.mouseDisplayContainer.Objects[1].(*widget.Label).SetText(r.widget.mouseDisplayStr)

	r.topLeftDesc.Text = r.widget.topLeftDesc
	r.topCenteredDesc.Text = r.widget.topCenteredDesc
	r.topRightDesc.Text = r.widget.topRightDesc
	r.bottomLeftDesc.Text = r.widget.bottomLeftDesc
	r.bottomCenteredDesc.Text = r.widget.bottomCenteredDesc
	r.bottomRightDesc.Text = r.widget.bottomRightDesc

	r.leftMiddleBox.RemoveAll()
	for _, c := range r.widget.leftMiddleDesc {
		z := canvas.NewText(
			strings.ToUpper(string(c)),
			theme.PrimaryColorNamed(string(theme.ColorNameForeground)))
		z.TextSize = 12
		r.leftMiddleBox.Add(z)
	}
	r.rightMiddleBox.RemoveAll()
	for _, c := range r.widget.rightMiddleDesc {
		z := canvas.NewText(
			strings.ToUpper(string(c)),
			theme.PrimaryColorNamed(string(theme.ColorNameForeground)))
		z.TextSize = 12
		r.rightMiddleBox.Add(z)
	}

	r.topLeftDesc.Refresh()
	r.topCenteredDesc.Refresh()
	r.topRightDesc.Refresh()

	r.bottomLeftDesc.Refresh()
	r.bottomCenteredDesc.Refresh()
	r.bottomRightDesc.Refresh()
}

// Layout Given the size required by the fyne application
// move and re-size the all custom widget canvas objects here
func (r *sknLineChartRenderer) Layout(s fyne.Size) {
	fmt.Println("sknLineChartRenderer::Layout() called")
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
		if nil == r.dataPoints[key] {
			r.dataPoints[key] = []*canvas.Line{}
			r.dataPointMarkers[key] = []*canvas.Circle{}
			fmt.Println(" ====> sknLineChartRenderer::Layout() called, add series: ", key)
		}
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

			if idx > (len(r.dataPoints[key]) - 1) { // concurrency error !!!
				x := canvas.NewLine(theme.PrimaryColorNamed(point.ColorName()))
				x.StrokeWidth = 2.0
				r.dataPoints[key] = append(r.dataPoints[key], x)
				z := canvas.NewCircle(theme.PrimaryColorNamed(point.ColorName()))
				z.StrokeWidth = 4.0
				r.dataPointMarkers[key] = append(r.dataPointMarkers[key], z)
				fmt.Println(" ====> sknLineChartRenderer::Layout() called, add points: ", key, ", index: ", idx)
			}
			r.dataPoints[key][idx].Position1 = thisPoint
			r.dataPoints[key][idx].Position2 = lastPoint
			lastPoint = thisPoint

			zt := fyne.NewPos(thisPoint.X-2, thisPoint.Y-2)
			r.dataPointMarkers[key][idx].Position1 = zt
			zb := fyne.NewPos(thisPoint.X+2, thisPoint.Y+2)
			r.dataPointMarkers[key][idx].Position2 = zb
			point.SetMarkerPosition(&zt, &zb)
			r.dataPointMarkers[key][idx].Resize(fyne.NewSize(5, 5))
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

}

// MinSize Create a minimum size for the widget.
// The smallest size is can be overridden by user
func (r *sknLineChartRenderer) MinSize() fyne.Size {
	fmt.Println("sknLineChartRenderer::MinSize() called")
	return r.widget.minSize
}

// Objects Return a list of each canvas object.
// but only the objects that have been enabled or are not at default value; i.e. ""
func (r *sknLineChartRenderer) Objects() []fyne.CanvasObject {
	fmt.Println("sknLineChartRenderer::Objects() called")
	objs := []fyne.CanvasObject{
		r.chartFrame,
	}
	if r.topLeftDesc.Text != "" {
		objs = append(objs, r.topLeftDesc)
	}
	if r.topCenteredDesc.Text != "" {
		objs = append(objs, r.topCenteredDesc)
	}
	if r.topRightDesc.Text != "" {
		objs = append(objs, r.topRightDesc)
	}
	if r.widget.leftMiddleDesc != "" {
		objs = append(objs, r.leftMiddleBox)
	}
	if r.widget.rightMiddleDesc != "" {
		objs = append(objs, r.rightMiddleBox)
	}
	if r.bottomLeftDesc.Text != "" {
		objs = append(objs, r.bottomLeftDesc)
	}
	if r.bottomCenteredDesc.Text != "" {
		objs = append(objs, r.bottomCenteredDesc)
	}
	if r.bottomRightDesc.Text != "" {
		objs = append(objs, r.bottomRightDesc)
	}
	for idx, line := range r.xLines {
		if r.widget.enableHorizGridLines {
			objs = append(objs, line)
		}
		if r.widget.enableVertGridLines {
			objs = append(objs, r.yLines[idx])
		}
	}
	for _, lbl := range r.yLabels {
		objs = append(objs, lbl)
	}
	for _, lbl := range r.xLabels {
		objs = append(objs, lbl)
	}
	for key, lines := range r.dataPoints {
		for idx, line := range lines {
			objs = append(objs, line)
			if r.widget.enableDataPointMarkers {
				objs = append(objs, r.dataPointMarkers[key][idx])
			}
		}
	}
	if r.mouseDisplayContainer.Objects[1].(*widget.Label).Text != "" {
		objs = append(objs, r.mouseDisplayContainer)
	}

	return objs
}

// Destroy Cleanup if resources have been allocated
func (r *sknLineChartRenderer) Destroy() {}

// VerifyDataPoints Renderer method to inject newly add data series or points
// called by Refresh() to ensure new data is recognized
func (r *sknLineChartRenderer) VerifyDataPoints() {
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
}
