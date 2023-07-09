package interfaces

import "fyne.io/fyne/v2"

type SknDatapoint interface {
	Value() float32
	MarkerPosition() (*fyne.Position, *fyne.Position)
	Timestamp() string
	ColorName() string
	SetValue(y float32)
	SetMarkerPosition(top *fyne.Position, bottom *fyne.Position)
	SetColorName(n string)
	SetTimestamp(t string)
}
