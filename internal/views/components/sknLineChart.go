package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"ggApcMon/internal/interfaces"
	"image/color"
	"strconv"
)

// Widget code starts here
//
// A text widget with themed background and foreground
type sknLineChart struct {
	widget.BaseWidget  // Inherit from BaseWidget
	dataPoints         []interfaces.SknDataSeries
	dataPointColor     color.Color
	dataPointScale     fyne.Size
	topLeftDesc        string // The text to display in the widget
	topCenteredDesc    string
	topRightDesc       string
	leftMiddleDesc     string
	rightMiddleDesc    string
	bottomCenteredDesc string
	bottomLeftDesc     string
	bottomRightDesc    string
}

var _ interfaces.SknLineChart = (*sknLineChart)(nil)

// Create a Widget and Extend (initialiase) the BaseWidget
func NewSknLineChart(title, xScale, yScale string, dataPoints []interfaces.SknDataSeries) *sknLineChart {
	w := &sknLineChart{ // Create this widget with an initial text value
		dataPoints:         dataPoints,
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
	}
	w.ExtendBaseWidget(w) // Initialize the BaseWidget
	return w
}

// Create the renderer. This is called by the fyne application
func (w *sknLineChart) CreateRenderer() fyne.WidgetRenderer {
	// Pass this widget to the renderer so it can access the text field
	return newSknLineChartRenderer(w)
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
func (w *sknLineChart) UpdateDataSeries(newData []interfaces.SknDataSeries) {
	for _, point := range newData {
		w.dataPoints = append(w.dataPoints, point)
	}
}

// Widget Renderer code starts here
type sknLineChartRenderer struct {
	widget             *sknLineChart     // Reference to the widget holding the current state
	background         *canvas.Rectangle // A background rectangle
	dataPointLines     []*canvas.Line
	xInc               float32
	yInc               float32
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
}

// Create the renderer with a reference to the widget
// Note: The background and foreground colours are set from the current theme.
//
// Do not size or move canvas objects here.
func newSknLineChartRenderer(lineChart *sknLineChart) *sknLineChartRenderer {
	background := canvas.NewRectangle(color.Transparent)
	background.StrokeWidth = 1
	background.StrokeColor = theme.PrimaryColorNamed(theme.ColorBlue)

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
		x.StrokeWidth = 0.5
		y := canvas.NewLine(theme.PrimaryColorNamed(theme.ColorGreen))
		y.StrokeWidth = 0.5
		xlines = append(xlines, x)
		ylines = append(ylines, y)
		xt := strconv.Itoa((10 - i) * 10)
		yt := strconv.Itoa((10 - i) * 10)
		xl := canvas.NewText(xt, theme.ForegroundColor())
		yl := canvas.NewText(yt, theme.ForegroundColor())
		xLabels = append(xLabels, xl)
		yLabels = append(yLabels, yl)
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
	leftMiddleDesc.Resize(fyne.NewSize(16, 200))
	rightMiddleDesc := canvas.NewText(lineChart.rightMiddleDesc, theme.ForegroundColor())
	rightMiddleDesc.TextSize = 18

	return &sknLineChartRenderer{
		widget:             lineChart,
		background:         background,
		xLines:             xlines,
		yLines:             ylines,
		xLabels:            xLabels,
		yLabels:            yLabels,
		topLeftDesc:        canvas.NewText(lineChart.topLeftDesc, theme.ForegroundColor()),
		topCenteredDesc:    topCenteredDesc,
		topRightDesc:       canvas.NewText(lineChart.topRightDesc, theme.ForegroundColor()),
		leftMiddleDesc:     leftMiddleDesc,
		rightMiddleDesc:    rightMiddleDesc,
		bottomLeftDesc:     canvas.NewText(lineChart.bottomLeftDesc, theme.ForegroundColor()),
		bottomCenteredDesc: bottomCenteredDesc,
		bottomRightDesc:    canvas.NewText(lineChart.bottomRightDesc, theme.ForegroundColor()),
	}
}

// The Refresh() method is called if the state of the widget changes or the
// theme is changed
//
// Note: The background and foreground colours are set from the current theme
func (r *sknLineChartRenderer) Refresh() {
	r.background.Refresh() // Redraw the background first
	for _, line := range r.xLines {
		line.Refresh()
	}
	for _, line := range r.yLines {
		line.Refresh()
	}
	for idx, xlbl := range r.xLabels {
		xlbl.Refresh()
		r.yLabels[idx].Refresh()
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

// Given the size required by the fyne application move and re-size the
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
		r.xLabels[idx+1].Move(fyne.NewPos((xp + (2 * r.xInc) - 8), yp+10))
		r.xLabels[idx+1].Refresh()
	}
	r.xLabels[0].Move(fyne.NewPos(((2 * r.xInc) - 10), yp+10))
	r.xLabels[0].Refresh()

	// Horiz lines
	xp := r.xInc
	for idx, line := range r.yLines {
		yp := float32(idx+1) * r.yInc

		line.Position1 = fyne.NewPos(xp, yp+r.yInc) // left
		line.Position2 = fyne.NewPos(xp*13+8, yp+r.yInc)
		line.Refresh()
		r.yLabels[idx+1].Move(fyne.NewPos(xp*13+10, (yp + r.yInc - 8)))
		r.yLabels[idx+1].Refresh()
	}
	r.yLabels[0].Move(fyne.NewPos(xp-10, (yp + 8)))
	r.yLabels[0].Refresh()

	// Make sure the background fills the widget
	r.background.Resize(fyne.NewSize(r.xInc*12, r.yInc*11))
	r.background.Move(fyne.NewPos(r.xInc, r.yInc))

	// Measure the size of the text so we can calculate the center offset.

	ts := fyne.MeasureText(r.topCenteredDesc.Text, r.topCenteredDesc.TextSize, r.topCenteredDesc.TextStyle)
	r.topCenteredDesc.Move(fyne.Position{X: (s.Width - ts.Width) / 2, Y: theme.Padding() / 4})

	ts = fyne.MeasureText(r.topRightDesc.Text, r.topRightDesc.TextSize, r.topRightDesc.TextStyle)
	r.topRightDesc.Move(fyne.Position{X: (s.Width - ts.Width) - theme.Padding(), Y: ts.Height / 4})
	r.topLeftDesc.Move(fyne.NewPos(theme.Padding(), (ts.Height / 4)))

	r.leftMiddleDesc.Move(fyne.NewPos(theme.Padding(), s.Height/3))
	ts = fyne.MeasureText(r.rightMiddleDesc.Text, r.rightMiddleDesc.TextSize, r.rightMiddleDesc.TextStyle)
	r.rightMiddleDesc.Move(fyne.NewPos((s.Width-ts.Width)-theme.Padding(), s.Height/3))

	ts = fyne.MeasureText(r.bottomCenteredDesc.Text, r.bottomCenteredDesc.TextSize, r.bottomCenteredDesc.TextStyle)
	r.bottomCenteredDesc.Move(fyne.Position{X: (s.Width - ts.Width) / 2, Y: (s.Height - ts.Height - theme.Padding())})

	ts = fyne.MeasureText(r.bottomRightDesc.Text, r.bottomRightDesc.TextSize, r.bottomRightDesc.TextStyle)
	r.bottomRightDesc.Move(fyne.NewPos((s.Width-ts.Width)-theme.Padding(), (s.Height - ts.Height - theme.Padding())))
	r.bottomLeftDesc.Move(fyne.NewPos(theme.Padding()+2.0, (s.Height - ts.Height - theme.Padding())))

}

// Create a minimum size for the widget.
// The smallest size is the size of the text with a border defined by the theme padding
func (r *sknLineChartRenderer) MinSize() fyne.Size {
	// Measure the size of the text so we can calculate a border size.
	ts := fyne.MeasureText(r.topCenteredDesc.Text, r.topCenteredDesc.TextSize, r.topCenteredDesc.TextStyle)
	// Use the theme padding to set a border size
	return fyne.NewSize(ts.Width+theme.Padding()*4, ts.Height+theme.Padding()*4)
}

// Return a list of each canvas object.
func (r *sknLineChartRenderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{
		r.background,
		r.topLeftDesc, r.topCenteredDesc, r.topRightDesc,
		r.leftMiddleDesc, r.rightMiddleDesc,
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

	return objs
}

// Cleanup if resources have been allocated
func (r *sknLineChartRenderer) Destroy() {}
