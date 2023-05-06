// Package coachwire provides operations to write and read values into wire (binary) format
//
// It wraps around the basic bitpacker package read and write functions
package coachwire

import (
	"errors"
	"fmt"
	"math"

	"github.com/trphume/coachbuf/internal/bitpacker"
)

// WriteInteger and ReadInteger are meant to be used together
// Assumptions made by ReadInteger regarding overflow are only valid for buffer written with WriteInteger

// WriteInteger writes an integer in the specified range [min, max] where min != max
func WriteInteger(writer *bitpacker.Writer, value, min, max int32) error {
	if min > max || value > max || value < min || min == max {
		return fmt.Errorf("value=%d, min=%d, max=%d: %w", value, min, max, ErrInvalidArgument)
	}
	bits := bitpacker.BitsRequired(uint32(max - min))
	unsignedValue := uint32(value - min)

	err := writer.Write(unsignedValue, bits)
	if err != nil {
		if errors.Is(err, bitpacker.ErrBitsInvalidRange) {
			panic("required bits error")
		}

		return err
	}

	return nil
}

// ReadInteger reads an integer in the specified range [min, max] where min != max
func ReadInteger(reader *bitpacker.Reader, min, max int32) (int32, error) {
	if min > max || min == max {
		return 0, fmt.Errorf("min=%d, max=%d: %w", min, max, ErrInvalidArgument)
	}
	bits := bitpacker.BitsRequired(uint32(max - min))

	unsignedValue, err := reader.Read(bits)
	if err != nil {
		if errors.Is(err, bitpacker.ErrBitsInvalidRange) {
			panic("required bits error")
		}

		return 0, err
	}

	value := int32(unsignedValue) + min

	return value, nil
}

// WriteFloat and ReadFloat are meant to be used together
// Assumptions made by ReadFloat regarding overflow are only valid for buffer written with WriteFloat

// WriteFloat writes float32 value with 32 bits as-is for full precision
func WriteFloat(writer *bitpacker.Writer, value float32) error {
	if err := writer.Write(math.Float32bits(value), 32); err != nil {
		if errors.Is(err, bitpacker.ErrBitsInvalidRange) {
			panic("required bits error")
		}

		return err
	}

	return nil
}

// ReadFloat float32  reads float32 value with 32 bits for full precision
func ReadFloat(reader *bitpacker.Reader) (float32, error) {
	value, err := reader.Read(32)
	if err != nil {
		if errors.Is(err, bitpacker.ErrBitsInvalidRange) {
			panic("required bits error")
		}

		return 0, err
	}

	return math.Float32frombits(value), nil
}

// WriteCompressedFloat and ReadCompressedFloat are meant to be used together
// Assumptions made by ReadCompressedFloat regarding overflow are only valid for buffer written with WriteCompressedFloat

// WriteCompressedFloat writes a float value given a value,min,max and a precision res (example res: 0.001)
// The function is suited for small-mid range values and precision
// There is a limit to precision value chosen but this function does not handle all edge cases - use within reason
func WriteCompressedFloat(writer *bitpacker.Writer, value, min, max, res float32) error {
	if min > max || value > max || value < min || min == max {
		return fmt.Errorf("value=%f, min=%f, max=%f: %w", value, min, max, ErrInvalidArgument)
	}

	diff := max - min
	maxIntegerValue := math.Ceil(float64(diff / res))
	normalizedValue := float64((value - min) / diff)
	switch {
	case normalizedValue < 0:
		normalizedValue = 0
	case normalizedValue > 1:
		normalizedValue = 1
	}
	integerValue := uint32(math.Floor(normalizedValue*maxIntegerValue + 0.5))

	bits := bitpacker.BitsRequired(uint32(maxIntegerValue))
	if err := writer.Write(integerValue, bits); err != nil {
		if errors.Is(err, bitpacker.ErrBitsInvalidRange) {
			panic("required bits error")
		}

		return err
	}

	return nil
}

// ReadCompressedFloat reads a compressed float32 value
func ReadCompressedFloat(reader *bitpacker.Reader, min, max, res float32) (float32, error) {
	if min > max || min == max {
		return 0, fmt.Errorf("min=%f, max=%f: %w", min, max, ErrInvalidArgument)
	}

	diff := max - min
	maxIntegerValue := math.Ceil(float64(diff / res))

	bits := bitpacker.BitsRequired(uint32(maxIntegerValue))
	integerValue, err := reader.Read(bits)
	if err != nil {
		if errors.Is(err, bitpacker.ErrBitsInvalidRange) {
			panic("required bits error")
		}

		return 0, err
	}

	normalizedValue := float64(integerValue) / maxIntegerValue
	value := float32(normalizedValue)*diff + min
	return value, nil
}
