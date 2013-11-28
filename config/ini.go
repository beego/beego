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
	DEFAULT_SECTION = "DEFAULT"
	bNumComment     = []byte{'#'} // number sign
	bSemComment     = []byte{';'} // semicolon
	bEmpty          = []byte{}
	bEqual          = []byte{'='}
	bDQuote         = []byte{'"'}
	sectionStart    = []byte{'['}
	sectionEnd      = []byte{']'}
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
		make(map[string]map[string]string),
		make(map[string]string),
		make(map[string]string),
		sync.RWMutex{},
	}
	cfg.Lock()
	defer cfg.Unlock()
	defer file.Close()

	var comment bytes.Buffer
	buf := bufio.NewReader(file)
	section := DEFAULT_SECTION
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if bytes.Equal(line, bEmpty) {
			continue
		}
		line = bytes.TrimSpace(line)

		var bComment []byte
		switch {
		case bytes.HasPrefix(line, bNumComment):
			bComment = bNumComment
		case bytes.HasPrefix(line, bSemComment):
			bComment = bSemComment
		}
		if bComment != nil {
			line = bytes.TrimLeft(line, string(bComment))
			line = bytes.TrimLeftFunc(line, unicode.IsSpace)
			comment.Write(line)
			comment.WriteByte('\n')
			continue
		}

		if bytes.HasPrefix(line, sectionStart) && bytes.HasSuffix(line, sectionEnd) {
			section = string(line[1 : len(line)-1])
			if comment.Len() > 0 {
				cfg.sectionComment[section] = comment.String()
				comment.Reset()
			}
			if _, ok := cfg.data[section]; !ok {
				cfg.data[section] = make(map[string]string)
			}
		} else {
			if _, ok := cfg.data[section]; !ok {
				cfg.data[section] = make(map[string]string)
			}
			keyval := bytes.SplitN(line, bEqual, 2)
			val := bytes.TrimSpace(keyval[1])
			if bytes.HasPrefix(val, bDQuote) {
				val = bytes.Trim(val, `"`)
			}

			key := string(bytes.TrimSpace(keyval[0]))
			cfg.data[section][key] = string(val)
			if comment.Len() > 0 {
				cfg.keycomment[section+"."+key] = comment.String()
				comment.Reset()
			}
		}

	}
	return cfg, nil
}

// A Config represents the configuration.
type IniConfigContainer struct {
	filename       string
	data           map[string]map[string]string //section=> key:val
	sectionComment map[string]string            //sction : comment
	keycomment     map[string]string            // id: []{comment, key...}; id 1 is for main comment.
	sync.RWMutex
}

// Bool returns the boolean value for a given key.
func (c *IniConfigContainer) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.getdata(key))
}

// Int returns the integer value for a given key.
func (c *IniConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.getdata(key))
}

func (c *IniConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.getdata(key), 10, 64)
}

// Float returns the float value for a given key.
func (c *IniConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.getdata(key), 64)
}

// String returns the string value for a given key.
func (c *IniConfigContainer) String(key string) string {
	return c.getdata(key)
}

// WriteValue writes a new value for key.
func (c *IniConfigContainer) Set(key, value string) error {
	c.Lock()
	defer c.Unlock()
	if len(key) == 0 {
		return errors.New("key is empty")
	}
	var section, k string
	sectionkey := strings.Split(key, ".")
	if len(sectionkey) >= 2 {
		section = sectionkey[0]
		k = sectionkey[1]
	} else {
		section = DEFAULT_SECTION
		k = sectionkey[0]
	}
	c.data[section][k] = value
	return nil
}

func (c *IniConfigContainer) DIY(key string) (v interface{}, err error) {
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return v, errors.New("key not find")
}

//section.key or key
func (c *IniConfigContainer) getdata(key string) string {
	c.RLock()
	defer c.RUnlock()
	if len(key) == 0 {
		return ""
	}
	var section, k string
	sectionkey := strings.Split(key, ".")
	if len(sectionkey) >= 2 {
		section = sectionkey[0]
		k = sectionkey[1]
	} else {
		section = DEFAULT_SECTION
		k = sectionkey[0]
	}
	if v, ok := c.data[section]; ok {
		if vv, o := v[k]; o {
			return vv
		}
	}
	return ""
}

func init() {
	Register("ini", &IniConfig{})
}
