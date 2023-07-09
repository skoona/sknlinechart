package interfaces

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// SknLineChart feature list
type SknLineChart interface {
	desktop.Hoverable
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
	ReplaceDataSeries(newData *map[string][]SknDatapoint) error

	ApplyNewDataSeries(seriesName string, newSeries []SknDatapoint) error
	ApplySingleDataPoint(seriesName string, newDataPoint SknDatapoint)
}
