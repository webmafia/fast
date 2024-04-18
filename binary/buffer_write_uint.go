// Borrowed from jsoniter (https://github.com/json-iterator/go)
package binary

import (
	"encoding/binary"
)

// Write uint8
func (b *BufferWriter) WriteUint8(v uint8) error {
	return b.WriteByte(v)
}

// Write uint16
func (b *BufferWriter) WriteUint16(v uint16) error {
	b.buf = binary.LittleEndian.AppendUint16(b.buf, v)
	return nil
}

// Write uint32
func (b *BufferWriter) WriteUint32(v uint32) error {
	b.buf = binary.LittleEndian.AppendUint32(b.buf, v)
	return nil
}

// Write uint64
func (b *BufferWriter) WriteUint64(v uint64) error {
	b.buf = binary.LittleEndian.AppendUint64(b.buf, v)
	return nil
}

// Write uint
func (b *BufferWriter) WriteUint(v uint) error {
	return b.WriteUint64(uint64(v))
}

func (b *BufferWriter) WriteUvarint(v uint64) error {
	b.buf = binary.AppendUvarint(b.buf, v)
	return nil
}
