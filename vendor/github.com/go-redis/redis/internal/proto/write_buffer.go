package proto

import (
	"encoding"
	"fmt"
	"strconv"
)

type WriteBuffer struct {
	b []byte
}

func NewWriteBuffer() *WriteBuffer {
	return &WriteBuffer{
		b: make([]byte, 0, 4096),
	}
}

func (w *WriteBuffer) Len() int      { return len(w.b) }
func (w *WriteBuffer) Bytes() []byte { return w.b }
func (w *WriteBuffer) Reset()        { w.b = w.b[:0] }

func (w *WriteBuffer) Append(args []interface{}) error {
	w.b = append(w.b, ArrayReply)
	w.b = strconv.AppendUint(w.b, uint64(len(args)), 10)
	w.b = append(w.b, '\r', '\n')

	for _, arg := range args {
		if err := w.append(arg); err != nil {
			return err
		}
	}
	return nil
}

func (w *WriteBuffer) append(val interface{}) error {
	switch v := val.(type) {
	case nil:
		w.AppendString("")
	case string:
		w.AppendString(v)
	case []byte:
		w.AppendBytes(v)
	case int:
		w.AppendString(formatInt(int64(v)))
	case int8:
		w.AppendString(formatInt(int64(v)))
	case int16:
		w.AppendString(formatInt(int64(v)))
	case int32:
		w.AppendString(formatInt(int64(v)))
	case int64:
		w.AppendString(formatInt(v))
	case uint:
		w.AppendString(formatUint(uint64(v)))
	case uint8:
		w.AppendString(formatUint(uint64(v)))
	case uint16:
		w.AppendString(formatUint(uint64(v)))
	case uint32:
		w.AppendString(formatUint(uint64(v)))
	case uint64:
		w.AppendString(formatUint(v))
	case float32:
		w.AppendString(formatFloat(float64(v)))
	case float64:
		w.AppendString(formatFloat(v))
	case bool:
		if v {
			w.AppendString("1")
		} else {
			w.AppendString("0")
		}
	case encoding.BinaryMarshaler:
		b, err := v.MarshalBinary()
		if err != nil {
			return err
		}
		w.AppendBytes(b)
	default:
		return fmt.Errorf(
			"redis: can't marshal %T (consider implementing encoding.BinaryMarshaler)", val)
	}
	return nil
}

func (w *WriteBuffer) AppendString(s string) {
	w.b = append(w.b, StringReply)
	w.b = strconv.AppendUint(w.b, uint64(len(s)), 10)
	w.b = append(w.b, '\r', '\n')
	w.b = append(w.b, s...)
	w.b = append(w.b, '\r', '\n')
}

func (w *WriteBuffer) AppendBytes(p []byte) {
	w.b = append(w.b, StringReply)
	w.b = strconv.AppendUint(w.b, uint64(len(p)), 10)
	w.b = append(w.b, '\r', '\n')
	w.b = append(w.b, p...)
	w.b = append(w.b, '\r', '\n')
}

func formatInt(n int64) string {
	return strconv.FormatInt(n, 10)
}

func formatUint(u uint64) string {
	return strconv.FormatUint(u, 10)
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
