package logs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// JLWriter implements beego LoggerInterface and is used to send jiaoliao webhook
type JLWriter struct {
	AuthorName  string `json:"authorname"`
	Title       string `json:"title"`
	WebhookURL  string `json:"webhookurl"`
	RedirectURL string `json:"redirecturl,omitempty"`
	ImageURL    string `json:"imageurl,omitempty"`
	Level       int    `json:"level"`

	formatter LogFormatter
	Formatter string `json:"formatter"`
}

// newJLWriter creates jiaoliao writer.
func newJLWriter() Logger {
	res := &JLWriter{Level: LevelTrace}
	res.formatter = res
	return res
}

// Init JLWriter with json config string
func (s *JLWriter) Init(config string) error {
	res := json.Unmarshal([]byte(config), s)
	if res == nil && len(s.Formatter) > 0 {
		fmtr, ok := GetFormatter(s.Formatter)
		if !ok {
			return errors.New(fmt.Sprintf("the formatter with name: %s not found", s.Formatter))
		}
		s.formatter = fmtr
	}
	return res
}

func (s *JLWriter) Format(lm *LogMsg) string {
	msg := lm.OldStyleFormat()
	msg = fmt.Sprintf("%s %s", lm.When.Format("2006-01-02 15:04:05"), msg)
	return msg
}

func (s *JLWriter) SetFormatter(f LogFormatter) {
	s.formatter = f
}

// WriteMsg writes message in smtp writer.
// Sends an email with subject and only this message.
func (s *JLWriter) WriteMsg(lm *LogMsg) error {
	if lm.Level > s.Level {
		return nil
	}

	text := s.formatter.Format(lm)

	form := url.Values{}
	form.Add("authorName", s.AuthorName)
	form.Add("title", s.Title)
	form.Add("text", text)
	if s.RedirectURL != "" {
		form.Add("redirectUrl", s.RedirectURL)
	}
	if s.ImageURL != "" {
		form.Add("imageUrl", s.ImageURL)
	}

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
func (s *JLWriter) Flush() {
}

// Destroy implementing method. empty.
func (s *JLWriter) Destroy() {
}

func init() {
	Register(AdapterJianLiao, newJLWriter)
}
