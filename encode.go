package coachbuf

import (
	"fmt"
	"math"
	"reflect"

	"github.com/trphume/coachbuf/internal/bitpacker"
	"github.com/trphume/coachbuf/internal/encoding/coachwire"
)

// Encode takes in a value and serializes it into a slice of byte in Coachbuf format
func Encode(v any) ([]byte, error) {
	writer := bitpacker.NewWriter()
	rv := reflect.ValueOf(v)
	if err := encodeValue(writer, rv); err != nil {
		return nil, err
	}

	if err := writer.FlushBits(); err != nil {
		return nil, fmt.Errorf("could not FlushBits() after writing: %w", ErrWriterInvalidState)
	}

	return writer.Bytes(), nil
}

func encodeValue(writer *bitpacker.Writer, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Struct:
		return encodeStruct(writer, rv)
	case reflect.Int32:
		return coachwire.WriteInteger(writer, int32(rv.Int()), math.MinInt32, math.MaxInt32)
	case reflect.Float32:
		return coachwire.WriteFloat(writer, float32(rv.Float()))
	default:
		return fmt.Errorf("encode type=%v: %w", rv.Type(), ErrUnsupportedType)
	}
}

func encodeStruct(writer *bitpacker.Writer, rv reflect.Value) error {
	rt := rv.Type()

	orderExistMap := make(map[int]struct{}) // is needed to keep track of ordering number (to avoid duplicate)
	for i := 0; i < rt.NumField(); i++ {
		structField := rt.Field(i)
		structTag := structField.Tag.Get(cbStructTagsKey)

		cbStructTags, order, err := getCoachbufTag(structField.Name, structTag)
		if err != nil {
			return err
		}

		if len(cbStructTags) < 1 {
			continue
		}
		if _, exist := orderExistMap[int(order)]; exist {
			return fmt.Errorf("field=%s: %w", structField.Name, ErrDuplicateOrdering)
		}
		orderExistMap[int(order)] = struct{}{}

		if err = coachwire.WriteInteger(writer, order, cbMinOrderingNumber, cbMaxOrderingNumber); err != nil {
			return fmt.Errorf("field=%s: %w", structField.Name, err)
		}

		fieldValue := rv.Field(i)
		switch fieldValue.Kind() {
		case reflect.Struct:
			err = encodeStruct(writer, fieldValue)
		case reflect.Int32:
			var min, max int32
			min, max, err = getMinAndMaxTags(cbStructTags)
			if err != nil {
				break
			}
			err = coachwire.WriteInteger(writer, int32(fieldValue.Int()), min, max)
		case reflect.Float32:
			err = coachwire.WriteFloat(writer, float32(fieldValue.Float()))
		default:
			err = fmt.Errorf("encode type=%v: %w", fieldValue.Type(), ErrUnsupportedType)
		}
		if err != nil {
			return fmt.Errorf("field=%s: %w", structField.Name, err)
		}

	}

	return nil
}
