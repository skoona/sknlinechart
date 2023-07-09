package components

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"ggApcMon/internal/commons"
	"ggApcMon/internal/interfaces"
	"image/color"
	"strconv"
	"strings"
)

// Widget code starts here
//
// A text widget with themed chartFrame and foreground
type sknLineChart struct {
	widget.BaseWidget       // Inherit from BaseWidget
	dataPoints              *map[string][]interfaces.SknDatapoint
	dataPointColor          color.Color
	dataPointScale          fyne.Size
	dataPointLimit          int
	enableDataPointMarkers  bool
	enableHorizGridLines    bool
	enableVertGridLines     bool
	enableMousePointDisplay bool
	mouseDisplayStr         string
	mouseDisplayPosition    *fyne.Position
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

var _ interfaces.SknLineChart = (*sknLineChart)(nil)

// NewSknLineChart Create a Widget and Extend (initialize) the BaseWidget
func NewSknLineChart(tTitle, bTitle string, dataPoints *map[string][]interfaces.SknDatapoint) (*sknLineChart, error) {
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

	w := &sknLineChart{ // Create this widget with an initial text value
		dataPoints:              dataPoints,
		dataPointLimit:          dpl,
		dataPointColor:          theme.PrimaryColor(),
		dataPointScale:          fyne.NewSize(120.0, 110.0),
		enableDataPointMarkers:  true,
		enableHorizGridLines:    true,
		enableVertGridLines:     true,
		enableMousePointDisplay: true,
		mouseDisplayStr:         "",
		mouseDisplayPosition:    &fyne.Position{},
		topLeftDesc:             "top left desc",
		topCenteredDesc:         tTitle,
		topRightDesc:            "",
		leftMiddleDesc:          "left middle desc",
		rightMiddleDesc:         "right middle desc",
		bottomLeftDesc:          "bottom left desc",
		bottomCenteredDesc:      bTitle,
		bottomRightDesc:         "bottom right desc",
		minSize:                 fyne.NewSize(200+theme.Padding()*4, 150+theme.Padding()*4),
	}
	w.ExtendBaseWidget(w) // Initialize the BaseWidget
	return w, err
}

// CreateRenderer Create the renderer. This is called by the fyne application
func (w *sknLineChart) CreateRenderer() fyne.WidgetRenderer {
	return newSknLineChartRenderer(w)
}
func (w *sknLineChart) SetMinSize(s fyne.Size) {
	w.minSize = s
}
func (w *sknLineChart) GetTopLeftDescription() string {
	return w.topLeftDesc
}
func (w *sknLineChart) GetTitle() string {
	return w.topCenteredDesc
}
func (w *sknLineChart) IsDataPointMarkersEnabled() bool {
	return w.enableDataPointMarkers
}
func (w *sknLineChart) IsHorizGridLinesEnabled() bool {
	return w.enableHorizGridLines
}
func (w *sknLineChart) IsVertGridLinesEnabled() bool {
	return w.enableVertGridLines
}
func (w *sknLineChart) IsMousePointDisplayEnabled() bool {
	return w.enableMousePointDisplay
}
func (w *sknLineChart) GetTopRightDescription() string {
	return w.topRightDesc
}
func (w *sknLineChart) GetMiddleLeftDescription() string {
	return w.leftMiddleDesc
}
func (w *sknLineChart) GetMiddleRightDescription() string {
	return w.rightMiddleDesc
}
func (w *sknLineChart) GetBottomLeftDescription() string {
	return w.bottomLeftDesc
}
func (w *sknLineChart) GetBottomCenteredDescription() string {
	return w.bottomCenteredDesc
}
func (w *sknLineChart) GetBottomRightDescription() string {
	return w.bottomRightDesc
}
func (w *sknLineChart) SetTopLeftDescription(newValue string) {
	w.topLeftDesc = newValue
}
func (w *sknLineChart) SetTitle(newValue string) {
	w.topCenteredDesc = newValue
}
func (w *sknLineChart) EnableDataPointMarkers(newValue bool) {
	w.enableDataPointMarkers = newValue
}
func (w *sknLineChart) EnabledHorizGridLines(newValue bool) {
	w.enableHorizGridLines = newValue
}
func (w *sknLineChart) EnableVertGridLine(newValue bool) {
	w.enableVertGridLines = newValue
}
func (w *sknLineChart) EnableMousePointDisplay(newValue bool) {
	w.enableMousePointDisplay = newValue
}
func (w *sknLineChart) SetTopRightDescription(newValue string) {
	w.topRightDesc = newValue
}
func (w *sknLineChart) SetMiddleLeftDescription(newValue string) {
	w.leftMiddleDesc = newValue
}
func (w *sknLineChart) SetMiddleRightDescription(newValue string) {
	w.rightMiddleDesc = newValue
}
func (w *sknLineChart) SetBottomLeftDescription(newValue string) {
	w.bottomLeftDesc = newValue
}
func (w *sknLineChart) SetBottomRightDescription(newValue string) {
	w.bottomRightDesc = newValue
}
func (w *sknLineChart) SetBottomCenteredDescription(newValue string) {
	w.bottomCenteredDesc = newValue
}
func (w *sknLineChart) SetDataSeriesColor(c color.Color) {
	w.dataPointColor = c
}
func (w *sknLineChart) ReplaceDataSeries(newData *map[string][]interfaces.SknDatapoint) error {
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
func (w *sknLineChart) ApplyNewDataSeries(seriesName string, newSeries []interfaces.SknDatapoint) error {
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
func (w *sknLineChart) ApplySingleDataPoint(seriesName string, newDataPoint interfaces.SknDatapoint) {
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
func (w *sknLineChart) MouseIn(me *desktop.MouseEvent) {

}
func (w *sknLineChart) MouseMoved(me *desktop.MouseEvent) {
	if !w.enableMousePointDisplay {
		return
	}
	for key, points := range *w.dataPoints {
		for idx, point := range points {
			top, bottom := point.MarkerPosition()
			if !me.Position.IsZero() && !top.IsZero() {
				if me.Position.X >= top.X && me.Position.X <= bottom.X &&
					me.Position.Y >= top.Y && me.Position.Y <= bottom.Y {
					value := fmt.Sprint("Match Series: ", key, ", Index: ", idx, ", Color: ", point.ColorName(), ", Value: ", point.Value())
					w.enableMouseContainer(value, &me.Position).Refresh()
				} else {
					w.disableMouseContainer()
				}
			}
		}
	}
}
func (w *sknLineChart) MouseOut() {
	w.disableMouseContainer()
}
func (w *sknLineChart) enableMouseContainer(value string, mousePosition *fyne.Position) *sknLineChart {
	w.mouseDisplayStr = value
	mp := &fyne.Position{X: mousePosition.X - 20, Y: mousePosition.Y - 20}
	w.mouseDisplayPosition = mp
	fmt.Println(value)
	return w
}
func (w *sknLineChart) disableMouseContainer() {
	w.mouseDisplayStr = ""
}

// Widget Renderer code starts here
type sknLineChartRenderer struct {
	widget                *sknLineChart     // Reference to the widget holding the current state
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
// Note: The chartFrame and foreground colours are set from the current theme.
//
// Do not size or move canvas objects here.
func newSknLineChartRenderer(lineChart *sknLineChart) *sknLineChartRenderer {
	background := canvas.NewRectangle(color.Transparent)
	background.StrokeWidth = 0.75
	background.StrokeColor = theme.PrimaryColorNamed(theme.ColorBlue)

	mouseDisplay := container.NewPadded(canvas.NewText("", theme.ForegroundColor()))

	dataPoints := map[string][]*canvas.Line{}
	xlines := []*canvas.Line{}
	ylines := []*canvas.Line{}
	xLabels := []*canvas.Text{}
	yLabels := []*canvas.Text{}
	dpMaker := map[string][]*canvas.Circle{}
	xl := canvas.NewText("110", theme.ForegroundColor())
	yl := canvas.NewText("120", theme.ForegroundColor())
	xLabels = append(xLabels, xl) // x12
	yLabels = append(yLabels, yl) // x13 hi-jacked
	for i := 0; i < 11; i++ {
		x := canvas.NewLine(theme.PrimaryColorNamed(theme.ColorGreen))
		x.StrokeWidth = 0.25
		y := canvas.NewLine(theme.PrimaryColorNamed(theme.ColorGreen))
		y.StrokeWidth = 0.25
		xlines = append(xlines, x)
		ylines = append(ylines, y)
		xt := strconv.Itoa((10 - i) * 10)
		yt := strconv.Itoa((10 - i) * 10)
		xl := canvas.NewText(xt, theme.ForegroundColor())
		yl := canvas.NewText(yt, theme.ForegroundColor())
		xLabels = append(xLabels, xl)
		yLabels = append(yLabels, yl)
	}

	for key, points := range *lineChart.dataPoints {
		// color
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
	bottomCenteredDesc.TextSize = 14
	bottomCenteredDesc.TextStyle = fyne.TextStyle{
		Bold:   false,
		Italic: true,
	}

	lBox := container.NewVBox()
	for _, c := range lineChart.leftMiddleDesc {
		lBox.Add(canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground))))
	}
	rBox := container.NewVBox()
	for _, c := range lineChart.rightMiddleDesc {
		rBox.Add(canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground))))
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
//
// Note: The chartFrame and foreground colours are set from the current theme
func (r *sknLineChartRenderer) Refresh() {
	r.VerifyDataPoints()

	r.chartFrame.Refresh()            // Redraw the chartFrame first
	for idx, line := range r.xLines { // grid
		line.Refresh()
		r.yLines[idx].Refresh()
	}
	for idx, xlbl := range r.xLabels { // labels
		xlbl.Refresh()
		r.yLabels[idx].Refresh()
	}
	for key, lines := range r.dataPoints { // data points
		for idx, point := range lines {
			point.Refresh()
			r.dataPointMarkers[key][idx].Refresh()
		}
	}

	r.mouseDisplayContainer.Objects[0].(*canvas.Text).Text = r.widget.mouseDisplayStr
	r.mouseDisplayContainer.Refresh()
	r.topLeftDesc.Text = r.widget.topLeftDesc
	r.topCenteredDesc.Text = r.widget.topCenteredDesc
	r.topRightDesc.Text = r.widget.topRightDesc
	r.bottomLeftDesc.Text = r.widget.bottomLeftDesc
	r.bottomCenteredDesc.Text = r.widget.bottomCenteredDesc
	r.bottomRightDesc.Text = r.widget.bottomRightDesc

	r.leftMiddleBox.RemoveAll()
	for _, c := range r.widget.leftMiddleDesc {
		z := canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground)))
		z.TextSize = 14
		r.leftMiddleBox.Add(z)
	}
	r.rightMiddleBox.RemoveAll()
	for _, c := range r.widget.rightMiddleDesc {
		z := canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground)))
		z.TextSize = 14
		r.rightMiddleBox.Add(z)
	}

	r.topLeftDesc.Refresh()
	r.topCenteredDesc.Refresh()
	r.topRightDesc.Refresh()
	r.bottomCenteredDesc.Refresh()
	r.bottomLeftDesc.Refresh()
	r.bottomRightDesc.Refresh()
}

// Layout Given the size required by the fyne application move and re-size the
// canvas objects.
func (r *sknLineChartRenderer) Layout(s fyne.Size) {
	r.xInc = (s.Width - theme.Padding()) / 14.0
	r.yInc = (s.Height - theme.Padding()) / 14.0

	// Vert lines
	yp := 12.0 * r.yInc
	for idx, line := range r.xLines {
		xp := float32(idx+1) * r.xInc
		line.Position1 = fyne.NewPos(xp+r.xInc, r.yInc) //top
		line.Position2 = fyne.NewPos(xp+r.xInc, yp+8)
		line.Refresh()
		r.xLabels[idx+1].Move(fyne.NewPos(xp+(2*r.xInc)-8, yp+10))
		r.xLabels[idx+1].Refresh()
	}
	r.xLabels[0].Move(fyne.NewPos((2*r.xInc)-10, yp+10))
	r.xLabels[0].Refresh()

	// Horiz lines
	xp := r.xInc
	for idx, line := range r.yLines {
		yp := float32(idx+1) * r.yInc

		line.Position1 = fyne.NewPos(xp, yp+r.yInc) // left
		line.Position2 = fyne.NewPos(xp*13+8, yp+r.yInc)
		line.Refresh()
		r.yLabels[idx+1].Move(fyne.NewPos(xp*13+10, yp+r.yInc-8))
		r.yLabels[idx+1].Refresh()
	}
	r.yLabels[0].Move(fyne.NewPos(xp-10, yp+8))
	r.yLabels[0].Refresh()

	// data points
	xp = r.xInc * 13
	yp = r.yInc * 12
	yScale := (r.yInc * 10) / 100
	xScale := (r.xInc * 10) / 100
	dp := float32(1.0)
	for key, data := range *r.widget.dataPoints { // datasource
		lastPoint := fyne.NewPos(xp, yp)
		if nil == r.dataPoints[key] {
			r.dataPoints[key] = []*canvas.Line{}
			r.dataPointMarkers[key] = []*canvas.Circle{}
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
			xx := xp - (float32(idx) * xScale)
			thisPoint := fyne.NewPos(xx, yy)
			if idx == 0 {
				lastPoint.Y = yy
			}

			if idx > (len(r.dataPoints[key]) - 1) {
				x := canvas.NewLine(theme.PrimaryColorNamed(point.ColorName()))
				x.StrokeWidth = 2.0
				r.dataPoints[key][idx] = x
				z := canvas.NewCircle(theme.PrimaryColorNamed(point.ColorName()))
				z.StrokeWidth = 4.0
				r.dataPointMarkers[key][idx] = z
			}
			r.dataPoints[key][idx].Position1 = thisPoint
			r.dataPoints[key][idx].Position2 = lastPoint
			r.dataPoints[key][idx].Refresh()
			lastPoint = thisPoint
			zt := fyne.NewPos(thisPoint.X-2, thisPoint.Y-2)
			r.dataPointMarkers[key][idx].Position1 = zt
			zb := fyne.NewPos(thisPoint.X+2, thisPoint.Y+2)
			r.dataPointMarkers[key][idx].Position2 = zb
			point.SetMarkerPosition(&zt, &zb)
			r.dataPointMarkers[key][idx].Resize(fyne.NewSize(5, 5))
			r.dataPointMarkers[key][idx].Refresh()
		}
	}

	r.chartFrame.Resize(fyne.NewSize(r.xInc*12, r.yInc*11))
	r.chartFrame.Move(fyne.NewPos(r.xInc, r.yInc))

	ts := fyne.MeasureText(r.topCenteredDesc.Text, r.topCenteredDesc.TextSize, r.topCenteredDesc.TextStyle)
	r.topCenteredDesc.Move(fyne.Position{X: (s.Width - ts.Width) / 2, Y: theme.Padding() / 4})

	ts = fyne.MeasureText(r.topRightDesc.Text, r.topRightDesc.TextSize, r.topRightDesc.TextStyle)
	r.topRightDesc.Move(fyne.Position{X: (s.Width - ts.Width) - theme.Padding(), Y: ts.Height / 4})
	r.topLeftDesc.Move(fyne.NewPos(theme.Padding(), ts.Height/4))

	r.mouseDisplayContainer.Move(*r.widget.mouseDisplayPosition)

	r.leftMiddleBox.Resize(fyne.NewSize(ts.Height, s.Height*0.75))
	r.leftMiddleBox.Move(fyne.NewPos(2*theme.Padding(), s.Height*0.1))

	r.rightMiddleBox.Resize(fyne.NewSize(ts.Height, s.Height*0.75))
	r.rightMiddleBox.Move(fyne.NewPos(s.Width-ts.Height, s.Height*0.1))

	ts = fyne.MeasureText(r.bottomCenteredDesc.Text, r.bottomCenteredDesc.TextSize, r.bottomCenteredDesc.TextStyle)
	r.bottomCenteredDesc.Move(fyne.Position{X: (s.Width - ts.Width) / 2, Y: s.Height - ts.Height - theme.Padding()})

	ts = fyne.MeasureText(r.bottomRightDesc.Text, r.bottomRightDesc.TextSize, r.bottomRightDesc.TextStyle)
	r.bottomRightDesc.Move(fyne.NewPos((s.Width-ts.Width)-theme.Padding(), s.Height-ts.Height-theme.Padding()))
	r.bottomLeftDesc.Move(fyne.NewPos(theme.Padding()+2.0, s.Height-ts.Height-theme.Padding()))

}

// MinSize Create a minimum size for the widget.
// The smallest size is the size of the text with a border defined by the theme padding
func (r *sknLineChartRenderer) MinSize() fyne.Size {
	return r.widget.minSize
}

// Objects Return a list of each canvas object.
func (r *sknLineChartRenderer) Objects() []fyne.CanvasObject {
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
	if r.mouseDisplayContainer.Objects[0].(*canvas.Text).Text != "" {
		objs = append(objs, r.mouseDisplayContainer)
	}
	for idx, line := range r.xLines {
		if r.widget.enableHorizGridLines {
			objs = append(objs, line)
		}
		if r.widget.enableVertGridLines {
			objs = append(objs, r.yLines[idx])
		}
	}
	for idx, lbl := range r.yLabels {
		objs = append(objs, lbl)
		objs = append(objs, r.xLabels[idx])
	}
	for key, lines := range r.dataPoints {
		for idx, line := range lines {
			objs = append(objs, line)
			if r.widget.enableDataPointMarkers {
				objs = append(objs, r.dataPointMarkers[key][idx])
			}
		}
	}
	return objs
}

// Destroy Cleanup if resources have been allocated
func (r *sknLineChartRenderer) Destroy() {}

// VerifyDataPoints
//
// Renderer method to inject newly add data series or points
func (r *sknLineChartRenderer) VerifyDataPoints() {
	for key, points := range *r.widget.dataPoints {
		if nil == r.dataPoints[key] {
			r.dataPoints[key] = []*canvas.Line{}
			r.dataPointMarkers[key] = []*canvas.Circle{}
		}
		for idx, point := range points {
			if idx > (len(r.dataPoints[key]) - 1) {
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
