package encoder

var encoders = make(map[string]Encoder)

type Encoder interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
	String() string
}

func Register(name string, enc Encoder)  {
	encoders[name] = enc
}

func GetEncoder(name string) Encoder {
	if _, ok := encoders[name]; !ok {
		panic("encoders: Register called twice for encoder " + name)
	}

	return encoders[name]
}