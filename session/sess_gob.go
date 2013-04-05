package session

import (
	"bytes"
	"encoding/gob"
)

func init() {
	gob.Register([]interface{}{})
	gob.Register(map[int]interface{}{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[interface{}]interface{}{})
	gob.Register(map[string]string{})
	gob.Register(map[int]string{})
	gob.Register(map[int]int{})
	gob.Register(map[int]int64{})
}

func encodeGob(obj map[interface{}]interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(obj)
	if err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}

func decodeGob(encoded []byte) (map[interface{}]interface{}, error) {
	buf := bytes.NewBuffer(encoded)
	dec := gob.NewDecoder(buf)
	var out map[interface{}]interface{}
	err := dec.Decode(&out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
