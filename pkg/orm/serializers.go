package orm

import (
	"encoding/json"
)

// SerializerMap maps from the serializer's string representation to the Serializer
var serializerMap = make(map[string]FieldSerializer)

// FieldSerializer is the interface which all serializers must implement
type FieldSerializer interface {
	Serialize(obj interface{}) string
	Deserialize(fieldContent string,originalField interface{}) interface{}
}

// RegisterSerializer is used to register a serializer into the SerializerMap
func RegisterSerializer(key string, serializer FieldSerializer) {
	serializerMap[key] = serializer
}

// jsonFieldSerializer is an inbuit serializer that implements the fieldSerializer interface
type jsonFieldSerializer struct{}

// Serialize - serializes objects to a storable string format
func (serializer *jsonFieldSerializer) Serialize(xyz interface{}) string {
	v, _ := json.Marshal(xyz)
	return string(v)
}

// Deserialize -- deserializes objects to their original json format.
func (serializer *jsonFieldSerializer) Deserialize(content string,originalField interface{}) interface{} {
	json.Unmarshal([]byte(content), originalField)
	return originalField
}

func init() {
	// registering all the inbuilt serializers
	RegisterSerializer("json", new(jsonFieldSerializer))
}
