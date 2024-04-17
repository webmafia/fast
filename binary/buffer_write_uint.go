// Borrowed from jsoniter (https://github.com/json-iterator/go)
package binary

import (
	"encoding/binary"
)

// Write uint8
func (b *BufferWriter) WriteUint8(v uint8) {
	b.WriteByte(v)
}

// Write uint16
func (b *BufferWriter) WriteUint16(v uint16) {
	b.buf = binary.LittleEndian.AppendUint16(b.buf, v)
}

// Write uint32
func (b *BufferWriter) WriteUint32(v uint32) {
	b.buf = binary.LittleEndian.AppendUint32(b.buf, v)
}

// Write uint64
func (b *BufferWriter) WriteUint64(v uint64) {
	b.buf = binary.LittleEndian.AppendUint64(b.buf, v)
}

// Write uint
func (b *BufferWriter) WriteUint(v uint) {
	b.WriteUint64(uint64(v))
}

func (b *BufferWriter) WriteUvarint(v uint64) {
	b.buf = binary.AppendUvarint(b.buf, v)
}
