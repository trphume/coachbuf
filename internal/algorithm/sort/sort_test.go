package algorithmsort_test

import (
	"testing"

	"github.com/trphume/coachbuf/internal/algorithm/sort"
)

func TestPDQSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
		want  []int
	}{
		{name: "small input (insertion sort only)", input: []int{37, 17, -7, 7, 57, 67, 47, 27}, want: []int{-7, 7, 17, 27, 37, 47, 57, 67}},
		{
			name:  "requires partition (not small input)",
			input: []int{-47, 37, 17, -7, 7, 57, -37, 67, 47, 27, -17, -27, 77, 87, 97, 107},
			want:  []int{-47, -37, -27, -17, -7, 7, 17, 27, 37, 47, 57, 67, 77, 87, 97, 107},
		},
		{
			name: "large input (tukey's ninther)",
			input: []int{147, 157, -87, -47, 37, -117, 17, 127, -197, -207, -177, -7, -57, 57, -67, 7, -37, 67, 47, 137, 27, -17,
				-27, 77, 87, 97, -77, -97, 107, -107, -137, -147, -127, 167, -157, -167, -187, 117, 207, 177, 187, 197,
			},
			want: []int{-207, -197, -187, -177, -167, -157, -147, -137, -127, -117, -107, -97, -87, -77, -67, -57, -47, -37, -27, -17, -7,
				7, 17, 27, 37, 47, 57, 67, 77, 87, 97, 107, 117, 127, 137, 147, 157, 167, 177, 187, 197, 207,
			},
		},
		{
			name:  "already sorted",
			input: []int{-27, -17, -7, 7, 17, 27, 37, 47, 57, 67, 77, 87, 97, 107},
			want:  []int{-27, -17, -7, 7, 17, 27, 37, 47, 57, 67, 77, 87, 97, 107},
		},
		{name: "one element", input: []int{9999}, want: []int{9999}},
		{name: "empty slice", input: []int{}, want: []int{}},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := make([]int, len(tt.input))
			copy(input, tt.input)
			algorithmsort.PDQSort(tt.input, func(a, b int) bool { return a < b })
			shouldEqual(t, "PDQSort(): got %v, want %v", tt.input, tt.want)
		})
	}
}

func TestInsertionSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
		want  []int
	}{
		{name: "all positive", input: []int{5, 2, 6, 1, 10}, want: []int{1, 2, 5, 6, 10}},
		{name: "all negative", input: []int{-5, -2, -6, -1, -10}, want: []int{-10, -6, -5, -2, -1}},
		{name: "positive and negative", input: []int{5, -3, 4, -2, 10}, want: []int{-3, -2, 4, 5, 10}},
		{name: "already sorted", input: []int{1, 2, 3, 4, 5}, want: []int{1, 2, 3, 4, 5}},
		{name: "sorted in reverse", input: []int{5, 4, 3, 2, 1}, want: []int{1, 2, 3, 4, 5}},
		{name: "one element", input: []int{9999}, want: []int{9999}},
		{name: "empty slice", input: []int{}, want: []int{}},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := make([]int, len(tt.input))
			copy(input, tt.input)
			algorithmsort.InsertionSort(tt.input, func(a, b int) bool { return a < b })
			shouldEqual(t, "InsertionSort(): got %v, want %v", tt.input, tt.want)
		})
	}
}

func TestHeapSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
		want  []int
	}{
		{name: "all positive", input: []int{5, 2, 6, 1, 10}, want: []int{1, 2, 5, 6, 10}},
		{name: "all negative", input: []int{-5, -2, -6, -1, -10}, want: []int{-10, -6, -5, -2, -1}},
		{name: "positive and negative", input: []int{5, -3, 4, -2, 10}, want: []int{-3, -2, 4, 5, 10}},
		{name: "already sorted", input: []int{1, 2, 3, 4, 5}, want: []int{1, 2, 3, 4, 5}},
		{name: "sorted in reverse", input: []int{5, 4, 3, 2, 1}, want: []int{1, 2, 3, 4, 5}},
		{name: "one element", input: []int{9999}, want: []int{9999}},
		{name: "empty slice", input: []int{}, want: []int{}},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := make([]int, len(tt.input))
			copy(input, tt.input)
			algorithmsort.HeapSort(tt.input, func(a, b int) bool { return a < b })
			shouldEqual(t, "InsertionSort(): got %v, want %v", tt.input, tt.want)
		})
	}
}

// shouldEqual is a testing helper function that checks if two slices are equal
func shouldEqual[T comparable](t testing.TB, errString string, a, b []T) {
	if len(a) != len(b) {
		t.Errorf(errString, a, b)
		return
	}

	for i := range a {
		if a[i] != b[i] {
			t.Errorf(errString, a, b)
			return
		}
	}
}
