package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SLACKWriter implements beego LoggerInterface and is used to send jiaoliao webhook
type SLACKWriter struct {
	WebhookURL string `json:"webhookurl"`
	Level      int    `json:"level"`
	formatter  LogFormatter
	Formatter  string `json:"formatter"`
}

// newSLACKWriter creates jiaoliao writer.
func newSLACKWriter() Logger {
	res := &SLACKWriter{Level: LevelTrace}
	res.formatter = res
	return res
}

func (s *SLACKWriter) Format(lm *LogMsg) string {
	// text := fmt.Sprintf("{\"text\": \"%s\"}", msg)
	return lm.When.Format("2006-01-02 15:04:05") + " " + lm.OldStyleFormat()
}

func (s *SLACKWriter) SetFormatter(f LogFormatter) {
	s.formatter = f
}

// Init SLACKWriter with json config string
func (s *SLACKWriter) Init(config string) error {
	res := json.Unmarshal([]byte(config), s)

	if res == nil && len(s.Formatter) > 0 {
		fmtr, ok := GetFormatter(s.Formatter)
		if !ok {
			return fmt.Errorf("the formatter with name: %s not found", s.Formatter)
		}
		s.formatter = fmtr
	}

	return res
}

// WriteMsg write message in smtp writer.
// Sends an email with subject and only this message.
func (s *SLACKWriter) WriteMsg(lm *LogMsg) error {
	if lm.Level > s.Level {
		return nil
	}
	msg := s.Format(lm)
	m := make(map[string]string, 1)
	m["text"] = msg

	body, _ := json.Marshal(m)
	// resp, err := http.PostForm(s.WebhookURL, form)
	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewReader(body))
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
