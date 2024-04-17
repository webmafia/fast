package binary

import (
	"bytes"
	"testing"
)

func FuzzBuffer(f *testing.F) {
	f.Add("foobar", int64(123), 456.789)
	f.Add("räksmörgås", int64(-123), -456.789)

	f.Fuzz(func(t *testing.T, s string, i int64, f float64) {
		w := NewBufferWriter(64)
		w.WriteString(s)
		w.WriteVarint(i)
		w.WriteFloat64(f)

		r := NewBufferReader(w.Bytes())

		if res := r.ReadString(len(s)); res != s {
			t.Errorf("expected '%s', got '%s'", s, res)
		}

		if res := r.ReadVarint(); res != i {
			t.Errorf("expected '%d', got '%d'", i, res)
		}

		if res := r.ReadFloat64(); res != f {
			t.Errorf("expected '%f', got '%f'", f, res)
		}
	})
}

func FuzzStreamReader(f *testing.F) {
	f.Add("foobar", int64(123), 456.789)
	f.Add("räksmörgås", int64(-123), -456.789)

	f.Fuzz(func(t *testing.T, s string, i int64, f float64) {
		w := NewBufferWriter(64)
		w.WriteString(s)
		w.WriteVarint(i)
		w.WriteFloat64(f)

		b := bytes.NewBuffer(w.Bytes())
		r := NewStreamReader(b)

		if res := r.ReadString(len(s)); res != s {
			t.Errorf("expected '%s', got '%s'", s, res)
		}

		if res := r.ReadVarint(); res != i {
			t.Errorf("expected '%d', got '%d'", i, res)
		}

		if res := r.ReadFloat64(); res != f {
			t.Errorf("expected '%f', got '%f'", f, res)
		}
	})
}
