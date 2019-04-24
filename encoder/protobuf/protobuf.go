package protobuf

import (
	"github.com/gogo/protobuf/proto"

	"github.com/astaxie/beego/encoder"
)

func init() {
	e := NewEncoder()
	encoder.Register(e.String(), e)
}

type protoEncoder struct{}

func Encode(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (p protoEncoder) Encode(v interface{}) ([]byte, error) {
	return Encode(v)
}

func Decode(d []byte, v interface{}) error {
	return proto.Unmarshal(d, v.(proto.Message))
}

func (p protoEncoder) Decode(d []byte, v interface{}) error {
	return Decode(d, v)
}

func (p protoEncoder) String() string {
	return "protobuf"
}

func NewEncoder() encoder.Encoder {
	return protoEncoder{}
}
