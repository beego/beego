// Copyright 2012-2013 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

//	x2m_bulk.go: Process files with multiple XML messages.

package x2j

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"regexp"
	"sync"
)

// XmlMsgsFromFile()
//	'fname' is name of file
//	'phandler' is the map processing handler. Return of 'false' stops further processing.
//	'ehandler' is the parsing error handler. Return of 'false' stops further processing and returns error.
//	Note: phandler() and ehandler() calls are blocking, so reading and processing of messages is serialized.
//	      This means that you can stop reading the file on error or after processing a particular message.
//	      To have reading and handling run concurrently, pass arguments to a go routine in handler and return true.
func XmlMsgsFromFile(fname string, phandler func(map[string]interface{})(bool), ehandler func(error)(bool), recast ...bool) error {
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
		m, merr := XmlBufferToMap(b,r)
		if merr != nil && merr != io.EOF {
			if ok := ehandler(merr); !ok {
				// caused reader termination
				return merr
			 }
		}
		if m != nil {
			if ok := phandler(m); !ok {
				break
			}
		}
		if merr == io.EOF {
			break
		}
	}
	return nil
}

// XmlBufferToMap - process XML message from a bytes.Buffer
//	'b' is the buffer
//	Optional argument 'recast' coerces map values to float64 or bool where possible.
func XmlBufferToMap(b *bytes.Buffer,recast ...bool) (map[string]interface{},error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}

	n,err := XmlBufferToTree(b)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m[n.key] = n.treeToMap(r)

	return m,nil
}

// BufferToTree - derived from DocToTree()
func XmlBufferToTree(b *bytes.Buffer) (*Node, error) {
	p := xml.NewDecoder(b)
	p.CharsetReader = X2jCharsetReader
	n, berr := xmlToTree("",nil,p)
	if berr != nil {
		return nil, berr
	}

	return n,nil
}

// XmlBuffer - create XML decoder buffer for a string from anywhere, not necessarily a file.
type XmlBuffer struct {
	cnt uint64
	str *string
	buf *bytes.Buffer
}
var mtx sync.Mutex
var cnt uint64
var activeXmlBufs = make(map[uint64]*XmlBuffer)

// NewXmlBuffer() - creates a bytes.Buffer from a string with multiple messages
//	Use Close() function to release the buffer for garbage collection.
func NewXmlBuffer(s string) *XmlBuffer {
	// xml.Decoder doesn't properly handle whitespace in some doc
	// see songTextString.xml test case ... 
	reg,_ := regexp.Compile("[ \t\n\r]*<")
	s = reg.ReplaceAllString(s,"<")
	b := bytes.NewBufferString(s)
	buf := new(XmlBuffer)
	buf.str = &s
	buf.buf = b
	mtx.Lock()
	defer mtx.Unlock()
	buf.cnt = cnt ; cnt++
	activeXmlBufs[buf.cnt] = buf
	return buf
}

// BytesNewXmlBuffer() - creates a bytes.Buffer from b with possibly multiple messages
//	Use Close() function to release the buffer for garbage collection.
func BytesNewXmlBuffer(b []byte) *XmlBuffer {
	bb := bytes.NewBuffer(b)
	buf := new(XmlBuffer)
	buf.buf = bb
	mtx.Lock()
	defer mtx.Unlock()
	buf.cnt = cnt ; cnt++
	activeXmlBufs[buf.cnt] = buf
	return buf
}

// Close() - release the buffer address for garbage collection
func (buf *XmlBuffer)Close() {
	mtx.Lock()
	defer mtx.Unlock()
	delete(activeXmlBufs,buf.cnt)
}

// NextMap() - retrieve next XML message in buffer as a map[string]interface{} value.
//	The optional argument 'recast' will try and coerce values to float64 or bool as appropriate.
func (buf *XmlBuffer)NextMap(recast ...bool) (map[string]interface{}, error) {
		var r bool
		if len(recast) == 1 {
			r = recast[0]
		}
		if _, ok := activeXmlBufs[buf.cnt]; !ok {
			return nil, errors.New("Buffer is not active.")
		}
		return XmlBufferToMap(buf.buf,r)
}


// =============================  io.Reader version for stream processing  ======================

// XmlMsgsFromReader() - io.Reader version of XmlMsgsFromFile
//	'rdr' is an io.Reader for an XML message (stream)
//	'phandler' is the map processing handler. Return of 'false' stops further processing.
//	'ehandler' is the parsing error handler. Return of 'false' stops further processing and returns error.
//	Note: phandler() and ehandler() calls are blocking, so reading and processing of messages is serialized.
//	      This means that you can stop reading the file on error or after processing a particular message.
//	      To have reading and handling run concurrently, pass arguments to a go routine in handler and return true.
func XmlMsgsFromReader(rdr io.Reader, phandler func(map[string]interface{})(bool), ehandler func(error)(bool), recast ...bool) error {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}

	for {
		m, merr := ToMap(rdr,r)
		if merr != nil && merr != io.EOF {
			if ok := ehandler(merr); !ok {
				// caused reader termination
				return merr
			 }
		}
		if m != nil {
			if ok := phandler(m); !ok {
				break
			}
		}
		if merr == io.EOF {
			break
		}
	}
	return nil
}

