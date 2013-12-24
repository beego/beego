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
	DEFAULT_SECTION = "default"   // default section means if some ini items not in a section, make them in default section,
	bNumComment     = []byte{'#'} // number signal
	bSemComment     = []byte{';'} // semicolon signal
	bEmpty          = []byte{}
	bEqual          = []byte{'='} // equal signal
	bDQuote         = []byte{'"'} // quote signal
	sectionStart    = []byte{'['} // section start signal
	sectionEnd      = []byte{']'} // section end signal
)

// IniConfig implements Config to parse ini file.
type IniConfig struct {
}

// ParseFile creates a new Config and parses the file configuration from the named file.
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
			section = strings.ToLower(section) // section name case insensitive
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

			key := string(bytes.TrimSpace(keyval[0])) // key name case insensitive
			key = strings.ToLower(key)
			cfg.data[section][key] = string(val)
			if comment.Len() > 0 {
				cfg.keycomment[section+"."+key] = comment.String()
				comment.Reset()
			}
		}

	}
	return cfg, nil
}

// A Config represents the ini configuration.
// When set and get value, support key as section:name type.
type IniConfigContainer struct {
	filename       string
	data           map[string]map[string]string // section=> key:val
	sectionComment map[string]string            // section : comment
	keycomment     map[string]string            // id: []{comment, key...}; id 1 is for main comment.
	sync.RWMutex
}

// Bool returns the boolean value for a given key.
func (c *IniConfigContainer) Bool(key string) (bool, error) {
	key = strings.ToLower(key)
	return strconv.ParseBool(c.getdata(key))
}

// Int returns the integer value for a given key.
func (c *IniConfigContainer) Int(key string) (int, error) {
	key = strings.ToLower(key)
	return strconv.Atoi(c.getdata(key))
}

// Int64 returns the int64 value for a given key.
func (c *IniConfigContainer) Int64(key string) (int64, error) {
	key = strings.ToLower(key)
	return strconv.ParseInt(c.getdata(key), 10, 64)
}

// Float returns the float value for a given key.
func (c *IniConfigContainer) Float(key string) (float64, error) {
	key = strings.ToLower(key)
	return strconv.ParseFloat(c.getdata(key), 64)
}

// String returns the string value for a given key.
func (c *IniConfigContainer) String(key string) string {
	key = strings.ToLower(key)
	return c.getdata(key)
}

// WriteValue writes a new value for key.
// if write to one section, the key need be "section::key".
// if the section is not existed, it panics.
func (c *IniConfigContainer) Set(key, value string) error {
	c.Lock()
	defer c.Unlock()
	if len(key) == 0 {
		return errors.New("key is empty")
	}

	var section, k string
	key = strings.ToLower(key)
	sectionkey := strings.Split(key, "::")
	if len(sectionkey) >= 2 {
		section = sectionkey[0]
		k = sectionkey[1]
	} else {
		section = DEFAULT_SECTION
		k = sectionkey[0]
	}
	if _, ok := c.data[section]; !ok {
		c.data[section] = make(map[string]string)
	}
	c.data[section][k] = value
	return nil
}

// DIY returns the raw value by a given key.
func (c *IniConfigContainer) DIY(key string) (v interface{}, err error) {
	key = strings.ToLower(key)
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return v, errors.New("key not find")
}

// section.key or key
func (c *IniConfigContainer) getdata(key string) string {
	c.RLock()
	defer c.RUnlock()
	if len(key) == 0 {
		return ""
	}

	var section, k string
	key = strings.ToLower(key)
	sectionkey := strings.Split(key, "::")
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
