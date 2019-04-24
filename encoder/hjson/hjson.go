package hhjson

import (
	"github.com/astaxie/beego/encoder"
	"github.com/hjson/hjson-go"
)

func init()  {
	e := NewEncoder()
	encoder.Register(e.String(), e)
}

type hjsonEncoder struct{}

func Encode(v interface{}) ([]byte, error) {
	return hjson.Marshal(v)
}

func (j hjsonEncoder) Encode(v interface{}) ([]byte, error) {
	return Encode(v)
}

func Decode(d []byte, v interface{}) error {
	return hjson.Unmarshal(d, v)
}

func (j hjsonEncoder) Decode(d []byte, v interface{}) error {
	return Decode(d, v)
}

func (j hjsonEncoder) String() string {
	return "hjson"
}

func NewEncoder() encoder.Encoder {
	return hjsonEncoder{}
}
