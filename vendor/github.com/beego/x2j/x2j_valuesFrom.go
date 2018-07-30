// Copyright 2012-2013 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

//	x2j_valuesFrom.go: Extract values from an arbitrary XML doc. Tag path can include wildcard characters.

package x2j

import (
	"strings"
)

// ------------------- sweep up everything for some point in the node tree ---------------------

// ValuesFromTagPath - deliver all values for a path node from a XML doc
// If there are no values for the path 'nil' is returned.
// A return value of (nil, nil) means that there were no values and no errors parsing the doc.
//   'doc' is the XML document
//   'path' is a dot-separated path of tag nodes
//   'getAttrs' can be set 'true' to return attribute values for "*"-terminated path
//          If a node is '*', then everything beyond is scanned for values.
//          E.g., "doc.books' might return a single value 'book' of type []interface{}, but
//                "doc.books.*" could return all the 'book' entries as []map[string]interface{}.
//                "doc.books.*.author" might return all the 'author' tag values as []string - or
//            		"doc.books.*.author.lastname" might be required, depending on he schema.
func ValuesFromTagPath(doc, path string, getAttrs ...bool) ([]interface{}, error) {
	var a bool
	if len(getAttrs) == 1 {
		a = getAttrs[0]
	}
	m, err := DocToMap(doc)
	if err != nil {
		return nil, err
	}

	v := ValuesFromKeyPath(m, path, a)
	return v, nil
}

// ValuesFromKeyPath - deliver all values for a path node from a map[string]interface{}
// If there are no values for the path 'nil' is returned.
//   'm' is the map to be walked
//   'path' is a dot-separated path of key values
//   'getAttrs' can be set 'true' to return attribute values for "*"-terminated path
//          If a node is '*', then everything beyond is walked.
//          E.g., see ValuesFromTagPath documentation.
func ValuesFromKeyPath(m map[string]interface{}, path string, getAttrs ...bool) []interface{} {
	var a bool
	if len(getAttrs) == 1 {
		a = getAttrs[0]
	}
	keys := strings.Split(path, ".")
	ret := make([]interface{}, 0)
	valuesFromKeyPath(&ret, m, keys, a)
	if len(ret) == 0 {
		return nil
	}
	return ret
}

func valuesFromKeyPath(ret *[]interface{}, m interface{}, keys []string, getAttrs bool) {
	lenKeys := len(keys)

	// load 'm' values into 'ret'
	// expand any lists
	if lenKeys == 0 {
		switch m.(type) {
		case map[string]interface{}:
			*ret = append(*ret, m)
		case []interface{}:
			for _, v := range m.([]interface{}) {
				*ret = append(*ret, v)
			}
		default:
			*ret = append(*ret, m)
		}
		return
	}

	// key of interest
	key := keys[0]
	switch key {
	case "*": // wildcard - scan all values
		switch m.(type) {
		case map[string]interface{}:
			for k, v := range m.(map[string]interface{}) {
				if string(k[:1]) == "-" && !getAttrs { // skip attributes?
					continue
				}
				valuesFromKeyPath(ret, v, keys[1:], getAttrs)
			}
		case []interface{}:
			for _, v := range m.([]interface{}) {
				switch v.(type) {
				// flatten out a list of maps - keys are processed
				case map[string]interface{}:
					for kk, vv := range v.(map[string]interface{}) {
						if string(kk[:1]) == "-" && !getAttrs { // skip attributes?
							continue
						}
						valuesFromKeyPath(ret, vv, keys[1:], getAttrs)
					}
				default:
					valuesFromKeyPath(ret, v, keys[1:], getAttrs)
				}
			}
		}
	default: // key - must be map[string]interface{}
		switch m.(type) {
		case map[string]interface{}:
			if v, ok := m.(map[string]interface{})[key]; ok {
				valuesFromKeyPath(ret, v, keys[1:], getAttrs)
			}
		case []interface{}: // may be buried in list
			for _, v := range m.([]interface{}) {
				switch v.(type) {
				case map[string]interface{}:
					if vv, ok := v.(map[string]interface{})[key]; ok {
						valuesFromKeyPath(ret, vv, keys[1:], getAttrs)
					}
				}
			}
		}
	}
}
