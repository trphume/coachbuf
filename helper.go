package coachbuf

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// getCocahbufTag extracts the values from a comma separated string in coachbuf expected format
// return the tag values, the ordering number, and an error
func getCoachbufTag(structFieldName string, structTag string) ([]string, int32, error) {
	if structTag == "" {
		return nil, 0, nil
	}

	cbStructTags := strings.Split(structTag, ",")
	order, err := strconv.ParseInt(cbStructTags[0], 10, 32)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"ordering number must be first value in the comma separated tag, field=%s: %w",
			structFieldName, ErrInvalidTagFormat,
		)
	}

	if order < cbMinOrderingNumber || order > cbMaxOrderingNumber {
		return nil, 0, fmt.Errorf("order=%d: %w", order, ErrOutOfRangeOrdering)
	}

	return cbStructTags, int32(order), nil
}

// getMinAndMaxTags is a helper function to find and retrieve min and max tags from a slice of string
// return values min, max, err in this order
func getMinAndMaxTags(tags []string) (int32, int32, error) {
	min, max := int32(math.MinInt32), int32(math.MaxInt32)

	var minSet, maxSet bool
	for _, tag := range tags {
		if strings.HasPrefix(tag, "max=") || strings.HasPrefix(tag, "min=") {
			if minSet && maxSet {
				break
			}
			if len(tag) <= 4 {
				return min, max, fmt.Errorf("min and max tag value missing value, tag=%s: %w", tag, ErrInvalidTagFormat)
			}

			if value, err := strconv.ParseInt(tag[4:], 10, 32); err == nil {
				switch tag[:4] {
				case "max=":
					max = int32(value)
					maxSet = true
				case "min=":
					min = int32(value)
					minSet = true
				}
			} else {
				return min, max, fmt.Errorf("min and max tag value must be a int32 number, tag=%s: %w", tag, ErrInvalidTagFormat)
			}
		}
	}

	return min, max, nil
}
