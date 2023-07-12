package sknlinechart

// RemoveIndexFromSlice remove the given index from any type of slice
func RemoveIndexFromSlice[K comparable](index int, slice []K) []K {
	var idx int

	if index > len(slice) {
		idx = len(slice) - 1
	} else if index < 0 {
		idx = 0
	} else {
		idx = index
	}
	return append(slice[:idx], slice[idx+1:]...)
}

// ShiftSlice drops index 0 and append newData to any type of slice
func ShiftSlice[K comparable](newData K, slice []K) []K {
	idx := 0

	shorter := append(slice[:idx], slice[idx+1:]...)
	shorter = append(shorter, newData)
	return shorter
}
