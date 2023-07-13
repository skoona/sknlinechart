package sknlinechart

import (
	"fyne.io/fyne/v2"
)

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
	d.value = v
}
func (d *lineChartDatapoint) SetMarkerPosition(top *fyne.Position, bottom *fyne.Position) {
	d.markerTopPosition = top
	d.markerBottomPosition = bottom
}
func (d *lineChartDatapoint) SetColorName(n string) {
	d.colorName = n
}
func (d *lineChartDatapoint) SetTimestamp(t string) {
	d.timestamp = t
}
