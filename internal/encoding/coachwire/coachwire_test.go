package coachwire_test

import (
	"bytes"
	"errors"
	"math"
	"testing"

	"github.com/trphume/coachbuf/internal/bitpacker"
	"github.com/trphume/coachbuf/internal/encoding/coachwire"
)

func TestWriteInteger(t *testing.T) {
	tests := []struct {
		name       string
		inputValue int32
		inputMin   int32
		inputMax   int32
		err        error
		want       []byte
	}{
		{name: "min > max", inputValue: 50, inputMin: 100, inputMax: 0, err: coachwire.ErrInvalidArgument, want: nil},
		{name: "value > max", inputValue: 10000, inputMin: 0, inputMax: 100, err: coachwire.ErrInvalidArgument, want: nil},
		{name: "value < min", inputValue: -100, inputMin: 0, inputMax: 100, err: coachwire.ErrInvalidArgument, want: nil},
		{name: "min == max", inputValue: 0, inputMin: 0, inputMax: 0, err: coachwire.ErrInvalidArgument, want: nil},
		{name: "valid input", inputValue: 50, inputMin: 0, inputMax: 100, err: nil, want: []byte{50, 0, 0, 0}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := bitpacker.NewWriter()

			if err := coachwire.WriteInteger(w, tt.inputValue, tt.inputMin, tt.inputMax); !errors.Is(err, tt.err) {
				t.Errorf("WriteInteger() = %v, want %v", err.Error(), tt.err.Error())
			}
			if err := w.FlushBits(); err != nil {
				t.Errorf("FlushBit() = %v, want %v", err.Error(), nil)
			}

			result := w.Bytes()
			if string(result) != string(tt.want) {
				t.Errorf("Bytes() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestReadInteger(t *testing.T) {
	tests := []struct {
		name     string
		reader   *bytes.Reader
		numBytes int
		inputMin int32
		inputMax int32
		err      error
		want     int32
	}{
		{name: "min > max", reader: bytes.NewReader([]byte{255, 0, 0, 0}), numBytes: 1, inputMin: 255, inputMax: 0, err: coachwire.ErrInvalidArgument, want: 0},
		{name: "min == max", reader: bytes.NewReader([]byte{255, 0, 0, 0}), numBytes: 1, inputMin: 255, inputMax: 255, err: coachwire.ErrInvalidArgument, want: 0},
		{name: "positive min", reader: bytes.NewReader([]byte{200, 0, 0, 0}), numBytes: 1, inputMin: 10, inputMax: 255, err: nil, want: 210},
		{name: "negative min", reader: bytes.NewReader([]byte{50, 0, 0, 0}), numBytes: 1, inputMin: -10, inputMax: 100, err: nil, want: 40},
		{name: "negative min and max", reader: bytes.NewReader([]byte{50, 0, 0, 0}), numBytes: 1, inputMin: -100, inputMax: -10, err: nil, want: -50},
		{name: "overflow wrap around", reader: bytes.NewReader([]byte{255, 255, 255, 255}), numBytes: 4, inputMin: math.MinInt32, inputMax: math.MaxInt32, err: nil, want: math.MaxInt32},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := bitpacker.NewReader(tt.reader, tt.numBytes)

			result, err := coachwire.ReadInteger(r, tt.inputMin, tt.inputMax)
			if !errors.Is(err, tt.err) {
				t.Errorf("ReadInteger() = %v, want %v", err, tt.err)
			}
			if tt.want != result {
				t.Errorf("ReadInteger() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestWriteAndReadInteger(t *testing.T) {
	tests := []struct {
		name     string
		value    int32
		inputMin int32
		inputMax int32
	}{
		{name: "positive min", value: 50000, inputMin: 10000, inputMax: 100000},
		{name: "negative max", value: -50000, inputMin: -100000, inputMax: -1000},
		{name: "negative min and positive max", value: 0, inputMin: -1000000, inputMax: 1000000},
		{name: "smallest value", value: math.MinInt32, inputMin: math.MinInt32, inputMax: 0},
		{name: "biggest value", value: math.MaxInt32, inputMin: 0, inputMax: math.MaxInt32},
		{name: "overflow wrap around", value: math.MaxInt32 - 100, inputMin: math.MinInt32, inputMax: math.MaxInt32},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// write
			w := bitpacker.NewWriter()
			if err := coachwire.WriteInteger(w, tt.value, tt.inputMin, tt.inputMax); err != nil {
				t.Errorf("WriteInteger() = %v, want %v", err, nil)
			}
			if err := w.FlushBits(); err != nil {
				t.Errorf("FlushBits() = %v, want %v", err, nil)
			}

			b := w.Bytes()

			// read
			r := bitpacker.NewReader(bytes.NewReader(b), len(b))
			result, err := coachwire.ReadInteger(r, tt.inputMin, tt.inputMax)
			if err != nil {
				t.Errorf("ReadInteger() = %v, want %v", err, nil)
			}

			if result != tt.value {
				t.Errorf("WriteInteger() and ReadInteger() = %v, want %v", result, tt.value)
			}
		})
	}
}

func TestWriteAndReadFloat(t *testing.T) {
	tests := []struct {
		name  string
		value float32
	}{
		{name: "negative value", value: -123.33},
		{name: "positive value", value: 424359.4349},
		{name: "max value", value: math.MaxFloat32},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// write
			w := bitpacker.NewWriter()
			if err := coachwire.WriteFloat(w, tt.value); err != nil {
				t.Errorf("WriteFloat() = %v, want %v", err, nil)
			}
			if err := w.FlushBits(); err != nil {
				t.Errorf("FlushBits() = %v, want %v", err, nil)
			}

			b := w.Bytes()

			// read
			r := bitpacker.NewReader(bytes.NewReader(b), len(b))
			result, err := coachwire.ReadFloat(r)
			if err != nil {
				t.Errorf("ReadFloat() = %v, want %v", err, nil)
			}

			if result != tt.value {
				t.Errorf("WriteFloat() and ReadFloat() = %v, want %v", result, tt.value)
			}
		})
	}
}

func TestWriteCompressedFloat(t *testing.T) {
	tests := []struct {
		name       string
		inputValue float32
		inputMin   float32
		inputMax   float32
		inputRes   float32
		err        error
		want       []byte
	}{
		{name: "min > max", inputValue: 100, inputMin: 100000, inputMax: 50, inputRes: 0.01, err: coachwire.ErrInvalidArgument, want: nil},
		{name: "value > max", inputValue: 5000000, inputMin: 50, inputMax: 100000, inputRes: 0.01, err: coachwire.ErrInvalidArgument, want: nil},
		{name: "value < min", inputValue: 100, inputMin: 100000, inputMax: 5000000, inputRes: 0.01, err: coachwire.ErrInvalidArgument, want: nil},
		{name: "min == max", inputValue: 100, inputMin: 100, inputMax: 100, inputRes: 0.01, err: coachwire.ErrInvalidArgument, want: nil},
		{name: "valid input", inputValue: 5000.12345, inputMin: 1000.54321, inputMax: 10000.54321, inputRes: 0.001, err: nil, want: []byte{92, 7, 61, 0}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := bitpacker.NewWriter()

			if err := coachwire.WriteCompressedFloat(w, tt.inputValue, tt.inputMin, tt.inputMax, tt.inputRes); !errors.Is(err, tt.err) {
				t.Errorf("WriteCompressedFloat() = %v, want %v", err.Error(), tt.err.Error())
			}
			if err := w.FlushBits(); err != nil {
				t.Errorf("FlushBit() = %v, want %v", err.Error(), nil)
			}

			result := w.Bytes()
			if string(result) != string(tt.want) {
				t.Errorf("Bytes() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestReadCompressedFloat(t *testing.T) {
	tests := []struct {
		name     string
		reader   *bytes.Reader
		numBytes int
		inputMin float32
		inputMax float32
		inputRes float32
		err      error
		want     float32
	}{
		{name: "min > max", reader: bytes.NewReader([]byte{0, 0, 0, 0}), numBytes: 1, inputMin: 10000, inputMax: 1000, inputRes: 0.01, err: coachwire.ErrInvalidArgument, want: 0},
		{name: "min == max", reader: bytes.NewReader([]byte{0, 0, 0, 0}), numBytes: 1, inputMin: 10000, inputMax: 10000, inputRes: 0.01, err: coachwire.ErrInvalidArgument, want: 0},
		{name: "valid all positive", reader: bytes.NewReader([]byte{100, 0, 0, 0}), numBytes: 1, inputMin: 0, inputMax: 2, inputRes: 0.01, err: nil, want: 1},
		{name: "valid all negative", reader: bytes.NewReader([]byte{100, 0, 0, 0}), numBytes: 1, inputMin: -2, inputMax: 0, inputRes: 0.01, err: nil, want: -1},
		{name: "valid negative min and positive max", reader: bytes.NewReader([]byte{238, 2, 0, 0}), numBytes: 2, inputMin: -5, inputMax: 5, inputRes: 0.01, err: nil, want: 2.5},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := bitpacker.NewReader(tt.reader, tt.numBytes)

			result, err := coachwire.ReadCompressedFloat(r, tt.inputMin, tt.inputMax, tt.inputRes)
			if !errors.Is(err, tt.err) {
				t.Errorf("ReadCompressedFloat() = %v, want %v", err, tt.err)
			}
			if tt.want != result {
				t.Errorf("ReadCompressedFloat() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestWriteAndReadCompressedFloat(t *testing.T) {
	tests := []struct {
		name     string
		value    float32
		inputMin float32
		inputMax float32
		inputRes float32
	}{
		{name: "all positive", value: 8, inputMin: 1, inputMax: 9, inputRes: 0.001},
		{name: "all negative", value: -8, inputMin: -9, inputMax: -1, inputRes: 0.01},
		{name: "negative min and positive max", value: 3, inputMin: -9, inputMax: 9, inputRes: 0.001},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// write
			w := bitpacker.NewWriter()
			if err := coachwire.WriteCompressedFloat(w, tt.value, tt.inputMin, tt.inputMax, tt.inputRes); err != nil {
				t.Errorf("WriteCompressedFloat() = %v, want %v", err, nil)
			}
			if err := w.FlushBits(); err != nil {
				t.Errorf("FlushBits() = %v, want %v", err, nil)
			}

			b := w.Bytes()

			// read
			r := bitpacker.NewReader(bytes.NewReader(b), len(b))
			result, err := coachwire.ReadCompressedFloat(r, tt.inputMin, tt.inputMax, tt.inputRes)
			if err != nil {
				t.Errorf("ReadCompressedFloat() = %v, want %v", err, nil)
			}

			// we need to account for imprecision at higher res and number range thus we find difference between input and result
			if result != tt.value && math.Abs(float64(tt.value)-float64(result)) > 0.001 {
				t.Errorf("WriteCompressedFloat() and ReadFloat() = %v, want %v", result, tt.value)
			}
		})
	}
}
