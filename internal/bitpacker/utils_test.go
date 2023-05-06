package bitpacker_test

import (
	"math"
	"testing"

	"github.com/trphume/coachbuf/internal/bitpacker"
)

func TestBitsRequired1(t *testing.T) {
	tests := []struct {
		name  string
		input uint32
		want  int
	}{
		{name: "smallest number input", input: 0, want: 0},
		{name: "normal number input", input: 256, want: 9},
		{name: "largest number input", input: math.MaxUint32, want: 32},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := bitpacker.BitsRequired(tt.input); got != tt.want {
				t.Errorf("BitsRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Exported (global) variable to store function results
// during benchmarking to ensure side effect free calls
// are not optimized away.
var Output int

func BenchmarkBitsRequired(b *testing.B) {
	var tmp int
	for i := 0; i < b.N; i++ {
		tmp = bitpacker.BitsRequired(uint32(i))
	}

	Output = tmp
}
