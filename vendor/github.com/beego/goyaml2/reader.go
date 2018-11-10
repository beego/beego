package goyaml2

import (
	"bufio"
	"fmt"
	"github.com/wendal/errors"
	"io"
	"log"
	"strings"
)

const (
	DEBUG        = true
	MAP_KEY_ONLY = iota
)

func Read(r io.Reader) (interface{}, error) {
	yr := &yamlReader{}
	yr.br = bufio.NewReader(r)
	obj, err := yr.ReadObject(0)
	if err == io.EOF {
		err = nil
	}
	if obj == nil {
		log.Println("Obj == nil")
	}
	return obj, err
}

type yamlReader struct {
	br       *bufio.Reader
	nodes    []interface{}
	lineNum  int
	lastLine string
}

func (y *yamlReader) ReadObject(minIndent int) (interface{}, error) {
	line, err := y.NextLine()
	if err != nil {
		if err == io.EOF && line != "" {
			//log.Println("Read EOF , but still some data here")
		} else {
			//log.Println("ReadERR", err)
			return nil, err
		}
	}
	y.lastLine = line
	indent, str := getIndent(line)
	if indent < minIndent {
		//log.Println("Current Indent Unexpect : ", str, indent, minIndent)
		return nil, y.Error("Unexpect Indent", nil)
	}
	if indent > minIndent {
		//log.Println("Change minIndent from %d to %d", minIndent, indent)
		minIndent = indent
	}
	switch str[0] {
	case '-':
		return y.ReadList(minIndent)
	case '[':
		fallthrough
	case '{':
		y.lastLine = ""
		_, value, err := y.asMapKeyValue("tmp:" + str)
		if err != nil {
			return nil, y.Error("Err inline map/list", nil)
		}
		return value, nil
	}
	//log.Println("Read Objcet as Map", indent, str)

	return y.ReadMap(minIndent)

}

func (y *yamlReader) ReadList(minIndent int) ([]interface{}, error) {
	list := []interface{}{}
	for {
		line, err := y.NextLine()
		if err != nil {
			return list, err
		}
		indent, str := getIndent(line)
		switch {
		case indent < minIndent:
			y.lastLine = line
			if len(list) == 0 {
				return nil, nil
			}
			return list, nil
		case indent == minIndent:
			if str[0] != '-' {
				y.lastLine = line
				return list, nil
			}
			if len(str) < 2 {
				return nil, y.Error("ListItem is Emtry", nil)
			}
			key, value, err := y.asMapKeyValue(str[1:])
			if err != nil {
				return nil, err
			}
			switch value {
			case nil:
				list = append(list, key)
			case MAP_KEY_ONLY:
				return nil, y.Error("Not support List-Map yet", nil)
			default:
				_map := map[string]interface{}{key.(string): value}
				list = append(list, _map)

				_line, _err := y.NextLine()
				if _err != nil && _err != io.EOF {
					return nil, err
				}
				if _line == "" {
					return list, nil
				}
				y.lastLine = _line
				_indent, _str := getIndent(line)
				if _indent >= minIndent+2 {
					switch _str[0] {
					case '-':
						return nil, y.Error("Unexpect", nil)
					case '[':
						return nil, y.Error("Unexpect", nil)
					case '{':
						return nil, y.Error("Unexpect", nil)
					}
					// look like a map
					_map2, _err := y.ReadMap(_indent)
					if _map2 != nil {
						_map2[key.(string)] = value
					}
					if err != nil {
						return list, _err
					}
				}
			}
			continue
		default:
			return nil, y.Error("Bad Indent\n"+line, nil)
		}
	}
	panic("ERROR")
	return nil, errors.New("Impossible")
}

func (y *yamlReader) ReadMap(minIndent int) (map[string]interface{}, error) {
	_map := map[string]interface{}{}
	//log.Println("ReadMap", minIndent)
OUT:
	for {
		line, err := y.NextLine()
		if err != nil {
			return _map, err
		}
		indent, str := getIndent(line)
		//log.Printf("Indent : %d, str = %s", indent, str)
		switch {
		case indent < minIndent:
			y.lastLine = line
			if len(_map) == 0 {
				return nil, nil
			}
			return _map, nil
		case indent == minIndent:
			key, value, err := y.asMapKeyValue(str)
			if err != nil {
				return nil, err
			}
			//log.Println("Key=", key, "value=", value)
			switch value {
			case nil:
				return nil, y.Error("Unexpect", nil)
			case MAP_KEY_ONLY:
				//log.Println("KeyOnly, read inner Map", key)

				//--------------------------------------
				_line, err := y.NextLine()
				if err != nil {
					if err == io.EOF {
						if _line == "" {
							// Emtry map item?
							_map[key.(string)] = nil
							return _map, err
						}
					} else {
						return nil, y.Error("ERR?", err)
					}
				}
				y.lastLine = _line
				_indent, _str := getIndent(_line)
				if _indent < minIndent {
					return _map, nil
				}
				////log.Println("##>>", _indent, _str)
				if _indent == minIndent {
					if _str[0] == '-' {
						//log.Println("Read Same-Indent ListItem for Map")
						_list, err := y.ReadList(minIndent)
						if _list != nil {
							_map[key.(string)] = _list
						}
						if err != nil {
							return _map, nil
						}
						continue OUT
					} else {
						// Emtry map item?
						_map[key.(string)] = nil
						continue OUT
					}
				}
				//--------------------------------------
				//log.Println("Read Map Item", _indent, _str)

				obj, err := y.ReadObject(_indent)
				if obj != nil {
					_map[key.(string)] = obj
				}
				if err != nil {
					return _map, err
				}
			default:
				_map[key.(string)] = value
			}
		default:
			//log.Println("Bad", indent, str)
			return nil, y.Error("Bad Indent\n"+line, nil)
		}
	}
	panic("ERROR")
	return nil, errors.New("Impossible")

}

func (y *yamlReader) NextLine() (line string, err error) {
	if y.lastLine != "" {
		line = y.lastLine
		y.lastLine = ""
		//log.Println("Return lastLine", line)
		return
	}
	for {
		line, err = y.br.ReadString('\n')
		y.lineNum++
		if err != nil {
			return
		}
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimRight(line, "\n\t\r ")
		if line == "" {
			continue
		}
		//log.Println("Return Line", line)
		return
	}
	//log.Println("Impossbible : " + line)
	return // impossbile!
}

func getIndent(str string) (int, string) {
	indent := 0
	for i, s := range str {
		switch s {
		case ' ':
			indent++
		case '\t':
			indent += 4
		default:
			return indent, str[i:]
		}
	}
	panic("Invalid indent : " + str)
	return -1, ""
}

func (y *yamlReader) asMapKeyValue(str string) (key interface{}, val interface{}, err error) {
	tokens := splitToken(str)
	key = tokens[0]
	if len(tokens) == 1 {
		return key, nil, nil
	}
	if tokens[1] != ":" {
		return "", nil, y.Error("Unexpect "+str, nil)
	}

	if len(tokens) == 2 {
		return key, MAP_KEY_ONLY, nil
	}
	if len(tokens) == 3 {
		return key, tokens[2], nil
	}
	switch tokens[2] {
	case "[":
		list := []interface{}{}
		for i := 3; i < len(tokens)-1; i++ {
			list = append(list, tokens[i])
		}
		return key, list, nil
	case "{":
		_map := map[string]interface{}{}
		for i := 3; i < len(tokens)-1; i += 4 {
			//log.Println(">>>", i, tokens[i])
			if i > len(tokens)-2 {
				return "", nil, y.Error("Unexpect "+str, nil)
			}
			if tokens[i+1] != ":" {
				return "", nil, y.Error("Unexpect "+str, nil)
			}
			_map[tokens[i].(string)] = tokens[i+2]
			if (i + 3) < (len(tokens) - 1) {
				if tokens[i+3] != "," {
					return "", "", y.Error("Unexpect "+str, nil)
				}
			} else {
				break
			}
		}
		return key, _map, nil
	}
	//log.Println(str, tokens)
	return "", nil, y.Error("Unexpect "+str, nil)
}

func splitToken(str string) (tokens []interface{}) {
	str = strings.Trim(str, "\r\t\n ")
	if str == "" {
		panic("Impossbile")
		return
	}

	tokens = []interface{}{}
	lastPos := 0
	for i := 0; i < len(str); i++ {
		switch str[i] {
		case ':':
			fallthrough
		case '{':
			fallthrough
		case '[':
			fallthrough
		case '}':
			fallthrough
		case ']':
			fallthrough
		case ',':
			if i > lastPos {
				tokens = append(tokens, str[lastPos:i])
			}
			tokens = append(tokens, str[i:i+1])
			lastPos = i + 1
		case ' ':
			if i > lastPos {
				tokens = append(tokens, str[lastPos:i])
			}
			lastPos = i + 1
		case '\'':
			//log.Println("Scan End of String")
			i++
			start := i
			for ; i < len(str); i++ {
				if str[i] == '\'' {
					//log.Println("Found End of String", start, i)
					break
				}
			}
			tokens = append(tokens, str[start:i])
			lastPos = i + 1
		case '"':
			i++
			start := i
			for ; i < len(str); i++ {
				if str[i] == '"' {
					break
				}
			}
			tokens = append(tokens, str[start:i])
			lastPos = i + 1
		}
	}
	////log.Println("last", lastPos)
	if lastPos < len(str) {
		tokens = append(tokens, str[lastPos:])
	}

	if len(tokens) == 1 {
		tokens[0] = string2Val(tokens[0].(string))
		return
	}

	if tokens[1] == ":" {
		if len(tokens) == 2 {
			return
		}
		if tokens[2] == "{" || tokens[2] == "[" {
			return
		}
		str = strings.Trim(strings.SplitN(str, ":", 2)[1], "\t ")
		if len(str) > 2 {
			if str[0] == '\'' && str[len(str)-1] == '\'' {
				str = str[1 : len(str)-1]
			} else if str[0] == '"' && str[len(str)-1] == '"' {
				str = str[1 : len(str)-1]
			}
		}
		val := string2Val(str)
		tokens = []interface{}{tokens[0], tokens[1], val}
		return
	}

	if len(str) > 2 {
		if str[0] == '\'' && str[len(str)-1] == '\'' {
			str = str[1 : len(str)-1]
		} else if str[0] == '"' && str[len(str)-1] == '"' {
			str = str[1 : len(str)-1]
		}
	}
	val := string2Val(str)
	tokens = []interface{}{val}
	return
}

func (y *yamlReader) Error(msg string, err error) error {
	if err != nil {
		return errors.New(fmt.Sprintf("line %d : %s : %v", y.lineNum, msg, err.Error()))
	}
	return errors.New(fmt.Sprintf("line %d >> %s", y.lineNum, msg))
}
