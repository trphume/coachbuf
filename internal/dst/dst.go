// Package dst implements common data structures and related functions and methods
package dst

import "github.com/trphume/coachbuf/constraints"

// CopyMap is a helper function that copies the value of a map to a new map (that won't affect the original)
func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	res := make(map[K]V)
	for k, v := range m {
		res[k] = v
	}

	return res
}

// OrderedKeyValue is a data structure where the key is an ordered type
type OrderedKeyValue[K constraints.Ordered, V any] struct {
	Key   K
	Value V
}

// OrderedKeyValueLessFunc returns true if the first parameter is less than the second
func OrderedKeyValueLessFunc[K constraints.Ordered, V any](v1, v2 OrderedKeyValue[K, V]) bool {
	return v1.Key < v2.Key
}
