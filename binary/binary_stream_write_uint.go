// Borrowed from jsoniter (https://github.com/json-iterator/go)
package fast

import (
	"encoding/binary"
)

// Write uint8
func (b BinaryStreamWriter) WriteUint8(v uint8) {
	b.WriteByte(v)
}

// Write uint16
func (b BinaryStreamWriter) WriteUint16(v uint16) {
	buf := b.buf.AvailableBuffer()
	buf = binary.LittleEndian.AppendUint16(buf, v)
	b.buf.Write(buf)
}

// Write uint32
func (b BinaryStreamWriter) WriteUint32(v uint32) {
	buf := b.buf.AvailableBuffer()
	buf = binary.LittleEndian.AppendUint32(buf, v)
	b.buf.Write(buf)
}

// Write uint64
func (b BinaryStreamWriter) WriteUint64(v uint64) {
	buf := b.buf.AvailableBuffer()
	buf = binary.LittleEndian.AppendUint64(buf, v)
	b.buf.Write(buf)
}

// Write uint
func (b BinaryStreamWriter) WriteUint(v uint) {
	b.WriteUint64(uint64(v))
}

func (b BinaryStreamWriter) WriteUvarint(v uint64) {
	buf := b.buf.AvailableBuffer()
	buf = binary.AppendUvarint(buf, v)
	b.buf.Write(buf)
}
