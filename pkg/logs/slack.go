package logs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/astaxie/beego/pkg/common"
)

// SLACKWriter implements beego LoggerInterface and is used to send jiaoliao webhook
type SLACKWriter struct {
	WebhookURL string `json:"webhookurl"`
	Level      int    `json:"level"`
}

// newSLACKWriter create jiaoliao writer.
func newSLACKWriter() Logger {
	return &SLACKWriter{Level: LevelTrace}
}

// Init SLACKWriter with json config string
func (s *SLACKWriter) Init(jsonconfig string) error {
	return json.Unmarshal([]byte(jsonconfig), s)
}

// WriteMsg write message in smtp writer.
// it will send an email with subject and only this message.
func (s *SLACKWriter) WriteMsg(lm *LogMsg, opts ...common.SimpleKV) error {
	if lm.Level > s.Level {
		return nil
	}

	msg := ""
	for _, elem := range opts {
		if elem.Key == "formatterFunc" {
			formatterFunc := elem.Value.(func(*LogMsg) string)
			msg = formatterFunc(lm)
		} else {
			msg = lm.Msg
		}
	}

	text := fmt.Sprintf("{\"text\": \"%s %s\"}", lm.When.Format("2006-01-02 15:04:05"), msg)

	form := url.Values{}
	form.Add("payload", text)

	resp, err := http.PostForm(s.WebhookURL, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Post webhook failed %s %d", resp.Status, resp.StatusCode)
	}
	return nil
}

// Flush implementing method. empty.
func (s *SLACKWriter) Flush() {
}

// Destroy implementing method. empty.
func (s *SLACKWriter) Destroy() {
}

func init() {
	Register(AdapterSlack, newSLACKWriter)
}
