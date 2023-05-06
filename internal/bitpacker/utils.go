package bitpacker

import "github.com/trphume/coachbuf/internal/bitmath"

// BitsRequired calculate number of bits required to represent an uint32 number
func BitsRequired(x uint32) int {
	if x == 0 {
		return 0
	}

	return bitmath.Log2For32(x) + 1
}
