package sknlinechart

import "fyne.io/fyne/v2"

// LineChartDatapoint data container interface for SknLineChart
type LineChartDatapoint interface {
	Value() float32
	SetValue(y float32)

	ColorName() string
	SetColorName(n string)

	Timestamp() string
	SetTimestamp(t string)

	// MarkerPosition internal use only: current data point marker location
	MarkerPosition() (*fyne.Position, *fyne.Position)

	// SetMarkerPosition internal use only: screen location of where data point marker is located
	SetMarkerPosition(top *fyne.Position, bottom *fyne.Position)
}

// SknLineChart feature list
type SknLineChart interface {
	// Chart Attributes

	IsDataPointMarkersEnabled() bool // mouse button 2 toggles
	IsHorizGridLinesEnabled() bool
	IsVertGridLinesEnabled() bool
	IsMousePointDisplayEnabled() bool // hoverable and mouse button one

	SetDataPointMarkers(enable bool)
	SetHorizGridLines(enable bool)
	SetVertGridLines(enable bool)
	SetMousePointDisplay(enable bool)

	// Info labels

	GetTopLeftLabel() string
	GetTitle() string
	GetTopRightLabel() string

	// Scale legend

	GetMiddleLeftLabel() string
	GetMiddleRightLabel() string

	// Info Labels

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
	// expect this will rarely be used, since loading more than 120 point will raise error
	ApplyDataSeries(seriesName string, newSeries []*LineChartDatapoint) error

	// ApplyDataPoint primary method to add another data point to any series
	// If series has more than 120 points, point 0 will be rolled out making room for this one
	ApplyDataPoint(seriesName string, newDataPoint *LineChartDatapoint)

	// SetMinSize sets the minimum widget size respond when asked
	SetMinSize(s fyne.Size)

	// EnableDebugLogging turns method entry/exit logging on or off
	EnableDebugLogging(enable bool)

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
