package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var (
	bComment = []byte{'#'}
	bEmpty   = []byte{}
	bEqual   = []byte{'='}
	bDQuote  = []byte{'"'}
)

type IniConfig struct {
}

// ParseFile creates a new Config and parses the file configuration from the
// named file.
func (ini *IniConfig) Parse(name string) (ConfigContainer, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	cfg := &IniConfigContainer{
		file.Name(),
		make(map[int][]string),
		make(map[string]string),
		make(map[string]int64),
		sync.RWMutex{},
	}
	cfg.Lock()
	defer cfg.Unlock()
	defer file.Close()

	var comment bytes.Buffer
	buf := bufio.NewReader(file)

	for nComment, off := 0, int64(1); ; {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if bytes.Equal(line, bEmpty) {
			continue
		}

		off += int64(len(line))

		if bytes.HasPrefix(line, bComment) {
			line = bytes.TrimLeft(line, "#")
			line = bytes.TrimLeftFunc(line, unicode.IsSpace)
			comment.Write(line)
			comment.WriteByte('\n')
			continue
		}
		if comment.Len() != 0 {
			cfg.comment[nComment] = []string{comment.String()}
			comment.Reset()
			nComment++
		}

		val := bytes.SplitN(line, bEqual, 2)
		if bytes.HasPrefix([]byte(strings.TrimSpace(string(val[1]))), bDQuote) {
			val[1] = bytes.Trim([]byte(strings.TrimSpace(string(val[1]))), `"`)
		}

		key := strings.TrimSpace(string(val[0]))
		cfg.comment[nComment-1] = append(cfg.comment[nComment-1], key)
		cfg.data[key] = strings.TrimSpace(string(val[1]))
		cfg.offset[key] = off
	}
	return cfg, nil
}

// A Config represents the configuration.
type IniConfigContainer struct {
	filename string
	comment  map[int][]string  // id: []{comment, key...}; id 1 is for main comment.
	data     map[string]string // key: value
	offset   map[string]int64  // key: offset; for editing.
	sync.RWMutex
}

// Bool returns the boolean value for a given key.
func (c *IniConfigContainer) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.data[key])
}

// Int returns the integer value for a given key.
func (c *IniConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.data[key])
}

func (c *IniConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.data[key], 10, 64)
}

// Float returns the float value for a given key.
func (c *IniConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.data[key], 64)
}

// String returns the string value for a given key.
func (c *IniConfigContainer) String(key string) string {
	return c.data[key]
}

// WriteValue writes a new value for key.
func (c *IniConfigContainer) Set(key, value string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = value
	return nil
}

func (c *IniConfigContainer) DIY(key string) (v interface{}, err error) {
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return v, errors.New("key not find")
}

func init() {
	Register("ini", &IniConfig{})
}
