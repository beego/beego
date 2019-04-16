package goes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Requester implements Request which builds an HTTP request for Elasticsearch
type Requester interface {
	// Request should set the URL and Body (if needed). The host of the URL will be overwritten by the client.
	Request() (*http.Request, error)
}

// Request holds a single request to elasticsearch
type Request struct {
	// A search query
	Query interface{}

	// Which index to search into
	IndexList []string

	// Which type to search into
	TypeList []string

	// HTTP Method to user (GET, POST ...)
	Method string

	// Which api keyword (_search, _bulk, etc) to use
	API string

	// Bulk data
	BulkData []byte

	// Request body
	Body []byte

	// A list of extra URL arguments
	ExtraArgs url.Values

	// Used for the id field when indexing a document
	ID string
}

// URL builds a URL for a Request
func (req *Request) URL() *url.URL {
	var path string
	if len(req.IndexList) > 0 {
		path = "/" + strings.Join(req.IndexList, ",")
	}

	if len(req.TypeList) > 0 {
		path += "/" + strings.Join(req.TypeList, ",")
	}

	// XXX : for indexing documents using the normal (non bulk) API
	if len(req.ID) > 0 {
		path += "/" + req.ID
	}

	path += "/" + req.API

	u := url.URL{
		//Scheme:   "http",
		//Host:     fmt.Sprintf("%s:%s", req.Conn.Host, req.Conn.Port),
		Path:     path,
		RawQuery: req.ExtraArgs.Encode(),
	}

	return &u
}

// Request generates an http.Request based on the contents of the Request struct
func (req *Request) Request() (*http.Request, error) {
	postData := []byte{}

	// XXX : refactor this
	if len(req.Body) > 0 {
		postData = req.Body
	} else if req.API == "_bulk" {
		postData = req.BulkData
	} else if req.Query != nil {
		b, err := json.Marshal(req.Query)
		if err != nil {
			return nil, err
		}
		postData = b
	}

	newReq, err := http.NewRequest(req.Method, "", nil)
	if err != nil {
		return nil, err
	}
	newReq.URL = req.URL()
	newReq.Body = ioutil.NopCloser(bytes.NewReader(postData))
	newReq.ContentLength = int64(len(postData))

	if req.Method == "POST" || req.Method == "PUT" {
		newReq.Header.Set("Content-Type", "application/json")
	}
	return newReq, nil
}

var _ Requester = (*Request)(nil)
