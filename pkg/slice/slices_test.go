package slice

import (
	"testing"
)

func TestSameValues(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want bool
	}{
		{
			name: "same values and order",
			a:    []int{1, 2, 3},
			b:    []int{1, 2, 3},
			want: true,
		},
		{
			name: "same values different order",
			a:    []int{1, 2, 3},
			b:    []int{3, 2, 1},
			want: true,
		},
		{
			name: "duplicate values match",
			a:    []int{1, 2, 2, 3},
			b:    []int{2, 1, 3, 2},
			want: true,
		},
		{
			name: "different lengths",
			a:    []int{1, 2, 3},
			b:    []int{1, 2},
			want: false,
		},
		{
			name: "different values",
			a:    []int{1, 2, 3},
			b:    []int{4, 5, 6},
			want: false,
		},
		{
			name: "one empty slice",
			a:    []int{},
			b:    []int{1},
			want: false,
		},
		{
			name: "both empty slices",
			a:    []int{},
			b:    []int{},
			want: true,
		},
		{
			name: "extra duplicates in one slice",
			a:    []int{1, 2, 3, 3},
			b:    []int{1, 2, 3},
			want: false,
		},
		{
			name: "repeated value mismatch",
			a:    []int{1, 1, 2, 3},
			b:    []int{1, 2, 3, 3},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SameValues(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("SameValues(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}
