package slice

import "cmp"

// SameValues - check if two slices have the same values
func SameValues[T cmp.Ordered](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	aCountMap := map[T]int{}
	for _, v := range a {
		aCountMap[v]++
	}
	for _, v := range b {
		if aCountMap[v] == 0 {
			return false
		}
		aCountMap[v]--
	}

	return true
}
