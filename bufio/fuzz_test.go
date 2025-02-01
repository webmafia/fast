package bufio_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/webmafia/fast/bufio"
)

func buildData(n int, letters []byte) []byte {
	data := make([]byte, n)
	l := len(letters)

	for i := range data {
		data[i] = letters[i%l]
	}

	return data
}

func FuzzFoo(f *testing.F) {
	f.Add(1)
	f.Fuzz(func(t *testing.T, n int) {})
}

func FuzzReadByte(f *testing.F) {
	f.Add(uint16(32), []byte("hello world"))

	f.Fuzz(func(t *testing.T, bufMaxSize uint16, data []byte) {
		r := bytes.NewReader(data)
		br := bufio.NewReader(r, int(bufMaxSize))

		for i := range data {
			c, err := br.ReadByte()

			if err != nil {
				if err != io.EOF || i != len(data)-1 {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if c != data[i] {
				t.Errorf("expected 0x%2x, got 0x%2x", data[i], c)
			}
		}
	})
}

func FuzzDiscard(f *testing.F) {
	f.Add(uint16(32), []byte("hello world"), uint16(6))

	f.Fuzz(func(t *testing.T, bufMaxSize uint16, data []byte, discard uint16) {
		r := bytes.NewReader(data)
		br := bufio.NewReader(r, int(bufMaxSize))

		n, err := br.Discard(int(discard))

		if err != nil {
			if err != io.EOF || n < len(data) {
				t.Error(err)
			}
		}

		if n != int(discard) && err != io.EOF {
			t.Errorf("Discard(%d): expected %d, got %d", discard, discard, n)
		}

		if n < len(data) {
			c, err := br.ReadByte()

			if c != data[n] {
				t.Errorf("ReadByte: expected 0x%2x, got 0x%2x", data[n], c)
			}

			if err != nil && err != io.EOF {
				t.Error(err)
			}
		}
	})
}

func FuzzDiscardUntil(f *testing.F) {
	f.Add(uint16(32), []byte("hello world"), byte('w'))

	f.Fuzz(func(t *testing.T, bufMaxSize uint16, data []byte, discardUntil byte) {
		r := bytes.NewReader(data)
		br := bufio.NewReader(r, int(bufMaxSize))

		n, err := br.DiscardUntil(discardUntil)

		if err != nil {
			if err != io.EOF || n < len(data) {
				t.Error(err)
			}
		}

		idx := bytes.IndexByte(data, discardUntil)

		if idx >= 0 {
			if n != idx {
				t.Errorf("Discard(0x%2x): expected %d, got %d", discardUntil, idx, n)
			}

			c, err := br.ReadByte()

			if c != data[n] {
				t.Errorf("ReadByte: expected 0x%2x ('%s'), got 0x%2x ('%s')", data[n], string(data[n]), c, string(c))
			}

			if err != nil && err != io.EOF {
				t.Error(err)
			}
		}
	})
}

func FuzzReadSlice(f *testing.F) {
	f.Add(uint16(32), []byte("hello world"), byte('w'))

	f.Fuzz(func(t *testing.T, bufMaxSize uint16, data []byte, discardUntil byte) {
		r := bytes.NewReader(data)
		br := bufio.NewReader(r, int(bufMaxSize))

		buf, err := br.ReadSlice(discardUntil)

		if err != nil {
			if err != io.EOF && err != bufio.ErrBufferFull {
				t.Error(err)
			}
		}

		if err != bufio.ErrBufferFull {
			idx := bytes.IndexByte(data, discardUntil)

			if idx >= 0 {
				if !bytes.Equal(buf, data[:idx]) {
					t.Errorf("Discard(0x%2x): expected %v, got %v", discardUntil, data[:idx], buf)
				}
			}
		}
	})
}

func FuzzReader(f *testing.F) {
	const mb = 1024 * 1024

	type testCase struct {
		dataLen         uint
		bufMaxSize      uint
		discard         uint16
		discardToLetter uint8
	}

	cases := []testCase{
		{
			dataLen:    1024,
			bufMaxSize: 1024,
			discard:    0,
		},
	}

	for _, c := range cases {
		f.Add(
			c.dataLen,
			c.bufMaxSize,
			c.discard,
		)
	}

	letters := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ012345")
	data := buildData(mb, letters)

	f.Fuzz(func(
		t *testing.T,
		dataLen uint,
		bufMaxSize uint,
		discard uint16,
	) {
		if dataLen > mb {
			t.Skipf("dataLen too high (%d)", dataLen)
		}

		if bufMaxSize > mb {
			t.Skipf("dataLen too high (%d)", bufMaxSize)
		}

		r := bytes.NewReader(data[:dataLen])
		br := bufio.NewReader(r, int(bufMaxSize))

		c, err := br.ReadByte()

		if err != nil {
			t.Error(err)
		}

		if c != 'A' {
			t.Errorf("expected 0x%2x, got 0x%2x", 'A', c)
		}

		n, err := br.Discard(int(discard))

		if err != nil {
			t.Error(err)
		}

		if n != int(discard) {
			t.Errorf("Discard(%d): expected %d, got %d", discard, discard, n)
		}

		c, err = br.ReadByte()

		if err != nil {
			t.Error(err)
		}

		if letter := letters[(n+1)%len(letters)]; c != letter {
			t.Errorf("expected 0x%2x, got 0x%2x", letter, c)
		}

	})
}
