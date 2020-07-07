package es

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"

	"github.com/astaxie/beego/logs"
)

// NewES return a LoggerInterface
func NewES() logs.Logger {
	cw := &esLogger{
		Level: logs.LevelDebug,
	}
	return cw
}

// esLogger will log msg into ES
// before you using this implementation,
// please import this package
// usually means that you can import this package in your main package
// for example, anonymous:
// import _ "github.com/astaxie/beego/logs/es"
type esLogger struct {
	*elasticsearch.Client
	DSN         string `json:"dsn"`
	Level       int    `json:"level"`
	Index       string `json:"index"`
	IndexPrefix string `json:"index_prefix"`
}

// {"dsn":"http://localhost:9200/","index":"index-name","index_prefix":"test-index-","level":1}
func (el *esLogger) Init(jsonconfig string) error {
	err := json.Unmarshal([]byte(jsonconfig), el)
	if err != nil {
		return err
	}
	if el.IndexPrefix == "" && beego.BConfig.AppName != "" {
		el.IndexPrefix = beego.BConfig.AppName + "-"
	}
	if el.DSN == "" {
		return errors.New("empty dsn")
	} else if u, err := url.Parse(el.DSN); err != nil {
		return err
	} else if u.Path == "" {
		return errors.New("missing prefix")
	} else {
		conn, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{el.DSN},
		})
		if err != nil {
			return err
		}
		el.Client = conn
	}
	return nil
}

// WriteMsg will write the msg and level into es
func (el *esLogger) WriteMsg(when time.Time, msg string, level int) error {
	if level > el.Level {
		return nil
	}

	idx := LogDocument{
		Timestamp: when.Format(time.RFC3339),
		Msg:       msg,
	}

	body, err := json.Marshal(idx)
	if err != nil {
		return err
	}
	req := esapi.IndexRequest{
		Index:        el.IndexPrefix + fmt.Sprintf("%04d-%02d-%02d", when.Year(), when.Month(), when.Day()),
		DocumentType: "logs",
		Body:         strings.NewReader(string(body)),
	}
	if el.Index != "" {
		req.Index = el.Index
	}
	_, err = req.Do(context.Background(), el.Client)
	return err
}

// Destroy is a empty method
func (el *esLogger) Destroy() {
}

// Flush is a empty method
func (el *esLogger) Flush() {

}

type LogDocument struct {
	Timestamp string      `json:"timestamp"`
	Msg       interface{} `json:"msg"`
}

func init() {
	logs.Register(logs.AdapterEs, NewES)
}
