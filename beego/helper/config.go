package helper

import (
	"json"
	"io/ioutil"
	"log"
)

var config map[string]string

func ReadConfig(filename string) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Exitf("Impossible to read %s", filename, err)
	}
	data, err := json.Decode(string(contents))
	if err != nil {
		log.Exitf("Can't parse %s as JSON", filename, err)
	}
	config = map[string]string{ }
	for key, value := range data.(map[string]interface{ }) {
		config[key] = value.(string)
	}
}

func GetConfig(key string) string {
	return config[key];
}