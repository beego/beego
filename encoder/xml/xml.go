package xml

import (
	"encoding/xml"

	"github.com/astaxie/beego/encoder"
)

func init() {
	e := NewEncoder()
	encoder.Register(e.String(), e)
}

type xmlEncoder struct{}

func Encode(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (j xmlEncoder) Encode(v interface{}) ([]byte, error) {
	return Encode(v)
}

func Decode(d []byte, v interface{}) error {
	return xml.Unmarshal(d, v)
}

func (j xmlEncoder) Decode(d []byte, v interface{}) error {
	return Decode(d, v)
}

func (j xmlEncoder) String() string {
	return "xml"
}

func NewEncoder() encoder.Encoder {
	return xmlEncoder{}
}
