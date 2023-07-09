package interfaces

type SknDataSeries interface {
	Value() float32
	Timestamp() string
	ColorName() string
	SetValue(y float32)
	SetColorName(n string)
	SetTimestamp(t string)
}
