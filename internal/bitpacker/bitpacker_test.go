package bitpacker_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/trphume/coachbuf/internal/bitpacker"
)

func TestWriter(t *testing.T) {
	t.Run("successful when writing within first 32 bit word", func(t *testing.T) {
		t.Parallel()

		w := bitpacker.NewWriter()

		// all 1s for entire 32 bit, but we only take write 30 bits
		if err := w.Write(0b11111111111111111111111111111111, 30); err != nil {
			t.Errorf("Write() = %v, want nil", err.Error())
		}

		if err := w.FlushBits(); err != nil {
			t.Errorf("FlushBits() = %v, want nil", err.Error())
		}

		want := []byte{255, 255, 255, 63}
		result := w.Bytes()
		if string(want) != string(result) {
			t.Errorf("Bytes() = %v, want %v", result, want)
		}

		wantBitsWritten := 30
		bitsWritten := w.NumBitsWritten()
		if bitsWritten != wantBitsWritten {
			t.Errorf("NumBitsWritten() = %v, want %v", bitsWritten, wantBitsWritten)
		}
	})

	t.Run("successful when writing exceed first 32 bit word", func(t *testing.T) {
		t.Parallel()

		w := bitpacker.NewWriter()

		// for the first word 30 bits taken from first write then last 2 bit taken from second write
		// second word will have 28 bits from second write after flushing
		if err := w.Write(0b11111111111111111111111111111111, 30); err != nil {
			t.Errorf("Write() = %v, want nil", err.Error())
		}
		if err := w.Write(0b11111111111111111111111111111110, 30); err != nil {
			t.Errorf("Write() = %v, want nil", err.Error())
		}

		if err := w.FlushBits(); err != nil {
			t.Errorf("FlushBits() = %v, want nil", err.Error())
		}

		want := []byte{
			255, 255, 255, 191,
			255, 255, 255, 15}
		result := w.Bytes()
		if string(want) != string(result) {
			t.Errorf("Bytes() = %v, want %v", result, want)
		}

		wantBitsWritten := 60
		bitsWritten := w.NumBitsWritten()
		if bitsWritten != wantBitsWritten {
			t.Errorf("NumBitsWritten() = %v, want %v", bitsWritten, wantBitsWritten)
		}
	})

	t.Run("data missing when not using Writer.Flush at the end", func(t *testing.T) {
		t.Parallel()

		w := bitpacker.NewWriter()

		// for the first word 30 bits taken from first write then last 2 bit taken from second write
		if err := w.Write(0b11111111111111111111111111111111, 30); err != nil {
			t.Errorf("Write() = %v, want nil", err.Error())
		}
		if err := w.Write(0b11111111111111111111111111111110, 30); err != nil {
			t.Errorf("Write() = %v, want nil", err.Error())
		}

		want := []byte{255, 255, 255, 191}
		result := w.Bytes()
		if string(want) != string(result) {
			t.Errorf("Bytes() = %v, want %v", result, want)
		}

		wantBitsWritten := 60
		bitsWritten := w.NumBitsWritten()
		if bitsWritten != wantBitsWritten {
			t.Errorf("NumBitsWritten() = %v, want %v", bitsWritten, wantBitsWritten)
		}
	})

	t.Run("error when bits args exceeds 32", func(t *testing.T) {
		t.Parallel()

		w := bitpacker.NewWriter()
		if err := w.Write(0, 33); !errors.Is(err, bitpacker.ErrBitsInvalidRange) {
			t.Errorf("Write() = %v, want %v", err.Error(), bitpacker.ErrBitsInvalidRange.Error())
		}
	})

	t.Run("error when bits args is 0 or less", func(t *testing.T) {
		t.Parallel()

		w := bitpacker.NewWriter()
		if err := w.Write(0, 33); !errors.Is(err, bitpacker.ErrBitsInvalidRange) {
			t.Errorf("Write() = %v, want %v", err.Error(), bitpacker.ErrBitsInvalidRange.Error())
		}
	})

	t.Run("error when writer flushed more than once", func(t *testing.T) {
		t.Parallel()

		w := bitpacker.NewWriter()
		if err := w.Write(0, 32); err != nil {
			t.Errorf("Write() = %v, want %v", err.Error(), nil)
		}
		if err := w.FlushBits(); err != nil {
			t.Errorf("Write() = %v, want %v", err.Error(), nil)
		}
		if err := w.FlushBits(); err == nil {
			t.Errorf("Write() = %v, want %v", err.Error(), bitpacker.ErrMethodCallNotAllowed)
		}
	})

	t.Run("error when writing when writer was already flushed", func(t *testing.T) {
		t.Parallel()

		w := bitpacker.NewWriter()
		if err := w.Write(0, 32); err != nil {
			t.Errorf("Write() = %v, want %v", err.Error(), nil)
		}
		if err := w.FlushBits(); err != nil {
			t.Errorf("Write() = %v, want %v", err.Error(), nil)
		}
		if err := w.Write(0, 32); err == nil {
			t.Errorf("Write() = %v, want %v", err.Error(), bitpacker.ErrMethodCallNotAllowed)
		}
	})

}

func TestReader(t *testing.T) {
	t.Run("successful when reading from within first 32 bit", func(t *testing.T) {
		t.Parallel()

		bRdr := bytes.NewReader([]byte{255, 100, 255, 1})
		rdr := bitpacker.NewReader(bRdr, 4)

		// the first 16 bits are 0110010011111111 which is 25855 in decimal
		var want uint32 = 25855
		result, err := rdr.Read(16)
		if err != nil {
			t.Errorf("Read() = %v, want %v", err.Error(), 25855)
		}
		if want != result {
			t.Errorf("Read() = %v, want %v", result, want)
		}
	})

	t.Run("successful when reading from across multiple word", func(t *testing.T) {
		t.Parallel()

		bRdr := bytes.NewReader([]byte{
			255, 100, 255, 1, // 00000001111111110110010011111111
			100, 234, 90, 0, // 00000000010110101110101001100100
		})
		rdr := bitpacker.NewReader(bRdr, 8)

		var want uint32 = 0b01111111110110010011111111
		result, err := rdr.Read(26)
		if err != nil {
			t.Errorf("Read() = %v, want %v", err.Error(), want)
		}
		if want != result {
			t.Errorf("Read() = %v, want %v", result, want)
		}

		want = 0b1100100000000
		result, err = rdr.Read(13)
		if err != nil {
			t.Errorf("Read() = %v, want %v", err.Error(), want)
		}
		if want != result {
			t.Errorf("Read() = %v, want %v", result, want)
		}
	})

	t.Run("error when bits to read exceed 32", func(t *testing.T) {
		t.Parallel()

		bRdr := bytes.NewReader([]byte{255, 255, 255, 255})
		rdr := bitpacker.NewReader(bRdr, 8)

		_, err := rdr.Read(40)
		if !errors.Is(err, bitpacker.ErrBitsInvalidRange) {
			t.Errorf("Read() = %v, want %v", err.Error(), bitpacker.ErrBitsInvalidRange)
		}
	})

	t.Run("error when bits to read is less than 1", func(t *testing.T) {
		t.Parallel()

		bRdr := bytes.NewReader([]byte{255, 255, 255, 255})
		rdr := bitpacker.NewReader(bRdr, 4)

		_, err := rdr.Read(0)
		if !errors.Is(err, bitpacker.ErrBitsInvalidRange) {
			t.Errorf("Read() = %v, want %v", err.Error(), bitpacker.ErrBitsInvalidRange)
		}
	})

	t.Run("error when bits to read exceed number of total bits specified", func(t *testing.T) {
		t.Parallel()

		bRdr := bytes.NewReader([]byte{255, 255, 255, 255})
		rdr := bitpacker.NewReader(bRdr, 4)

		_, err := rdr.Read(32)
		if err != nil {
			t.Errorf("Read() = %v, want %v", err.Error(), nil)
		}
		_, err = rdr.Read(10)
		if !errors.Is(err, bitpacker.ErrBitsReadExceeded) {
			t.Errorf("Read() = %v, want %v", err.Error(), bitpacker.ErrBitsReadExceeded)
		}
	})
}
