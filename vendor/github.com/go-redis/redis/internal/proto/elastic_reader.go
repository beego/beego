package proto

import (
	"bytes"
	"errors"
	"io"
)

const defaultBufSize = 4096

// ElasticBufReader is like bufio.Reader but instead of returning ErrBufferFull
// it automatically grows the buffer.
type ElasticBufReader struct {
	buf  []byte
	rd   io.Reader // reader provided by the client
	r, w int       // buf read and write positions
	err  error
}

func NewElasticBufReader(rd io.Reader) *ElasticBufReader {
	return &ElasticBufReader{
		buf: make([]byte, defaultBufSize),
		rd:  rd,
	}
}

func (b *ElasticBufReader) Reset(rd io.Reader) {
	b.rd = rd
	b.r, b.w = 0, 0
	b.err = nil
}

func (b *ElasticBufReader) Buffer() []byte {
	return b.buf
}

func (b *ElasticBufReader) ResetBuffer(buf []byte) {
	b.buf = buf
	b.r, b.w = 0, 0
	b.err = nil
}

func (b *ElasticBufReader) reset(buf []byte, rd io.Reader) {
	*b = ElasticBufReader{
		buf: buf,
		rd:  rd,
	}
}

// Buffered returns the number of bytes that can be read from the current buffer.
func (b *ElasticBufReader) Buffered() int {
	return b.w - b.r
}

func (b *ElasticBufReader) Bytes() []byte {
	return b.buf[b.r:b.w]
}

var errNegativeRead = errors.New("bufio: reader returned negative count from Read")

// fill reads a new chunk into the buffer.
func (b *ElasticBufReader) fill() {
	// Slide existing data to beginning.
	if b.r > 0 {
		copy(b.buf, b.buf[b.r:b.w])
		b.w -= b.r
		b.r = 0
	}

	if b.w >= len(b.buf) {
		panic("bufio: tried to fill full buffer")
	}

	// Read new data: try a limited number of times.
	const maxConsecutiveEmptyReads = 100
	for i := maxConsecutiveEmptyReads; i > 0; i-- {
		n, err := b.rd.Read(b.buf[b.w:])
		if n < 0 {
			panic(errNegativeRead)
		}
		b.w += n
		if err != nil {
			b.err = err
			return
		}
		if n > 0 {
			return
		}
	}
	b.err = io.ErrNoProgress
}

func (b *ElasticBufReader) readErr() error {
	err := b.err
	b.err = nil
	return err
}

func (b *ElasticBufReader) Read(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		return 0, b.readErr()
	}
	if b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		if len(p) >= len(b.buf) {
			// Large read, empty buffer.
			// Read directly into p to avoid copy.
			n, b.err = b.rd.Read(p)
			if n < 0 {
				panic(errNegativeRead)
			}
			return n, b.readErr()
		}
		// One read.
		// Do not use b.fill, which will loop.
		b.r = 0
		b.w = 0
		n, b.err = b.rd.Read(b.buf)
		if n < 0 {
			panic(errNegativeRead)
		}
		if n == 0 {
			return 0, b.readErr()
		}
		b.w += n
	}

	// copy as much as we can
	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	return n, nil
}

func (b *ElasticBufReader) ReadSlice(delim byte) (line []byte, err error) {
	for {
		// Search buffer.
		if i := bytes.IndexByte(b.buf[b.r:b.w], delim); i >= 0 {
			line = b.buf[b.r : b.r+i+1]
			b.r += i + 1
			break
		}

		// Pending error?
		if b.err != nil {
			line = b.buf[b.r:b.w]
			b.r = b.w
			err = b.readErr()
			break
		}

		// Buffer full?
		if b.Buffered() >= len(b.buf) {
			b.grow(len(b.buf) + defaultBufSize)
		}

		b.fill() // buffer is not full
	}

	return
}

func (b *ElasticBufReader) ReadLine() (line []byte, err error) {
	line, err = b.ReadSlice('\n')
	if len(line) == 0 {
		if err != nil {
			line = nil
		}
		return
	}
	err = nil

	if line[len(line)-1] == '\n' {
		drop := 1
		if len(line) > 1 && line[len(line)-2] == '\r' {
			drop = 2
		}
		line = line[:len(line)-drop]
	}
	return
}

func (b *ElasticBufReader) ReadByte() (byte, error) {
	for b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		b.fill() // buffer is empty
	}
	c := b.buf[b.r]
	b.r++
	return c, nil
}

func (b *ElasticBufReader) ReadN(n int) ([]byte, error) {
	b.grow(n)
	for b.Buffered() < n {
		// Pending error?
		if b.err != nil {
			buf := b.buf[b.r:b.w]
			b.r = b.w
			return buf, b.readErr()
		}

		b.fill()
	}

	buf := b.buf[b.r : b.r+n]
	b.r += n
	return buf, nil
}

func (b *ElasticBufReader) grow(n int) {
	if b.w-b.r >= n {
		return
	}

	// Slide existing data to beginning.
	if b.r > 0 {
		copy(b.buf, b.buf[b.r:b.w])
		b.w -= b.r
		b.r = 0
	}

	// Extend buffer if needed.
	if d := n - len(b.buf); d > 0 {
		b.buf = append(b.buf, make([]byte, d)...)
	}
}
