package sknlinechart

import "fyne.io/fyne/v2"

// LineChart feature list
type LineChart interface {
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

	ApplyDataSeries(seriesName string, newSeries []LineChartDatapoint) error
	ApplyDataPoint(seriesName string, newDataPoint LineChartDatapoint)
}
