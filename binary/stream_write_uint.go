// Borrowed from jsoniter (https://github.com/json-iterator/go)
package binary

import (
	"encoding/binary"
)

// Write uint8
func (b *StreamWriter) WriteUint8(v uint8) error {
	return b.WriteByte(v)
}

// Write uint16
func (b *StreamWriter) WriteUint16(v uint16) (err error) {
	_, err = b.Write([]byte{
		byte(v),
		byte(v >> 8),
	})
	return
}

// Write uint32
func (b *StreamWriter) WriteUint32(v uint32) (err error) {
	_, err = b.Write([]byte{
		byte(v),
		byte(v >> 8),
		byte(v >> 16),
		byte(v >> 24),
	})
	return
}

// Write uint64
func (b *StreamWriter) WriteUint64(v uint64) (err error) {
	_, err = b.Write([]byte{
		byte(v),
		byte(v >> 8),
		byte(v >> 16),
		byte(v >> 24),
		byte(v >> 32),
		byte(v >> 40),
		byte(v >> 48),
		byte(v >> 56),
	})
	return
}

// Write uint
func (b *StreamWriter) WriteUint(v uint) error {
	return b.WriteUint64(uint64(v))
}

func (b *StreamWriter) WriteUvarint(v uint64) (err error) {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], v)
	_, err = b.Write(buf[:n])
	return
}
