package dst_test

import (
	"testing"

	"github.com/trphume/coachbuf/internal/dst"
)

func TestCopyMap(t *testing.T) {
	original := map[int]string{1: "One", 2: "Two", 3: "Three"}
	result := dst.CopyMap(original)

	if len(result) != len(original) {
		t.Errorf("CopyMap() = %v, want %v", len(result), len(original))
	}
	for k, originalValue := range original {
		if copiedValue, exist := result[k]; exist && copiedValue != originalValue {
			// exist but not equal to original
			t.Errorf("CopyMap() value not equal = %v, want %v", copiedValue, originalValue)
		} else if !exist {
			// does not exist
			t.Errorf("CopyMap() key does not exist: want %v", k)
		}
	}

	result[1] = "Not One"
	if original[1] == result[1] {
		t.Errorf("result map mutation affected original map")
	}
}

func TestOrderedKeyValueLessFunc(t *testing.T) {
	tests := []struct {
		name   string
		input1 dst.OrderedKeyValue[int, int]
		input2 dst.OrderedKeyValue[int, int]
		want   bool
	}{
		{
			name:   "input1 < input2",
			input1: dst.OrderedKeyValue[int, int]{Key: -1, Value: 0},
			input2: dst.OrderedKeyValue[int, int]{Key: 0, Value: 0},
			want:   true,
		},
		{name: "input1 > input2",
			input1: dst.OrderedKeyValue[int, int]{Key: 1, Value: 0},
			input2: dst.OrderedKeyValue[int, int]{Key: 0, Value: 0},
			want:   false,
		},
		{name: "input1 == input2",
			input1: dst.OrderedKeyValue[int, int]{Key: 0, Value: 0},
			input2: dst.OrderedKeyValue[int, int]{Key: 0, Value: 0},
			want:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := dst.OrderedKeyValueLessFunc(tt.input1, tt.input2)
			if result != tt.want {
				t.Errorf("OrderedKeyValueLessFunc() = %v, want = %v", result, tt.want)
			}
		})
	}
}
