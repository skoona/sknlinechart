package sknlinechart

import "fyne.io/fyne/v2"

// GraphPointSmoothing support for different implementation
// of averaging or smooth data; current provides rolling average from last x reading.
type GraphPointSmoothing interface {
	AddValue(value float64) float64
	SeriesName() string
	IsNil() bool
	String() string
}

// ChartDatapoint data container interface for LineChart
type ChartDatapoint interface {
	Value() float32
	SetValue(y float32)

	ColorName() string
	SetColorName(n string)

	Timestamp() string
	SetTimestamp(t string)

	// ExternalID string uuid assigned when created
	ExternalID() string

	// Copy returns a cloned copy of current item
	Copy() ChartDatapoint

	// MarkerPosition internal use only: current data point marker location
	MarkerPosition() (*fyne.Position, *fyne.Position)

	// SetMarkerPosition internal use only: screen location of where data point marker is located
	SetMarkerPosition(top *fyne.Position, bottom *fyne.Position)
}

// LineChart feature list
type LineChart interface {
	// Chart Attributes

	GetLineStrokeSize() float32
	SetLineStrokeSize(newSize float32)

	IsDataPointMarkersEnabled() bool // mouse button 2 toggles
	IsHorizGridLinesEnabled() bool
	IsVertGridLinesEnabled() bool
	IsColorLegendEnabled() bool
	IsMousePointDisplayEnabled() bool // hoverable and mouse button one

	SetDataPointMarkers(enable bool)
	SetHorizGridLines(enable bool)
	SetVertGridLines(enable bool)
	SetColorLegend(enable bool)
	SetMousePointDisplay(enable bool)

	// Scale legend

	GetMiddleLeftLabel() string
	GetMiddleRightLabel() string

	// Info Labels

	GetTopLeftLabel() string
	GetTitle() string
	GetTopRightLabel() string
	GetBottomLeftLabel() string
	GetBottomCenteredLabel() string
	GetBottomRightLabel() string

	SetTopLeftLabel(newValue string)
	SetTitle(newValue string)
	SetTopRightLabel(newValue string)
	SetMiddleLeftLabel(newValue string)
	SetMiddleRightLabel(newValue string)
	SetBottomLeftLabel(newValue string)
	SetBottomCenteredLabel(newValue string)
	SetBottomRightLabel(newValue string)

	// ApplyDataSeries add a whole data series at once
	// expect this will rarely be used, since loading more than 130 point will raise error
	ApplyDataSeries(seriesName string, newSeries []*ChartDatapoint) error

	// ApplyDataPoint primary method to add another data point to any series
	// If series has more than 130 points, point 0 will be rolled out making room for this one
	ApplyDataPoint(seriesName string, newDataPoint *ChartDatapoint)

	// SetMinSize set the minimum size limit for the linechart
	SetMinSize(s fyne.Size)

	// EnableDebugLogging turns method entry/exit logging on or off
	EnableDebugLogging(enable bool)

	// SetHoverPointCallback method to call when a onscreen datapoint is hovered over by pointer
	SetOnHoverPointCallback(func(series string, dataPoint ChartDatapoint))

	// ObjectCount internal use only: return the default ui elements for testing
	ObjectCount() int

	// fyne.CanvasObject compliance
	// implemented by BaseWidget
	Hide()
	MinSize() fyne.Size
	Move(position fyne.Position)
	Position() fyne.Position
	Refresh()
	Resize(size fyne.Size)
	Show()
	Size() fyne.Size
	Visible() bool
}
