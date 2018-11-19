package proto

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/go-redis/redis/internal/util"
)

const (
	ErrorReply  = '-'
	StatusReply = '+'
	IntReply    = ':'
	StringReply = '$'
	ArrayReply  = '*'
)

//------------------------------------------------------------------------------

const Nil = RedisError("redis: nil")

type RedisError string

func (e RedisError) Error() string { return string(e) }

//------------------------------------------------------------------------------

type MultiBulkParse func(*Reader, int64) (interface{}, error)

type Reader struct {
	rd   *bufio.Reader
	_buf []byte
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{
		rd:   bufio.NewReader(rd),
		_buf: make([]byte, 64),
	}
}

func (r *Reader) Reset(rd io.Reader) {
	r.rd.Reset(rd)
}

func (r *Reader) ReadLine() ([]byte, error) {
	line, isPrefix, err := r.rd.ReadLine()
	if err != nil {
		return nil, err
	}
	if isPrefix {
		return nil, bufio.ErrBufferFull
	}
	if len(line) == 0 {
		return nil, fmt.Errorf("redis: reply is empty")
	}
	if isNilReply(line) {
		return nil, Nil
	}
	return line, nil
}

func (r *Reader) ReadReply(m MultiBulkParse) (interface{}, error) {
	line, err := r.ReadLine()
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case ErrorReply:
		return nil, ParseErrorReply(line)
	case StatusReply:
		return string(line[1:]), nil
	case IntReply:
		return util.ParseInt(line[1:], 10, 64)
	case StringReply:
		return r.readStringReply(line)
	case ArrayReply:
		n, err := parseArrayLen(line)
		if err != nil {
			return nil, err
		}
		return m(r, n)
	}
	return nil, fmt.Errorf("redis: can't parse %.100q", line)
}

func (r *Reader) ReadIntReply() (int64, error) {
	line, err := r.ReadLine()
	if err != nil {
		return 0, err
	}
	switch line[0] {
	case ErrorReply:
		return 0, ParseErrorReply(line)
	case IntReply:
		return util.ParseInt(line[1:], 10, 64)
	default:
		return 0, fmt.Errorf("redis: can't parse int reply: %.100q", line)
	}
}

func (r *Reader) ReadString() (string, error) {
	line, err := r.ReadLine()
	if err != nil {
		return "", err
	}
	switch line[0] {
	case ErrorReply:
		return "", ParseErrorReply(line)
	case StringReply:
		return r.readStringReply(line)
	case StatusReply:
		return string(line[1:]), nil
	case IntReply:
		return string(line[1:]), nil
	default:
		return "", fmt.Errorf("redis: can't parse reply=%.100q reading string", line)
	}
}

func (r *Reader) readStringReply(line []byte) (string, error) {
	if isNilReply(line) {
		return "", Nil
	}

	replyLen, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return "", err
	}

	b := make([]byte, replyLen+2)
	_, err = io.ReadFull(r.rd, b)
	if err != nil {
		return "", err
	}

	return util.BytesToString(b[:replyLen]), nil
}

func (r *Reader) ReadArrayReply(m MultiBulkParse) (interface{}, error) {
	line, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	switch line[0] {
	case ErrorReply:
		return nil, ParseErrorReply(line)
	case ArrayReply:
		n, err := parseArrayLen(line)
		if err != nil {
			return nil, err
		}
		return m(r, n)
	default:
		return nil, fmt.Errorf("redis: can't parse array reply: %.100q", line)
	}
}

func (r *Reader) ReadArrayLen() (int64, error) {
	line, err := r.ReadLine()
	if err != nil {
		return 0, err
	}
	switch line[0] {
	case ErrorReply:
		return 0, ParseErrorReply(line)
	case ArrayReply:
		return parseArrayLen(line)
	default:
		return 0, fmt.Errorf("redis: can't parse array reply: %.100q", line)
	}
}

func (r *Reader) ReadScanReply() ([]string, uint64, error) {
	n, err := r.ReadArrayLen()
	if err != nil {
		return nil, 0, err
	}
	if n != 2 {
		return nil, 0, fmt.Errorf("redis: got %d elements in scan reply, expected 2", n)
	}

	cursor, err := r.ReadUint()
	if err != nil {
		return nil, 0, err
	}

	n, err = r.ReadArrayLen()
	if err != nil {
		return nil, 0, err
	}

	keys := make([]string, n)
	for i := int64(0); i < n; i++ {
		key, err := r.ReadString()
		if err != nil {
			return nil, 0, err
		}
		keys[i] = key
	}

	return keys, cursor, err
}

func (r *Reader) ReadInt() (int64, error) {
	b, err := r.readTmpBytesReply()
	if err != nil {
		return 0, err
	}
	return util.ParseInt(b, 10, 64)
}

func (r *Reader) ReadUint() (uint64, error) {
	b, err := r.readTmpBytesReply()
	if err != nil {
		return 0, err
	}
	return util.ParseUint(b, 10, 64)
}

func (r *Reader) ReadFloatReply() (float64, error) {
	b, err := r.readTmpBytesReply()
	if err != nil {
		return 0, err
	}
	return util.ParseFloat(b, 64)
}

func (r *Reader) readTmpBytesReply() ([]byte, error) {
	line, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	switch line[0] {
	case ErrorReply:
		return nil, ParseErrorReply(line)
	case StringReply:
		return r._readTmpBytesReply(line)
	case StatusReply:
		return line[1:], nil
	default:
		return nil, fmt.Errorf("redis: can't parse string reply: %.100q", line)
	}
}

func (r *Reader) _readTmpBytesReply(line []byte) ([]byte, error) {
	if isNilReply(line) {
		return nil, Nil
	}

	replyLen, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, err
	}

	buf := r.buf(replyLen + 2)
	_, err = io.ReadFull(r.rd, buf)
	if err != nil {
		return nil, err
	}

	return buf[:replyLen], nil
}

func (r *Reader) buf(n int) []byte {
	if d := n - cap(r._buf); d > 0 {
		r._buf = append(r._buf, make([]byte, d)...)
	}
	return r._buf[:n]
}

func isNilReply(b []byte) bool {
	return len(b) == 3 &&
		(b[0] == StringReply || b[0] == ArrayReply) &&
		b[1] == '-' && b[2] == '1'
}

func ParseErrorReply(line []byte) error {
	return RedisError(string(line[1:]))
}

func parseArrayLen(line []byte) (int64, error) {
	if isNilReply(line) {
		return 0, Nil
	}
	return util.ParseInt(line[1:], 10, 64)
}
