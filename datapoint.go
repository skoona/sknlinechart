package sknlinechart

import (
	"fyne.io/fyne/v2"
)

// LineChartDatapoint data container interface for LineChart
type LineChartDatapoint interface {
	Value() float32
	SetValue(y float32)

	ColorName() string
	SetColorName(n string)

	Timestamp() string
	SetTimestamp(t string)

	// MarkerPosition internal use only: current data point marker location
	MarkerPosition() (*fyne.Position, *fyne.Position)
	// SetMarkerPosition internal use only: sace location of where data point marker is located
	SetMarkerPosition(top *fyne.Position, bottom *fyne.Position)
}

type lineChartDatapoint struct {
	value                float32
	markerTopPosition    *fyne.Position
	markerBottomPosition *fyne.Position
	colorName            string
	timestamp            string
}

func NewLineChartDatapoint(value float32, colorName, timestamp string) LineChartDatapoint {
	return &lineChartDatapoint{
		value:                value,
		colorName:            colorName,
		timestamp:            timestamp,
		markerTopPosition:    &fyne.Position{},
		markerBottomPosition: &fyne.Position{},
	}
}
func (d *lineChartDatapoint) Value() float32 {
	return d.value
}
func (d *lineChartDatapoint) MarkerPosition() (*fyne.Position, *fyne.Position) {
	return d.markerTopPosition, d.markerBottomPosition
}
func (d *lineChartDatapoint) ColorName() string {
	return d.colorName
}
func (d *lineChartDatapoint) Timestamp() string {
	return d.timestamp
}
func (d *lineChartDatapoint) SetValue(v float32) {
	(*d).value = v
}
func (d *lineChartDatapoint) SetMarkerPosition(top *fyne.Position, bottom *fyne.Position) {
	(*d).markerTopPosition = top
	(*d).markerBottomPosition = bottom
}
func (d *lineChartDatapoint) SetColorName(n string) {
	(*d).colorName = n
}
func (d *lineChartDatapoint) SetTimestamp(t string) {
	(*d).timestamp = t
}
