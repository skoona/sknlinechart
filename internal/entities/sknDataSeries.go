package entities

import "ggApcMon/internal/interfaces"

type sknDataSeries struct {
	value     float32
	colorName string
	timestamp string
}

func NewSknDataSeries(value float32, colorName, timestamp string) interfaces.SknDataSeries {
	return &sknDataSeries{
		value:     value,
		colorName: colorName,
		timestamp: timestamp,
	}
}
func (d *sknDataSeries) Value() float32 {
	return d.value
}
func (d *sknDataSeries) ColorName() string {
	return d.colorName
}
func (d *sknDataSeries) Timestamp() string {
	return d.timestamp
}
func (d *sknDataSeries) SetValue(v float32) {
	(*d).value = v
}
func (d *sknDataSeries) SetColorName(n string) {
	(*d).colorName = n
}
func (d *sknDataSeries) SetTimestamp(t string) {
	(*d).timestamp = t
}
