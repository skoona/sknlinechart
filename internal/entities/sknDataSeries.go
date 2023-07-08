package entities

import "ggApcMon/internal/interfaces"

type sknDataSeries struct {
	xValue    float32
	yValue    float32
	timestamp string
}

func NewSknDataSeries(xValue, yValue float32, timestamp string) interfaces.SknDataSeries {
	return &sknDataSeries{
		xValue:    xValue,
		yValue:    yValue,
		timestamp: timestamp,
	}
}
func (d sknDataSeries) YValue() float32 {
	return d.yValue
}
func (d sknDataSeries) XValue() float32 {
	return d.xValue
}
func (d sknDataSeries) Timestamp() string {
	return d.timestamp
}
func (d sknDataSeries) SetYValue(y float32) {
	d.yValue = y
}
func (d sknDataSeries) SetXValue(x float32) {
	d.xValue = x
}
func (d sknDataSeries) SetTimestamp(t string) {
	d.timestamp = t
}
