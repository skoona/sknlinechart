package interfaces

import "image/color"

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
	SetTopRightDescription(newValue string)
	SetMiddleLeftDescription(newValue string)
	SetMiddleRightDescription(newValue string)
	SetBottomLeftDescription(newValue string)
	SetBottomRightDescription(newValue string)
	SetBottomCenteredDescription(newValue string)
	SetDataSeriesColor(c color.Color)
	UpdateDataSeries(newData []SknDataSeries)
}
