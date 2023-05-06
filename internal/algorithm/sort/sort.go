// Package algorithmsort includes implementation for sorting algorithms utilizing generics
//
// Note: we could easily use the sorting algorithm in Go sorting package but for the sake of learning this exist
// this was originally going to be used to order the values during encoding but realised that is not necessary
// i've kept it here because a lot of effort was put into writing pdq sort :(
package algorithmsort

import (
	"math/bits"
)

const (
	PDQMaxInsert     = 12
	PDQNintherNumber = 40
)

// PDQSort is a minimal sorting algorithm implementation (not as complete) based on pattern-defeating quicksort
// paper: https://arxiv.org/pdf/2106.05123.pdf
// presentation: https://youtu.be/jz-PBiWwNjc
func PDQSort[T any](s []T, less func(T, T) bool) {
	// FYI the calculation method for fallback limit counter is taken from Go's "sort" package
	n := len(s)
	if n <= 1 {
		return
	}
	limit := bits.Len(uint(n))
	pdqSort(s, less, limit)
}

func pdqSort[T any](s []T, less func(T, T) bool, limit int) {
	n := len(s)

	// for small inputs we use insertion sort
	if n <= PDQMaxInsert {
		InsertionSort(s, less)
		return
	}

	// if too many bad pivots were made we fall back to heapsort
	if limit <= 0 {
		HeapSort(s, less)
		return
	}

	pivot := selectPivot(s, less)
	mid := partition(s, pivot, less)

	pdqSort(s[:mid+1], less, limit)
	pdqSort(s[mid+1:], less, limit)

	// data patterns:
	//
	// low cardinality:
	// 1. compare on element before current subarray
	// 2. if equal part-left (less than equal to on the left)  + recurse right otherwise part-right (more than equal to on the right) + recurse both
	//
	// pre-sorted or mostly sorted:
	// 1. check if perfect partitioning (already partitioned correctly according to the pivot)
	// 2. if perfect do optimistic insertion sort (if more than 8 element is out of place then abort insertion sort)
	//
	// self-similarity and malicious inputs (repeated bad pivot selection -> results in bad partition):
	// 1. check if bad partition (less than 10% of array is partitioned to one side)
	// 2. increment bad partition counter and check is counter more than log2(n)
	// 3. if more than we use fallback sort
	// 4. if not we need to introduce some new pivots
}

// assumption is that parameter s has more elements than PDQMaxInsert
func selectPivot[T any](s []T, less func(T, T) bool) int {
	n := len(s)
	i, j, k := 0, n/2, n-1
	if n > PDQNintherNumber {
		partSize := n / 6
		i = medianOfThree(s, 0, partSize, 2*partSize-1, less)
		j = medianOfThree(s, 2*partSize, 3*partSize, 4*partSize-1, less)
		k = medianOfThree(s, 4*partSize, 5*partSize, n-1, less)
	}

	return medianOfThree(s, i, j, k, less)
}

func medianOfThree[T any](s []T, i, j, k int, less func(T, T) bool) int {
	if less(s[i], s[j]) {
		switch {
		case less(s[j], s[k]):
			return j // i < j < k
		case less(s[i], s[k]):
			return k // i < k < j
		default:
			return i // k < i < j
		}
	} else {
		switch {
		case less(s[k], s[j]):
			return j // k < j < i
		case less(s[k], s[i]):
			return k // j < k < i
		default:
			return i // j < i < k
		}
	}
}

// partition uses hoare partitioning algorithm
// assumption is that parameter s has more elements than PDQMaxInsert
func partition[T any](s []T, pivot int, less func(T, T) bool) int {
	s[0], s[pivot] = s[pivot], s[0]
	pivotData := s[0]

	i, j := 0, len(s)-1
	for {
		for i <= j && less(s[i], pivotData) {
			i++
		}
		for i <= j && !less(s[j], pivotData) {
			j--
		}

		if i >= j {
			break
		}
		s[i], s[j] = s[j], s[i]
		i++
		j--
	}

	return j
}

// InsertionSort is an implementation of insertion sort algorithm
func InsertionSort[T any](s []T, less func(T, T) bool) {
	if len(s) <= 1 {
		return
	}

	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && less(s[j], s[j-1]); j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

// HeapSort is an implementation heap sort algorithm
// to construct a heap in O(n) we use floyd's algorithm to shift the elements down
// article: https://copyprogramming.com/howto/buildheap-the-algorithm-of-floyd
func HeapSort[T any](s []T, less func(T, T) bool) {
	n := len(s)

	// build heap
	for i := n/2 - 1; i >= 0; i-- {
		heapify(s, i, less)
	}

	// sorting the heap
	for i := len(s) - 1; i > 0; i-- {
		s[0], s[i] = s[i], s[0]
		heapify(s[:i], 0, less)
	}
}

func heapify[T any](s []T, i int, less func(T, T) bool) {
	n := len(s)
	for {
		largest := i
		left := 2*i + 1
		right := 2*i + 2

		if left < n && less(s[largest], s[left]) {
			largest = left
		}

		if right < n && less(s[largest], s[right]) {
			largest = right
		}

		if largest == i {
			break
		}

		s[i], s[largest] = s[largest], s[i]
		i = largest
	}
}
