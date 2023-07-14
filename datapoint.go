package sknlinechart

import (
	"fyne.io/fyne/v2"
	"strings"
)

type chartDatapoint struct {
	value                float32
	colorName            string
	timestamp            string
	markerTopPosition    *fyne.Position
	markerBottomPosition *fyne.Position
}

func NewChartDatapoint(value float32, colorName, timestamp string) ChartDatapoint {
	return &chartDatapoint{
		value:                value,
		colorName:            colorName,
		timestamp:            timestamp,
		markerTopPosition:    &fyne.Position{0, 0},
		markerBottomPosition: &fyne.Position{0, 0},
	}
}
func (d *chartDatapoint) Copy() ChartDatapoint {
	return NewChartDatapoint(d.value, strings.Clone(d.colorName), strings.Clone(d.timestamp))
}
func (d *chartDatapoint) Value() float32 {
	return d.value
}
func (d *chartDatapoint) MarkerPosition() (*fyne.Position, *fyne.Position) {
	return d.markerTopPosition, d.markerBottomPosition
}
func (d *chartDatapoint) ColorName() string {
	return d.colorName
}
func (d *chartDatapoint) Timestamp() string {
	return d.timestamp
}
func (d *chartDatapoint) SetValue(v float32) {
	d.value = v
}
func (d *chartDatapoint) SetMarkerPosition(top *fyne.Position, bottom *fyne.Position) {
	d.markerTopPosition = top
	d.markerBottomPosition = bottom
}
func (d *chartDatapoint) SetColorName(n string) {
	d.colorName = n
}
func (d *chartDatapoint) SetTimestamp(t string) {
	d.timestamp = t
}
