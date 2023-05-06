// Package bitpacker provides basic bitpacking operations
//
// This package operates in little endian byte order assuming a 32 bit word
package bitpacker

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Writer provides the ability to write some amount of bits to a buffer
type Writer struct {
	scratch        uint64
	scratchBits    int
	wordIndex      int
	buffer         *bytes.Buffer
	numBitsWritten int
	flushed        bool
}

// NewWriter returns a Writer with an empty buffer
func NewWriter() *Writer {
	return &Writer{buffer: new(bytes.Buffer)}
}

// NewWriterWithBuffer returns a Writer with a given buffer
func NewWriterWithBuffer(b *bytes.Buffer) *Writer {
	return &Writer{buffer: b}
}

// Write writes a binary value into the buffer given some desired number of bits to be written
func (w *Writer) Write(value uint32, bits int) error {
	if bits <= 0 || bits > 32 {
		return fmt.Errorf("bits should be in the range (0,32]: %w", ErrBitsInvalidRange)
	}
	if w.flushed {
		return fmt.Errorf("FlushBits() previously called: %w", ErrMethodCallNotAllowed)
	}

	// masking value and shifting into position
	value &= (uint32(1) << bits) - 1
	w.scratch |= uint64(value) << w.scratchBits
	w.scratchBits += bits

	for w.scratchBits >= 32 {
		if err := binary.Write(w.buffer, binary.LittleEndian, uint32(w.scratch&0xffffffff)); err != nil {
			return fmt.Errorf("could not write scratch with value %b to buffer: %w", w.scratch, err)
		}

		w.wordIndex++
		w.scratch >>= 32
		w.scratchBits -= 32
	}

	w.numBitsWritten += bits
	return nil
}

// FlushBits must be called ONLY once at the end to write any remaining value in scratch to the buffer
func (w *Writer) FlushBits() error {
	if w.flushed {
		return fmt.Errorf("FlushBits() previously called: %w", ErrMethodCallNotAllowed)
	}
	if w.scratchBits != 0 {
		err := binary.Write(w.buffer, binary.LittleEndian, uint32(w.scratch&0xffffffff))
		if err != nil {
			return fmt.Errorf("could not write scratch with value %b to buffer: %w", w.scratch, err)
		}

		w.wordIndex++
		w.scratch = 0
		w.scratchBits = 0
	}

	w.flushed = true
	return nil
}

// Bytes return a slice of bytes of value written to the buffer
func (w *Writer) Bytes() []byte {
	return w.buffer.Bytes()
}

func (w *Writer) NumBitsWritten() int {
	return w.numBitsWritten
}

// Reader provides the ability to read some amount bits from a buffer
type Reader struct {
	scratch     uint64
	scratchBits int
	totalBits   int
	numBitsRead int
	reader      *bytes.Reader
}

func NewReader(r *bytes.Reader, numBytes int) *Reader {
	if r.Len()%4 != 0 {
		panic("the reader needs to fit a 32 bit word")
	}
	return &Reader{
		totalBits: numBytes * 8,
		reader:    r,
	}
}

func (r *Reader) Read(bits int) (uint32, error) {
	if bits <= 0 || bits > 32 {
		return 0, fmt.Errorf("bits should be in the range (0,32]: %w", ErrBitsInvalidRange)
	}
	if r.numBitsRead+bits > r.totalBits {
		return 0, fmt.Errorf("totalBits specified = %d, bits + numBitsRead = %d : %w",
			r.totalBits, r.numBitsRead+bits, ErrBitsReadExceeded)
	}

	if r.scratchBits < bits {
		var tmp uint32
		if err := binary.Read(r.reader, binary.LittleEndian, &tmp); err != nil {
			return 0, fmt.Errorf("could not read into scratch: %w", err)
		}

		r.scratch |= uint64(tmp) << r.scratchBits
		r.scratchBits += 32
	}

	output := uint32(r.scratch) & ((uint32(1) << bits) - 1)
	r.scratch >>= bits
	r.scratchBits -= bits
	r.numBitsRead += bits

	return output, nil
}
