package bitmath_test

import (
	"math"
	"testing"

	"github.com/trphume/coachbuf/internal/bitmath"
)

func TestLog2For32(t *testing.T) {
	tests := []struct {
		name  string
		input uint32
		want  int
	}{
		{name: "smallest number input", input: 0, want: 0},
		{name: "log2 with small input", input: 8, want: 3},
		{name: "log2 with large input", input: 10000000, want: 23},
		{name: "largest number input", input: math.MaxUint32, want: 31},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := bitmath.Log2For32(tt.input)
			if got != tt.want {
				t.Errorf("Log2For32() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Exported (global) variable to store function results
// during benchmarking to ensure side effect free calls
// are not optimized away.
var Output int

func BenchmarkLog2For32(b *testing.B) {
	var tmp int
	for i := 0; i < b.N; i++ {
		tmp = bitmath.Log2For32(uint32(i))
	}
	Output = tmp
}
