package coachbuf_test

import (
	"errors"
	"testing"

	"github.com/trphume/coachbuf"
)

func TestEncodeDecodeStruct(t *testing.T) {
	t.Run("Encode and Decode", func(t *testing.T) {
		t.Parallel()

		t.Run("simple", func(t *testing.T) {
			t.Parallel()

			type TestStruct struct {
				Int32   int32   `coachbuf:"1,min=100"`
				Float32 float32 `coachbuf:"32,max=30000"`
				String  string  // unsupported type but not tagged with coachbuf
			}

			inputEncode := TestStruct{Int32: 10000, Float32: 10000.34, String: "Hello"}
			inputDecode, err := coachbuf.Encode(inputEncode)
			if err != nil {
				t.Errorf("Encode() = %v, want %v", err.Error(), nil)
			}

			result := TestStruct{}
			if err := coachbuf.Decode(inputDecode, &result); err != nil {
				t.Errorf("Decode() = %v, want %v", err.Error(), nil)
			}
			switch {
			case inputEncode.Int32 != result.Int32:
				t.Errorf("Decode() = %v, want %v", result.Int32, inputEncode.Int32)
			case inputEncode.Float32 != result.Float32:
				t.Errorf("Decode() = %v, want %v", result.Int32, inputEncode.Int32)
			}
		})

		t.Run("nested", func(t *testing.T) {
			t.Parallel()

			type TestStruct struct {
				Int32   int32   `coachbuf:"1,min=0,max=1000000"`
				Float32 float32 `coachbuf:"32"`
				String  string  // unsupported type but not tagged with coachbuf
				Nested  struct {
					Int32 int32 `coachbuf:"1"`
				} `coachbuf:"100"`
			}

			inputEncode := TestStruct{
				Int32:   10000,
				Float32: 10000.34,
				String:  "Hello",
				Nested: struct {
					Int32 int32 `coachbuf:"1"`
				}{Int32: 1000000},
			}
			inputDecode, err := coachbuf.Encode(inputEncode)
			if err != nil {
				t.Errorf("Encode() = %v, want %v", err.Error(), nil)
			}

			result := TestStruct{}
			if err := coachbuf.Decode(inputDecode, &result); err != nil {
				t.Errorf("Decode() = %v, want %v", err.Error(), nil)
			}
			switch {
			case inputEncode.Int32 != result.Int32:
				t.Errorf("Decode() = %v, want %v", result.Int32, inputEncode.Int32)
			case inputEncode.Float32 != result.Float32:
				t.Errorf("Decode() = %v, want %v", result.Int32, inputEncode.Int32)
			case inputEncode.Nested.Int32 != result.Nested.Int32:
				t.Errorf("Decode() = %v, want %v", result.Int32, inputEncode.Int32)
			}
		})
	})

	t.Run("Encode", func(t *testing.T) {
		t.Parallel()

		t.Run("invalid ordering number tag format", func(t *testing.T) {
			t.Parallel()

			input := struct {
				Int32   int32   `coachbuf:"1"`
				Float32 float32 `coachbuf:"min=100"`
				String  string  // unsupported type but not tagged with coachbuf
			}{Int32: 10000, Float32: 10000.34, String: "Hello"}
			want := coachbuf.ErrInvalidTagFormat

			_, err := coachbuf.Encode(input)
			if !errors.Is(err, want) {
				t.Errorf("Encode() = %v, want %v", err.Error(), want)
			}
		})

		t.Run("ordering number out of range", func(t *testing.T) {
			t.Parallel()

			input := struct {
				Int32   int32   `coachbuf:"1"`
				Float32 float32 `coachbuf:"100000000"`
				String  string  // unsupported type but not tagged with coachbuf
			}{Int32: 10000, Float32: 10000.34, String: "Hello"}
			want := coachbuf.ErrOutOfRangeOrdering

			_, err := coachbuf.Encode(input)
			if !errors.Is(err, want) {
				t.Errorf("Encode() = %v, want %v", err.Error(), want)
			}
		})

		t.Run("duplicate ordering number", func(t *testing.T) {
			t.Parallel()

			input := struct {
				Int32   int32   `coachbuf:"1"`
				Float32 float32 `coachbuf:"1"`
				String  string  // unsupported type but not tagged with coachbuf
			}{Int32: 10000, Float32: 10000.34, String: "Hello"}
			want := coachbuf.ErrDuplicateOrdering

			_, err := coachbuf.Encode(input)
			if !errors.Is(err, want) {
				t.Errorf("Encode() = %v, want %v", err.Error(), want)
			}
		})

		t.Run("min/max tag missing value", func(t *testing.T) {
			t.Parallel()

			input := struct {
				Int32   int32   `coachbuf:"1,min="`
				Float32 float32 `coachbuf:"2"`
				String  string  // unsupported type but not tagged with coachbuf
			}{Int32: 10000, Float32: 10000.34, String: "Hello"}
			want := coachbuf.ErrInvalidTagFormat

			_, err := coachbuf.Encode(input)
			if !errors.Is(err, want) {
				t.Errorf("Encode() = %v, want %v", err.Error(), want)
			}
		})

		t.Run("min/max tag not a number", func(t *testing.T) {
			t.Parallel()

			input := struct {
				Int32   int32   `coachbuf:"1,max=string"`
				Float32 float32 `coachbuf:"2"`
				String  string  // unsupported type but not tagged with coachbuf
			}{Int32: 10000, Float32: 10000.34, String: "Hello"}
			want := coachbuf.ErrInvalidTagFormat

			_, err := coachbuf.Encode(input)
			if !errors.Is(err, want) {
				t.Errorf("Encode() = %v, want %v", err.Error(), want)
			}
		})

		t.Run("unsupported type", func(t *testing.T) {
			t.Parallel()

			input := struct {
				Int32   int32   `coachbuf:"1"`
				Float32 float32 `coachbuf:"100"`
				String  string  `coachbuf:"2"`
			}{Int32: 10000, Float32: 10000.34, String: "Hello"}
			want := coachbuf.ErrUnsupportedType

			_, err := coachbuf.Encode(input)
			if !errors.Is(err, want) {
				t.Errorf("Encode() = %v, want %v", err.Error(), want)
			}
		})
	})
}
