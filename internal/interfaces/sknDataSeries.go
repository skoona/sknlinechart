package interfaces

type SknDataSeries interface {
	XValue() float32
	YValue() float32
	Timestamp() string
	SetXValue(x float32)
	SetYValue(y float32)
	SetTimestamp(t string)
}
