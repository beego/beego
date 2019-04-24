package yaml

import (
	"github.com/astaxie/beego/encoder"

	"gopkg.in/yaml.v2"
)

func init() {
	e := NewEncoder()
	encoder.Register(e.String(), e)
}

type yamlEncoder struct{}

func Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (j yamlEncoder) Encode(v interface{}) ([]byte, error) {
	return Encode(v)
}

func Decode(d []byte, v interface{}) error {
	return yaml.Unmarshal(d, v)
}

func (j yamlEncoder) Decode(d []byte, v interface{}) error {
	return Decode(d, v)
}

func (j yamlEncoder) String() string {
	return "yaml"
}

func NewEncoder() encoder.Encoder {
	return yamlEncoder{}
}
