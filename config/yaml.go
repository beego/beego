package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/wendal/goyaml2"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type YAMLConfig struct {
}

func (yaml *YAMLConfig) Parse(filename string) (ConfigContainer, error) {
	y := &YAMLConfigContainer{
		data: make(map[string]interface{}),
	}
	cnf, err := ReadYmlReader(filename)
	if err != nil {
		return nil, err
	}
	y.data = cnf
	return y, nil
}

// 从Reader读取YAML
func ReadYmlReader(path string) (cnf map[string]interface{}, err error) {
	err = nil
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	err = nil
	buf, err := ioutil.ReadAll(f)
	if err != nil || len(buf) < 3 {
		return
	}

	if string(buf[0:1]) == "{" {
		log.Println("Look lile a Json, try it")
		err = json.Unmarshal(buf, &cnf)
		if err == nil {
			log.Println("It is Json Map")
			return
		}
	}

	_map, _err := goyaml2.Read(bytes.NewBuffer(buf))
	if _err != nil {
		log.Println("Goyaml2 ERR>", string(buf), _err)
		//err = goyaml.Unmarshal(buf, &cnf)
		err = _err
		return
	}
	if _map == nil {
		log.Println("Goyaml2 output nil? Pls report bug\n" + string(buf))
	}
	cnf, ok := _map.(map[string]interface{})
	if !ok {
		log.Println("Not a Map? >> ", string(buf), _map)
		cnf = nil
	}
	return
}

type YAMLConfigContainer struct {
	data map[string]interface{}
	sync.Mutex
}

func (c *YAMLConfigContainer) Bool(key string) (bool, error) {
	if v, ok := c.data[key].(bool); ok {
		return v, nil
	}
	return false, errors.New("not bool value")
}

func (c *YAMLConfigContainer) Int(key string) (int, error) {
	if v, ok := c.data[key].(int64); ok {
		return int(v), nil
	}
	return 0, errors.New("not int value")
}

func (c *YAMLConfigContainer) Int64(key string) (int64, error) {
	if v, ok := c.data[key].(int64); ok {
		return v, nil
	}
	return 0, errors.New("not bool value")
}

func (c *YAMLConfigContainer) Float(key string) (float64, error) {
	if v, ok := c.data[key].(float64); ok {
		return v, nil
	}
	return 0.0, errors.New("not float64 value")
}

func (c *YAMLConfigContainer) String(key string) string {
	if v, ok := c.data[key].(string); ok {
		return v
	}
	return ""
}

func (c *YAMLConfigContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	return nil
}

func (c *YAMLConfigContainer) DIY(key string) (v interface{}, err error) {
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return nil, errors.New("not exist key")
}

func init() {
	Register("yaml", &YAMLConfig{})
}
