package components

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
// A text widget with themed background and foreground
type sknLineChart struct {
	widget.BaseWidget  // Inherit from BaseWidget
	dataPoints         *map[string][]interfaces.SknDataSeries
	dataPointColor     color.Color
	dataPointScale     fyne.Size
	dataPointLimit     int
	topLeftDesc        string // The text to display in the widget
	topCenteredDesc    string
	topRightDesc       string
	leftMiddleDesc     string
	rightMiddleDesc    string
	bottomCenteredDesc string
	bottomLeftDesc     string
	bottomRightDesc    string
	minSize            fyne.Size
}

var _ interfaces.SknLineChart = (*sknLineChart)(nil)

// NewSknLineChart Create a Widget and Extend (initialize) the BaseWidget
func NewSknLineChart(title, xScale, yScale string, dataPoints *map[string][]interfaces.SknDataSeries) (*sknLineChart, error) {
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
		dataPoints:         dataPoints,
		dataPointLimit:     dpl,
		dataPointColor:     theme.PrimaryColor(),
		dataPointScale:     fyne.NewSize(120.0, 110.0),
		topLeftDesc:        "top left description",
		topCenteredDesc:    title,
		topRightDesc:       "top right description",
		leftMiddleDesc:     "left middle description",
		rightMiddleDesc:    yScale,
		bottomLeftDesc:     "bottom left description",
		bottomCenteredDesc: xScale,
		bottomRightDesc:    "bottom right description",
		minSize:            fyne.NewSize(200+theme.Padding()*4, 150+theme.Padding()*4),
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
func (w *sknLineChart) ReplaceDataSeries(newData *map[string][]interfaces.SknDataSeries) error {
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
		w.validateDataPoints()
	} else {
		return fmt.Errorf("newData must be larger[%d] than or equal to existing[%d]", len(*newData), len(*w.dataPoints))
	}
	return nil
}
func (w *sknLineChart) ApplyNewDataSeries(seriesName string, newSeries []interfaces.SknDataSeries) error {
	if w == nil {
		return fmt.Errorf("ApplyNewDataSeries() no active widget")
	}

	if len(newSeries) < w.dataPointLimit {
		(*w.dataPoints)[seriesName] = newSeries
		w.validateDataPoints()
	} else {
		return fmt.Errorf("[%s] data series datapoints limit exceeded. limit:%d, count:%d", seriesName, w.dataPointLimit, len(newSeries))
	}
	return nil
}
func (w *sknLineChart) ApplySingleDataPoint(seriesName string, newDataPoint interfaces.SknDataSeries) {
	if w == nil {
		return
	}
	if len((*w.dataPoints)[seriesName]) < w.dataPointLimit {
		(*w.dataPoints)[seriesName] = append((*w.dataPoints)[seriesName], newDataPoint)
	} else {
		(*w.dataPoints)[seriesName] = commons.ShiftSlice(newDataPoint, (*w.dataPoints)[seriesName])
	}
	w.validateDataPoints()
}

// validateDataPoints
//
// Create missing line objects in renderer
func (w *sknLineChart) validateDataPoints() {
	w.Refresh()
}

// Widget Renderer code starts here
type sknLineChartRenderer struct {
	widget             *sknLineChart     // Reference to the widget holding the current state
	background         *canvas.Rectangle // A background rectangle
	xInc               float32
	yInc               float32
	dataPoints         map[string][]*canvas.Line
	xLines             []*canvas.Line
	yLines             []*canvas.Line
	xLabels            []*canvas.Text
	yLabels            []*canvas.Text
	topLeftDesc        *canvas.Text
	topCenteredDesc    *canvas.Text
	topRightDesc       *canvas.Text
	leftMiddleDesc     *canvas.Text
	rightMiddleDesc    *canvas.Text
	bottomCenteredDesc *canvas.Text
	bottomLeftDesc     *canvas.Text
	bottomRightDesc    *canvas.Text
	leftMiddleBox      *fyne.Container
	rightMiddleBox     *fyne.Container
}

// Create the renderer with a reference to the widget
// Note: The background and foreground colours are set from the current theme.
//
// Do not size or move canvas objects here.
func newSknLineChartRenderer(lineChart *sknLineChart) *sknLineChartRenderer {
	background := canvas.NewRectangle(color.Transparent)
	background.StrokeWidth = 0.75
	background.StrokeColor = theme.PrimaryColorNamed(theme.ColorBlue)

	dataPoints := map[string][]*canvas.Line{}
	xlines := []*canvas.Line{}
	ylines := []*canvas.Line{}
	xLabels := []*canvas.Text{}
	yLabels := []*canvas.Text{}
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

	leftMiddleDesc := canvas.NewText(lineChart.leftMiddleDesc, theme.ForegroundColor())
	leftMiddleDesc.TextSize = 18
	rightMiddleDesc := canvas.NewText(lineChart.rightMiddleDesc, theme.ForegroundColor())
	rightMiddleDesc.TextSize = 18

	lBox := container.NewVBox()
	for _, c := range lineChart.leftMiddleDesc {
		lBox.Add(canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground))))
	}
	rBox := container.NewVBox()
	for _, c := range lineChart.rightMiddleDesc {
		rBox.Add(canvas.NewText(strings.ToUpper(string(c)), theme.PrimaryColorNamed(string(theme.ColorNameForeground))))
	}

	return &sknLineChartRenderer{
		widget:             lineChart,
		background:         background,
		xLines:             xlines,
		yLines:             ylines,
		xLabels:            xLabels,
		yLabels:            yLabels,
		dataPoints:         dataPoints,
		topLeftDesc:        canvas.NewText(lineChart.topLeftDesc, theme.ForegroundColor()),
		topCenteredDesc:    topCenteredDesc,
		topRightDesc:       canvas.NewText(lineChart.topRightDesc, theme.ForegroundColor()),
		leftMiddleDesc:     leftMiddleDesc,
		rightMiddleDesc:    rightMiddleDesc,
		bottomLeftDesc:     canvas.NewText(lineChart.bottomLeftDesc, theme.ForegroundColor()),
		bottomCenteredDesc: bottomCenteredDesc,
		bottomRightDesc:    canvas.NewText(lineChart.bottomRightDesc, theme.ForegroundColor()),
		leftMiddleBox:      lBox,
		rightMiddleBox:     rBox,
	}
}

// Refresh method is called if the state of the widget changes or the
// theme is changed
//
// Note: The background and foreground colours are set from the current theme
func (r *sknLineChartRenderer) Refresh() {
	r.VerifyDataPoints()

	r.background.Refresh()            // Redraw the background first
	for idx, line := range r.xLines { // grid
		line.Refresh()
		r.yLines[idx].Refresh()
	}
	for idx, xlbl := range r.xLabels { // labels
		xlbl.Refresh()
		r.yLabels[idx].Refresh()
	}
	for _, lines := range r.dataPoints { // data points
		for _, point := range lines {
			point.Refresh()
		}
	}
	r.topLeftDesc.Text = r.widget.topLeftDesc
	r.topCenteredDesc.Text = r.widget.topCenteredDesc
	r.topRightDesc.Text = r.widget.topRightDesc
	r.leftMiddleDesc.Text = r.widget.leftMiddleDesc
	r.rightMiddleDesc.Text = r.widget.rightMiddleDesc
	r.bottomLeftDesc.Text = r.widget.bottomLeftDesc
	r.bottomCenteredDesc.Text = r.widget.bottomCenteredDesc
	r.bottomRightDesc.Text = r.widget.bottomRightDesc
	r.topLeftDesc.Refresh()
	r.topCenteredDesc.Refresh()
	r.topRightDesc.Refresh()
	r.leftMiddleDesc.Refresh()
	r.rightMiddleDesc.Refresh()
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

			if idx > (len(r.dataPoints[key]) - 1) {
				x := canvas.NewLine(theme.PrimaryColorNamed(point.ColorName()))
				x.StrokeWidth = 2.0
				r.dataPoints[key][idx] = x
			}
			r.dataPoints[key][idx].Position1 = thisPoint
			r.dataPoints[key][idx].Position2 = lastPoint
			r.dataPoints[key][idx].Refresh()
			lastPoint = thisPoint
		}
	}

	r.rightMiddleDesc.Resize(fyne.NewSize(12, 200))
	r.leftMiddleDesc.Resize(fyne.NewSize(12, 200))

	r.background.Resize(fyne.NewSize(r.xInc*12, r.yInc*11))
	r.background.Move(fyne.NewPos(r.xInc, r.yInc))

	ts := fyne.MeasureText(r.topCenteredDesc.Text, r.topCenteredDesc.TextSize, r.topCenteredDesc.TextStyle)
	r.topCenteredDesc.Move(fyne.Position{X: (s.Width - ts.Width) / 2, Y: theme.Padding() / 4})

	ts = fyne.MeasureText(r.topRightDesc.Text, r.topRightDesc.TextSize, r.topRightDesc.TextStyle)
	r.topRightDesc.Move(fyne.Position{X: (s.Width - ts.Width) - theme.Padding(), Y: ts.Height / 4})
	r.topLeftDesc.Move(fyne.NewPos(theme.Padding(), ts.Height/4))

	ts = fyne.MeasureText(r.leftMiddleDesc.Text, r.leftMiddleDesc.TextSize, r.leftMiddleDesc.TextStyle)
	r.leftMiddleDesc.Move(fyne.NewPos(theme.Padding(), s.Height/3))
	r.leftMiddleBox.Resize(fyne.NewSize(ts.Height, s.Height*0.75))
	r.leftMiddleBox.Move(fyne.NewPos(theme.Padding(), s.Height*0.1))

	ts = fyne.MeasureText(r.rightMiddleDesc.Text, r.rightMiddleDesc.TextSize, r.rightMiddleDesc.TextStyle)
	r.rightMiddleDesc.Move(fyne.NewPos((s.Width-ts.Width)-theme.Padding(), s.Height/3))
	r.rightMiddleBox.Resize(fyne.NewSize(ts.Height, s.Height*0.75))
	r.rightMiddleBox.Move(fyne.NewPos((s.Width - ts.Height), s.Height*0.1))

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
		r.background,
		r.topLeftDesc, r.topCenteredDesc, r.topRightDesc,
		r.leftMiddleBox, r.rightMiddleBox,
		r.bottomLeftDesc, r.bottomCenteredDesc, r.bottomRightDesc,
	}
	for idx, line := range r.xLines {
		objs = append(objs, line)
		objs = append(objs, r.yLines[idx])
	}
	for idx, lbl := range r.yLabels {
		objs = append(objs, lbl)
		objs = append(objs, r.xLabels[idx])
	}
	for _, lines := range r.dataPoints {
		for _, line := range lines {
			objs = append(objs, line)
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
		}
		for idx, point := range points {
			if idx > (len(r.dataPoints[key]) - 1) {
				x := canvas.NewLine(theme.PrimaryColorNamed(point.ColorName()))
				x.StrokeWidth = 2.0
				r.dataPoints[key] = append(r.dataPoints[key], x)
			}
		}
	}
}
