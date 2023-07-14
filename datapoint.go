package sknlinechart

import (
	"fyne.io/fyne/v2"
	"github.com/google/uuid"
	"strings"
)

type chartDatapoint struct {
	value                float32
	colorName            string
	timestamp            string
	externalID           string
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
		externalID:           uuid.New().String(),
	}
}
func (d *chartDatapoint) Copy() ChartDatapoint {
	return &chartDatapoint{
		value:      d.value,
		colorName:  strings.Clone(d.colorName),
		timestamp:  strings.Clone(d.timestamp),
		externalID: strings.Clone(d.externalID),
	}
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
func (d *chartDatapoint) ExternalID() string {
	return d.externalID
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
