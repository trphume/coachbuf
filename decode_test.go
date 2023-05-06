package coachbuf_test

import (
	"errors"
	"testing"

	"github.com/trphume/coachbuf"
)

// TestDecode uses encode to encode data to its binary representation - we do not test the internal detail
func TestDecode(t *testing.T) {
	t.Run("Int32", func(t *testing.T) {
		t.Parallel()
		value := int32(100)

		inputData, err := coachbuf.Encode(value)
		if err != nil {
			t.Errorf("Encode() = %v, want %v", err.Error(), nil)
		}

		var inputDecode int32
		if err = coachbuf.Decode(inputData, &inputDecode); err != nil {
			t.Errorf("Decode() = %v, want %v", err.Error(), nil)
		}

		if inputDecode != value {
			t.Errorf("Decode() = %v, want %v", inputDecode, value)
		}
	})

	t.Run("Float32", func(t *testing.T) {
		t.Parallel()
		value := float32(123.123)

		inputData, err := coachbuf.Encode(value)
		if err != nil {
			t.Errorf("Encode() = %v, want %v", err.Error(), nil)
		}

		var inputDecode float32
		if err = coachbuf.Decode(inputData, &inputDecode); err != nil {
			t.Errorf("Decode() = %v, want %v", err.Error(), nil)
		}

		if inputDecode != value {
			t.Errorf("Decode() = %v, want %v", inputDecode, value)
		}
	})

	t.Run("non pointer v argument", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Deode() expect panic for bad v argument")
			}
		}()

		inputDecode := int32(100)
		_ = coachbuf.Decode([]byte{255, 255, 255, 255}, inputDecode)
	})

	t.Run("nil pointer v argument", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Deode() expect panic for bad v argument")
			}
		}()

		_ = coachbuf.Decode([]byte{255, 255, 255, 255}, nil)
	})

	t.Run("Unsupported type", func(t *testing.T) {
		t.Parallel()
		value := int32(100)

		inputData, err := coachbuf.Encode(value)
		if err != nil {
			t.Errorf("Encode() = %v, want %v", err.Error(), nil)
		}

		var inputDecode string
		if err = coachbuf.Decode(inputData, &inputDecode); !errors.Is(err, coachbuf.ErrUnsupportedType) {
			t.Errorf("Decode() = %v, want %v", err.Error(), coachbuf.ErrUnsupportedType)
		}
	})
}
