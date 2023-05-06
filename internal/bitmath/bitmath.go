// Package bitmath provides math operation that utilises bit operations for fast computations
package bitmath

//go:generate go run make_tables.go

// Log2For32 finds the log base 2 of an uint32 integer with a lookup table
// Algorithm reference: https://graphics.stanford.edu/~seander/bithacks.html#IntegerLogLookup
func Log2For32(x uint32) int {
	var res int

	if idx1 := x >> 16; idx1 != 0 {
		idx2 := idx1 >> 8
		if idx2 != 0 {
			res = 24 + int(log2tab[idx2])
		} else {
			res = 16 + int(log2tab[idx1])
		}
	} else {
		idx := x >> 8
		if idx != 0 {
			res = 8 + int(log2tab[idx])
		} else {
			res = int(log2tab[x])
		}
	}

	return res
}
