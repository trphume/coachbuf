package coachbuf

import "errors"

var (
	// ErrUnsupportedType indicates that coachbuf has not implemented support for encoding the value of the type yet
	ErrUnsupportedType = errors.New("unsupported type")

	// ErrInvalidTagFormat indicates that the struct tag value received does not conform to the expected format
	ErrInvalidTagFormat = errors.New("invalid tag format")

	// ErrDuplicateOrdering indicates that the ordering number tag is given to more than one struct field
	ErrDuplicateOrdering = errors.New("given ordering number is used for more than one key")

	// ErrOutOfRangeOrdering indicates that the ordering number tag is not within the accepted min and max range
	ErrOutOfRangeOrdering = errors.New("given ordering number is not within accepted range")

	// ErrWriterInvalidState indicates that bitpacker.Writer is in an invalid state and could not continue the requested operation
	// This implies that there is a bug in coachbuf
	ErrWriterInvalidState = errors.New("invalid writer state")
)
