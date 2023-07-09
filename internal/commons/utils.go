package commons

// RemoveIndexFromSlice()
// -- No good way to handle various types of slices
//
//	since reflection does not support custom types
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

func ShiftSlice[K comparable](newData K, slice []K) []K {
	idx := 0

	shorter := append(slice[:idx], slice[idx+1:]...)
	shorter = append(shorter, newData)
	return shorter
}

// MapKeys returns a slice og the key in a map
//
// []K := MapKeys(m)
// []K := MapKeys[int, string](m)
func MapKeys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
