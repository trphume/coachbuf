package coachbuf_test

import (
	"errors"
	"testing"

	"github.com/trphume/coachbuf"
)

// TestEncode doesn't test the correctness of resulting []bytes - the binary representation is an internal detail
func TestEncode(t *testing.T) {
	t.Run("Int32", func(t *testing.T) {
		t.Parallel()

		input := int32(255)
		_, err := coachbuf.Encode(input)
		if err != nil {
			t.Errorf("Encode() = %v, want %v", err.Error(), nil)
		}
	})

	t.Run("Float32", func(t *testing.T) {
		t.Parallel()

		input := float32(100.456)
		_, err := coachbuf.Encode(input)
		if err != nil {
			t.Errorf("Encode() = %v, want %v", err.Error(), nil)
		}
	})

	t.Run("Unsupported type", func(t *testing.T) {
		t.Parallel()

		input := "Hello"
		want := coachbuf.ErrUnsupportedType

		_, err := coachbuf.Encode(input)
		if !errors.Is(err, want) {
			t.Errorf("Encode() = %v, want %v", err.Error(), want)
		}
	})
}
