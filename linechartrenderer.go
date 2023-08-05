package sknlinechart

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"math"
	"strconv"
	"strings"
	"time"
)

// Widget Renderer code starts here
type lineChartRenderer struct {
	widget                *LineChartSkn // Reference to the widget holding the current state
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
	lineChart.mapsLock.Lock()
	defer lineChart.mapsLock.Unlock()

	var (
		dataPoints       = map[string][]*canvas.Line{}
		dpMaker          = map[string][]*canvas.Circle{}
		objs             []fyne.CanvasObject
		xlines, ylines   []*canvas.Line
		xLabels, yLabels []*canvas.Text
	)

	// hover frame
	border := canvas.NewRectangle(theme.OverlayBackgroundColor())
	border.StrokeColor = theme.PrimaryColorNamed(lineChart.mouseDisplayFrameColor)
	border.StrokeWidth = 2.0

	// hover content
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

	// x & y frame lines
	for i := 0; i < 16; i++ { // vertical
		x := canvas.NewLine(theme.PrimaryColorNamed(theme.ColorGreen))
		x.StrokeWidth = 0.25
		xlines = append(xlines, x)
		objs = append(objs, x)
	}
	for i := 0; i < 14; i++ { // horiz line
		y := canvas.NewLine(theme.PrimaryColorNamed(theme.ColorGreen))
		y.StrokeWidth = 0.25
		ylines = append(ylines, y)
		objs = append(objs, y)
	}

	// Y scale labels
	for i := 0; i < 14; i++ {
		yt := strconv.Itoa((13 - i) * lineChart.chartScaleMultiplier)
		yl := canvas.NewText(yt, theme.ForegroundColor())
		yl.Alignment = fyne.TextAlignTrailing
		yLabels = append(yLabels, yl)
		objs = append(objs, yl)
	}
	// X scale labels
	for i := 0; i < 16; i++ {
		xt := strconv.Itoa(i * 10)
		xl := canvas.NewText(xt, theme.ForegroundColor())
		xl.Alignment = fyne.TextAlignTrailing
		xLabels = append(xLabels, xl)
		objs = append(objs, xl)
	}

	// series legend on bottom right
	colorLegend := container.NewHBox()
	strokeSize := lineChart.dataPointStrokeSize
	markerSize := strokeSize * 5
	for key, points := range lineChart.dataPoints {
		for _, point := range points {
			x := canvas.NewLine(theme.PrimaryColorNamed((*point).ColorName()))
			x.StrokeWidth = strokeSize
			dataPoints[key] = append(dataPoints[key], x)
			z := canvas.NewCircle(theme.PrimaryColorNamed((*point).ColorName()))
			z.StrokeWidth = strokeSize * 2
			z.Resize(fyne.NewSize(markerSize, markerSize))
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

	r.verifyDataPoints(true)

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

	r.widget.mapsLock.RLock()
	r.topLeftDesc.Text = r.widget.topLeftLabel
	r.topCenteredDesc.Text = r.widget.topCenteredLabel
	r.topRightDesc.Text = r.widget.topRightLabel
	r.bottomLeftDesc.Text = r.widget.bottomLeftLabel
	r.bottomCenteredDesc.Text = r.widget.bottomCenteredLabel
	r.bottomRightDesc.Text = r.widget.bottomRightLabel
	for _, v := range r.widget.objectsCache {
		v.Refresh()
	}

	r.manageLabelVisibility()

	r.widget.mapsLock.RUnlock()

	r.widget.mapsLock.Lock()

	r.mouseDisplayContainer.Hide()
	r.mouseDisplayContainer.Objects[0].(*canvas.Rectangle).StrokeColor = theme.PrimaryColorNamed(r.widget.mouseDisplayFrameColor)
	r.mouseDisplayContainer.Objects[1].(*widget.Label).SetText(r.widget.mouseDisplayStr)

	r.widget.mapsLock.Unlock()

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
	yp := r.yInc * 14.0
	yScale := (r.yInc * 10) / (10.0 * float32(r.widget.chartScaleMultiplier)) // 100
	xScale := (r.xInc * 10) / 100
	var dp float32
	data := r.widget.dataPoints[series] // datasource
	lastPoint := fyne.NewPos(xp, yp)

	for idx, point := range data { // one set of lines
		if (*point).Value() > r.widget.dataPointYLimit { // max y chart scale
			dp = r.widget.dataPointYLimit
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

	r.widget.mapsLock.Lock()
	defer r.widget.mapsLock.Unlock()

	r.xInc = (s.Width - (theme.Padding() * 4)) / 16.0
	r.yInc = (s.Height - (theme.Padding() * 3)) / 16.0

	r.xInc = float32(math.Trunc(float64(r.xInc)))
	r.yInc = float32(math.Trunc(float64(r.yInc)))

	// grid Vert lines
	yp := 14.0 * r.yInc
	for idx, line := range r.xLines {
		xp := float32(idx) * r.xInc
		line.Position1 = fyne.NewPos(xp+r.xInc, r.yInc) //top
		line.Position2 = fyne.NewPos(xp+r.xInc, yp+8)
	}

	// grid Horiz lines
	xp := r.xInc
	for idx, line := range r.yLines {
		yp := float32(idx) * r.yInc
		line.Position1 = fyne.NewPos(xp-8, yp+r.yInc) // left
		line.Position2 = fyne.NewPos(xp*16, yp+r.yInc)
	}

	// grid scale labels
	xp = r.xInc
	yp = 14.0 * r.yInc
	for idx, label := range r.xLabels {
		xxp := float32(idx+1) * r.xInc // starting at left
		label.Move(fyne.NewPos(xxp+8, yp+10))
	}
	for idx, label := range r.yLabels {
		yyp := float32(idx+1) * r.yInc // starting at top
		label.Move(fyne.NewPos(xp*0.80, yyp-8))
	}

	// handle new data points
	r.verifyDataPoints(false)

	// handle new data series
	if !r.widget.datapointAdded {
		for key := range r.widget.dataPoints { // datasource
			r.layoutSeries(key)
		}
	}
	r.widget.dataSeriesAdded = false
	r.widget.datapointAdded = false

	ts := fyne.MeasureText(
		r.topCenteredDesc.Text,
		r.topCenteredDesc.TextSize,
		r.topCenteredDesc.TextStyle)
	r.topCenteredDesc.Move(fyne.Position{X: (s.Width - ts.Width) / 2, Y: -4})

	ts = fyne.MeasureText(
		r.topRightDesc.Text,
		r.topRightDesc.TextSize,
		r.topRightDesc.TextStyle)
	r.topRightDesc.Move(fyne.Position{X: (s.Width - ts.Width) - theme.Padding(), Y: ts.Height / 4})
	r.topLeftDesc.Move(fyne.NewPos(theme.Padding(), ts.Height/4))

	msg := strings.Split(r.mouseDisplayContainer.Objects[1].(*widget.Label).Text, "[")
	ts = fyne.MeasureText(msg[0], 14, r.mouseDisplayContainer.Objects[1].(*widget.Label).TextStyle)
	//r.mouseDisplayContainer.Objects[1].(*widget.Label).Text = strings.Join(msg, "")
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

	r.widget.mapsLock.RLock()
	defer r.widget.mapsLock.RUnlock()

	var objs []fyne.CanvasObject
	objs = append(objs, r.widget.objectsCache...)

	for key, lines := range r.dataPoints {
		for idx, line := range lines {
			marker := r.dataPointMarkers[key][idx]
			objs = append(objs, marker, line)
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
func (r *lineChartRenderer) verifyDataPoints(protect bool) {
	startTime := time.Now()

	r.widget.debugLog("lineChartRenderer::VerifyDataPoints() ENTER")

	if protect {
		r.widget.mapsLock.Lock()
		defer r.widget.mapsLock.Unlock()
	}

	var changedKeys []string
	var changed bool
	strokeSize := r.widget.dataPointStrokeSize
	markerSize := strokeSize * 5
	for key, points := range r.widget.dataPoints {
		changed = false
		if nil == r.dataPoints[key] {
			r.dataPoints[key] = []*canvas.Line{}
			r.dataPointMarkers[key] = []*canvas.Circle{}
			changed = true
		}
		for idx, point := range points {
			if idx > (len(r.dataPoints[key]) - 1) { // add added points
				changed = true
				x := canvas.NewLine(theme.PrimaryColorNamed((*point).ColorName()))
				x.StrokeWidth = strokeSize
				r.dataPoints[key] = append(r.dataPoints[key], x)
				z := canvas.NewCircle(theme.PrimaryColorNamed((*point).ColorName()))
				z.StrokeWidth = strokeSize * 2
				z.Resize(fyne.NewSize(markerSize, markerSize))
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
