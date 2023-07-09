package interfaces

import (
	"fyne.io/fyne/v2"
	"image/color"
)

type SknLineChart interface {
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
	SetMinSize(s fyne.Size)
	SetTopRightDescription(newValue string)
	SetMiddleLeftDescription(newValue string)
	SetMiddleRightDescription(newValue string)
	SetBottomLeftDescription(newValue string)
	SetBottomRightDescription(newValue string)
	SetBottomCenteredDescription(newValue string)
	SetDataSeriesColor(c color.Color)
	ReplaceDataSeries(newData *map[string][]SknDataSeries) error
	ApplyNewDataSeries(seriesName string, newSeries []SknDataSeries) error
	ApplySingleDataPoint(seriesName string, newDataPoint SknDataSeries)
}
