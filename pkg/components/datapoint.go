package components

import (
	"fyne.io/fyne/v2"
)

// SknChartDatapoint data container interface for SknLineChart
type SknChartDatapoint interface {
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

type sknDatapoint struct {
	value                float32
	markerTopPosition    *fyne.Position
	markerBottomPosition *fyne.Position
	colorName            string
	timestamp            string
}

func NewSknDatapoint(value float32, colorName, timestamp string) SknChartDatapoint {
	return &sknDatapoint{
		value:                value,
		colorName:            colorName,
		timestamp:            timestamp,
		markerTopPosition:    &fyne.Position{},
		markerBottomPosition: &fyne.Position{},
	}
}
func (d *sknDatapoint) Value() float32 {
	return d.value
}
func (d *sknDatapoint) MarkerPosition() (*fyne.Position, *fyne.Position) {
	return d.markerTopPosition, d.markerBottomPosition
}
func (d *sknDatapoint) ColorName() string {
	return d.colorName
}
func (d *sknDatapoint) Timestamp() string {
	return d.timestamp
}
func (d *sknDatapoint) SetValue(v float32) {
	(*d).value = v
}
func (d *sknDatapoint) SetMarkerPosition(top *fyne.Position, bottom *fyne.Position) {
	(*d).markerTopPosition = top
	(*d).markerBottomPosition = bottom
}
func (d *sknDatapoint) SetColorName(n string) {
	(*d).colorName = n
}
func (d *sknDatapoint) SetTimestamp(t string) {
	(*d).timestamp = t
}
