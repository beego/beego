package hcl

import (
	"encoding/json"

	"github.com/astaxie/beego/encoder"
	"github.com/hashicorp/hcl"
)

func init()  {
	e := NewEncoder()
	encoder.Register(e.String(), e)
}

type hclEncoder struct{}

func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (h hclEncoder) Encode(v interface{}) ([]byte, error) {
	return Encode(v)
}

func Decode(d []byte, v interface{}) error {
	return hcl.Unmarshal(d, v)
}

func (h hclEncoder) Decode(d []byte, v interface{}) error {
	return Decode(d, v)
}

func (h hclEncoder) String() string {
	return "hcl"
}

func NewEncoder() encoder.Encoder {
	return hclEncoder{}
}