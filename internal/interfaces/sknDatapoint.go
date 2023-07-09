package interfaces

import "fyne.io/fyne/v2"

// SknDatapoint data container interface for SknLineChart
type SknDatapoint interface {
	Value() float32
	SetValue(y float32)

	ColorName() string
	SetColorName(n string)

	Timestamp() string
	SetTimestamp(t string)

	// mouse hover popup info; internal use only
	MarkerPosition() (*fyne.Position, *fyne.Position)
	SetMarkerPosition(top *fyne.Position, bottom *fyne.Position)
}
