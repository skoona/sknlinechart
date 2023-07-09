package entities

import (
	"fyne.io/fyne/v2"
	"ggApcMon/internal/interfaces"
)

type sknDatapoint struct {
	value                float32
	markerTopPosition    *fyne.Position
	markerBottomPosition *fyne.Position
	colorName            string
	timestamp            string
}

func NewSknDatapoint(value float32, colorName, timestamp string) interfaces.SknDatapoint {
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
