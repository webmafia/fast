package bufio

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"bufio"
)

func BenchmarkReader(b *testing.B) {
	dataSizes := [...]int{64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 252144, 524288, 1048576}
	data := make([]byte, dataSizes[len(dataSizes)-1])
	data[4095] = 1
	r := bytes.NewReader(data)
	// r := newChunkedReader(data, 512)
	br := NewReader(r, 4096)

	b.Run("Reset", func(b *testing.B) {
		for range b.N {
			r.Reset(data)
			br.Reset(r)
		}
	})

	b.Run("ReadSlice", func(b *testing.B) {
		for range b.N {
			r.Reset(data)
			br.Reset(r)
			_, err := br.ReadSlice(1)

			if err != nil {
				if err == ErrBufferFull {
					b.Skip(err)
				} else if err != io.EOF {
					b.Error(err)
				}
			}
		}
	})

	b.Run("DiscardUntil", func(b *testing.B) {
		for range b.N {
			r.Reset(data)
			br.Reset(r)
			_, err := br.DiscardUntil(1)

			if err != nil {
				if err == ErrBufferFull {
					b.Skip(err)
				} else if err != io.EOF {
					b.Error(err)
				}
			}
		}
	})

	b.Run("ReadByte", func(b *testing.B) {
		for range b.N {
			r.Reset(data)
			br.Reset(r)
			_, err := br.ReadByte()

			if err != nil {
				if err == ErrBufferFull {
					b.Skip(err)
				} else if err != io.EOF {
					b.Error(err)
				}
			}
		}
	})

	for _, size := range dataSizes {
		b.Run(fmt.Sprintf("Peek_%d", size), func(b *testing.B) {
			for range b.N {
				r.Reset(data)
				br.Reset(r)
				_, err := br.Peek(size)

				if err != nil {
					if err == ErrBufferFull {
						b.Skip(err)
					} else if err != io.EOF {
						b.Error(err)
					}
				}
			}
		})

		b.Run(fmt.Sprintf("ReadBytes_%d", size), func(b *testing.B) {
			for range b.N {
				r.Reset(data)
				br.Reset(r)
				_, err := br.ReadBytes(size)

				if err != nil {
					if err == ErrBufferFull {
						b.Skip(err)
					} else if err != io.EOF {
						b.Error(err)
					}
				}
			}
		})

		b.Run(fmt.Sprintf("Discard_%d", size), func(b *testing.B) {
			for range b.N {
				r.Reset(data)
				br.Reset(r)
				_, err := br.Discard(size)

				if err != nil {
					if err == ErrBufferFull {
						b.Skip(err)
					} else if err != io.EOF {
						b.Error(err)
					}
				}
			}
		})

		b.Run(fmt.Sprintf("Read_%d", size), func(b *testing.B) {
			dst := make([]byte, size)
			b.ResetTimer()

			for range b.N {
				r.Reset(data)
				br.Reset(r)
				_, err := br.Read(dst)

				if err != nil {
					if err == ErrBufferFull {
						b.Skip(err)
					} else if err != io.EOF {
						b.Error(err)
					}
				}
			}
		})
	}
}

func BenchmarkStandardBufio(b *testing.B) {
	dataSizes := [...]int{64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 252144, 524288, 1048576}
	data := make([]byte, dataSizes[len(dataSizes)-1])
	data[4095] = 1
	r := bytes.NewReader(data)
	// r := newChunkedReader(data, 512)
	br := bufio.NewReader(r)

	b.Run("Reset", func(b *testing.B) {
		for range b.N {
			r.Reset(data)
			br.Reset(r)
		}
	})

	b.Run("ReadSlice", func(b *testing.B) {
		for range b.N {
			r.Reset(data)
			br.Reset(r)
			_, err := br.ReadSlice(1)

			if err != nil {
				if err == ErrBufferFull {
					b.Skip(err)
				} else if err != io.EOF {
					b.Error(err)
				}
			}
		}
	})

	b.Run("ReadByte", func(b *testing.B) {
		for range b.N {
			r.Reset(data)
			br.Reset(r)
			_, err := br.ReadByte()

			if err != nil {
				if err == ErrBufferFull {
					b.Skip(err)
				} else if err != io.EOF {
					b.Error(err)
				}
			}
		}
	})

	for _, size := range dataSizes {
		b.Run(fmt.Sprintf("Peek_%d", size), func(b *testing.B) {
			for range b.N {
				r.Reset(data)
				br.Reset(r)
				_, err := br.Peek(size)

				if err != nil {
					if err == ErrBufferFull {
						b.Skip(err)
					} else if err != io.EOF {
						b.Error(err)
					}
				}
			}
		})

		b.Run(fmt.Sprintf("Discard_%d", size), func(b *testing.B) {
			for range b.N {
				r.Reset(data)
				br.Reset(r)
				_, err := br.Discard(size)

				if err != nil {
					if err == ErrBufferFull {
						b.Skip(err)
					} else if err != io.EOF {
						b.Error(err)
					}
				}
			}
		})

		b.Run(fmt.Sprintf("Read_%d", size), func(b *testing.B) {
			dst := make([]byte, size)
			b.ResetTimer()

			for range b.N {
				r.Reset(data)
				br.Reset(r)
				_, err := br.Read(dst)

				if err != nil {
					if err == ErrBufferFull {
						b.Skip(err)
					} else if err != io.EOF {
						b.Error(err)
					}
				}
			}
		})
	}
}
