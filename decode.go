package coachbuf

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"reflect"

	"github.com/trphume/coachbuf/internal/bitpacker"
	"github.com/trphume/coachbuf/internal/encoding/coachwire"
)

// Decode takes in a value and deserializes it into the value v
// argument v must be a non-nil pointer
func Decode(data []byte, v any) error {
	reader := bitpacker.NewReader(bytes.NewReader(data), len(data))
	pointerRv := reflect.ValueOf(v)
	if pointerRv.Kind() != reflect.Pointer || pointerRv.IsNil() {
		panic("argument v must be non-nil pointer type")
	}

	return decodeValue(reader, pointerRv)
}

func decodeValue(reader *bitpacker.Reader, pointerRv reflect.Value) error {
	rv := reflect.Indirect(pointerRv)
	switch rv.Kind() {
	case reflect.Struct:
		return decodeStruct(reader, rv)
	case reflect.Int32:
		if v, err := coachwire.ReadInteger(reader, math.MinInt32, math.MaxInt32); err != nil {
			return err
		} else {
			if !rv.CanSet() {
				panic("cannot set value")
			}
			rv.SetInt(int64(v))
			return nil
		}
	case reflect.Float32:
		if v, err := coachwire.ReadFloat(reader); err != nil {
			return err
		} else {
			if !rv.CanSet() {
				panic("cannot set value")
			}
			rv.SetFloat(float64(v))
			return nil
		}
	default:
		return fmt.Errorf("decode type=%v: %w", rv.Type(), ErrUnsupportedType)
	}
}

func decodeStruct(reader *bitpacker.Reader, rv reflect.Value) error {
	type orderingToFieldValue struct {
		fieldNumber  int
		cbStructTags []string
	}

	rt := rv.Type()

	orderingToFieldMap := make(map[int32]orderingToFieldValue)
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

		if _, exist := orderingToFieldMap[order]; exist {
			return fmt.Errorf("field=%s: %w", structField.Name, ErrDuplicateOrdering)
		}

		orderingToFieldMap[order] = orderingToFieldValue{fieldNumber: i, cbStructTags: cbStructTags}
	}

	// start reading in the ordering number from the reader then find how to read via the map previously built
	var readCounter int
	for readCounter < len(orderingToFieldMap) {
		order, err := coachwire.ReadInteger(reader, cbMinOrderingNumber, cbMaxOrderingNumber)
		if err != nil {
			return fmt.Errorf("error reading ordering number: %w", err)
		}

		otmValue := orderingToFieldMap[order]
		fieldValue := rv.Field(otmValue.fieldNumber)
		switch fieldValue.Kind() {
		case reflect.Struct:
			err = decodeStruct(reader, fieldValue)
		case reflect.Int32:
			var min, max, v int32
			min, max, err = getMinAndMaxTags(otmValue.cbStructTags)
			if err != nil {
				break
			}

			v, err = coachwire.ReadInteger(reader, min, max)
			log.Println("Hellllooooo", v)
			if err == nil {
				if !fieldValue.CanSet() {
					panic("cannot set value")
				}
				fieldValue.SetInt(int64(v))
			}
		case reflect.Float32:
			var v float32
			v, err = coachwire.ReadFloat(reader)
			if err == nil {
				if !fieldValue.CanSet() {
					panic("cannot set value")
				}
				fieldValue.SetFloat(float64(v))
			}
		default:
			err = fmt.Errorf("encode type=%v: %w", fieldValue.Type(), ErrUnsupportedType)
		}
		if err != nil {
			return fmt.Errorf("field=%s: %w", rt.Field(otmValue.fieldNumber).Name, err)
		}

		readCounter++
	}

	return nil
}
