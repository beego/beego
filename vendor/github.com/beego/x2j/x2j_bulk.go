// Copyright 2012-2013 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

//	x2j_bulk.go: Process files with multiple XML messages.
// Extends x2m_bulk.go to work with JSON strings rather than map[string]interface{}.

package x2j

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"regexp"
)

// XmlMsgsFromFileAsJson()
//	'fname' is name of file
//	'phandler' is the JSON string processing handler. Return of 'false' stops further processing.
//	'ehandler' is the parsing error handler. Return of 'false' stops further processing and returns error.
//	Note: phandler() and ehandler() calls are blocking, so reading and processing of messages is serialized.
//	      This means that you can stop reading the file on error or after processing a particular message.
//	      To have reading and handling run concurrently, pass arguments to a go routine in handler and return true.
func XmlMsgsFromFileAsJson(fname string, phandler func(string)(bool), ehandler func(error)(bool), recast ...bool) error {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	fi, fierr := os.Stat(fname)
	if fierr != nil {
		return fierr
	}
	fh, fherr := os.Open(fname)
	if fherr != nil {
		return fherr
	}
	defer fh.Close()
	buf := make([]byte,fi.Size())
	_, rerr  :=  fh.Read(buf)
	if rerr != nil {
		return rerr
	}
	doc := string(buf)

	// xml.Decoder doesn't properly handle whitespace in some doc
	// see songTextString.xml test case ... 
	reg,_ := regexp.Compile("[ \t\n\r]*<")
	doc = reg.ReplaceAllString(doc,"<")
	b := bytes.NewBufferString(doc)

	for {
		s, serr := XmlBufferToJson(b,r)
		if serr != nil && serr != io.EOF {
			if ok := ehandler(serr); !ok {
				// caused reader termination
				return serr
			 }
		}
		if s != "" {
			if ok := phandler(s); !ok {
				break
			}
		}
		if serr == io.EOF {
			break
		}
	}
	return nil
}

// XmlBufferToJson - process XML message from a bytes.Buffer
//	'b' is the buffer
//	Optional argument 'recast' coerces values to float64 or bool where possible.
func XmlBufferToJson(b *bytes.Buffer,recast ...bool) (string,error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}

	n,err := XmlBufferToTree(b)
	if err != nil {
		return "", err
	}

	m := make(map[string]interface{})
	m[n.key] = n.treeToMap(r)

	j, jerr := json.Marshal(m)
	return string(j), jerr
}

// =============================  io.Reader version for stream processing  ======================

// XmlMsgsFromReaderAsJson() - io.Reader version of XmlMsgsFromFileAsJson
//	'rdr' is an io.Reader for an XML message (stream)
//	'phandler' is the JSON string processing handler. Return of 'false' stops further processing.
//	'ehandler' is the parsing error handler. Return of 'false' stops further processing and returns error.
//	Note: phandler() and ehandler() calls are blocking, so reading and processing of messages is serialized.
//	      This means that you can stop reading the file on error or after processing a particular message.
//	      To have reading and handling run concurrently, pass arguments to a go routine in handler and return true.
func XmlMsgsFromReaderAsJson(rdr io.Reader, phandler func(string)(bool), ehandler func(error)(bool), recast ...bool) error {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}

	for {
		s, serr := ToJson(rdr,r)
		if serr != nil && serr != io.EOF {
			if ok := ehandler(serr); !ok {
				// caused reader termination
				return serr
			 }
		}
		if s != "" {
			if ok := phandler(s); !ok {
				break
			}
		}
		if serr == io.EOF {
			break
		}
	}
	return nil
}

