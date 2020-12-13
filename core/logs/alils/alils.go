package alils

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/beego/beego/core/logs"
)

const (
	// CacheSize sets the flush size
	CacheSize int = 64
	// Delimiter defines the topic delimiter
	Delimiter string = "##"
)

// Config is the Config for Ali Log
type Config struct {
	Project   string   `json:"project"`
	Endpoint  string   `json:"endpoint"`
	KeyID     string   `json:"key_id"`
	KeySecret string   `json:"key_secret"`
	LogStore  string   `json:"log_store"`
	Topics    []string `json:"topics"`
	Source    string   `json:"source"`
	Level     int      `json:"level"`
	FlushWhen int      `json:"flush_when"`
	Formatter string   `json:"formatter"`
}

// aliLSWriter implements LoggerInterface.
// Writes messages in keep-live tcp connection.
type aliLSWriter struct {
	store    *LogStore
	group    []*LogGroup
	withMap  bool
	groupMap map[string]*LogGroup
	lock     *sync.Mutex
	Config
	formatter logs.LogFormatter
}

// NewAliLS creates a new Logger
func NewAliLS() logs.Logger {
	alils := new(aliLSWriter)
	alils.Level = logs.LevelTrace
	alils.formatter = alils
	return alils
}

// Init parses config and initializes struct
func (c *aliLSWriter) Init(config string) error {
	err := json.Unmarshal([]byte(config), c)
	if err != nil {
		return err
	}

	if c.FlushWhen > CacheSize {
		c.FlushWhen = CacheSize
	}

	prj := &LogProject{
		Name:            c.Project,
		Endpoint:        c.Endpoint,
		AccessKeyID:     c.KeyID,
		AccessKeySecret: c.KeySecret,
	}

	store, err := prj.GetLogStore(c.LogStore)
	if err != nil {
		return err
	}

	c.store = store

	// Create default Log Group
	c.group = append(c.group, &LogGroup{
		Topic:  proto.String(""),
		Source: proto.String(c.Source),
		Logs:   make([]*Log, 0, c.FlushWhen),
	})

	// Create other Log Group
	c.groupMap = make(map[string]*LogGroup)
	for _, topic := range c.Topics {

		lg := &LogGroup{
			Topic:  proto.String(topic),
			Source: proto.String(c.Source),
			Logs:   make([]*Log, 0, c.FlushWhen),
		}

		c.group = append(c.group, lg)
		c.groupMap[topic] = lg
	}

	if len(c.group) == 1 {
		c.withMap = false
	} else {
		c.withMap = true
	}

	c.lock = &sync.Mutex{}

	if len(c.Formatter) > 0 {
		fmtr, ok := logs.GetFormatter(c.Formatter)
		if !ok {
			return errors.New(fmt.Sprintf("the formatter with name: %s not found", c.Formatter))
		}
		c.formatter = fmtr
	}

	return nil
}

func (c *aliLSWriter) Format(lm *logs.LogMsg) string {
	return lm.OldStyleFormat()
}

func (c *aliLSWriter) SetFormatter(f logs.LogFormatter) {
	c.formatter = f
}

// WriteMsg writes a message in connection.
// If connection is down, try to re-connect.
func (c *aliLSWriter) WriteMsg(lm *logs.LogMsg) error {
	if lm.Level > c.Level {
		return nil
	}

	var topic string
	var content string
	var lg *LogGroup
	if c.withMap {

		// Topic，LogGroup
		strs := strings.SplitN(lm.Msg, Delimiter, 2)
		if len(strs) == 2 {
			pos := strings.LastIndex(strs[0], " ")
			topic = strs[0][pos+1 : len(strs[0])]
			lg = c.groupMap[topic]
		}

		// send to empty Topic
		if lg == nil {
			lg = c.group[0]
		}
	} else {
		lg = c.group[0]
	}

	content = c.formatter.Format(lm)

	c1 := &LogContent{
		Key:   proto.String("msg"),
		Value: proto.String(content),
	}

	l := &Log{
		Time: proto.Uint32(uint32(lm.When.Unix())),
		Contents: []*LogContent{
			c1,
		},
	}

	c.lock.Lock()
	lg.Logs = append(lg.Logs, l)
	c.lock.Unlock()

	if len(lg.Logs) >= c.FlushWhen {
		c.flush(lg)
	}
	return nil
}

// Flush implementing method. empty.
func (c *aliLSWriter) Flush() {

	// flush all group
	for _, lg := range c.group {
		c.flush(lg)
	}
}

// Destroy destroy connection writer and close tcp listener.
func (c *aliLSWriter) Destroy() {
}

func (c *aliLSWriter) flush(lg *LogGroup) {

	c.lock.Lock()
	defer c.lock.Unlock()
	err := c.store.PutLogs(lg)
	if err != nil {
		return
	}

	lg.Logs = make([]*Log, 0, c.FlushWhen)
}

func init() {
	logs.Register(logs.AdapterAliLS, NewAliLS)
}
