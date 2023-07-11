package components

import "fyne.io/fyne/v2"

// SknLineChart feature list
type SknLineChart interface {
	IsDataPointMarkersEnabled() bool // mouse button 2 toggles
	IsHorizGridLinesEnabled() bool
	IsVertGridLinesEnabled() bool
	IsMousePointDisplayEnabled() bool // hoverable and mouse button one

	SetDataPointMarkers(enable bool)
	SetHorizGridLines(enable bool)
	SetVertGridLines(enable bool)
	SetMousePointDisplay(enable bool)

	GetTopLeftLabel() string
	GetTitle() string
	GetTopRightLabel() string
	GetMiddleLeftLabel() string
	GetMiddleRightLabel() string
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

	SetMinSize(s fyne.Size)
	Refresh()
	Resize(s fyne.Size)

	ApplyDataSeries(seriesName string, newSeries []SknChartDatapoint) error
	ApplyDataPoint(seriesName string, newDataPoint SknChartDatapoint)
}
