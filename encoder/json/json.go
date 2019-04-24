package json

import (
	"encoding/json"

	"github.com/astaxie/beego/encoder"
)

func init()  {
	e := NewEncoder()
	encoder.Register(e.String(), e)
}

type jsonEncoder struct{}


func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}


func (j jsonEncoder) Encode(v interface{}) ([]byte, error) {
	return Encode(v)
}

func Decode(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

func (j jsonEncoder) Decode(d []byte, v interface{}) error {
	return Decode(d, v)
}

func (j jsonEncoder) String() string {
	return "json"
}

func NewEncoder() encoder.Encoder {
	return jsonEncoder{}
}