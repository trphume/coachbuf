package bitpacker

import "errors"

var (
	ErrBitsInvalidRange     = errors.New("number of bits provided is not in a valid range")
	ErrBitsReadExceeded     = errors.New("number of bits to read exceeded total number of bits specifies")
	ErrMethodCallNotAllowed = errors.New("calling this method is not allowed")
)
