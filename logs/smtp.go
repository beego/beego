package logs

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

const (
	subjectPhrase = "Diagnostic message from server"
)

// smtpWriter is used to send emails via given SMTP-server.
type SmtpWriter struct {
	Username           string   `json:"Username"`
	Password           string   `json:"password"`
	Host               string   `json:"Host"`
	Subject            string   `json:"subject"`
	RecipientAddresses []string `json:"sendTos"`
	Level              int      `json:"level"`
}

func NewSmtpWriter() LoggerInterface {
	return &SmtpWriter{Level: LevelTrace}
}

func (s *SmtpWriter) Init(jsonconfig string) error {
	err := json.Unmarshal([]byte(jsonconfig), s)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmtpWriter) WriteMsg(msg string, level int) error {
	if level < s.Level {
		return nil
	}

	hp := strings.Split(s.Host, ":")

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		s.Username,
		s.Password,
		hp[0],
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	content_type := "Content-Type: text/plain" + "; charset=UTF-8"
	mailmsg := []byte("To: " + strings.Join(s.RecipientAddresses, ";") + "\r\nFrom: " + s.Username + "<" + s.Username +
		">\r\nSubject: " + s.Subject + "\r\n" + content_type + "\r\n\r\n" + fmt.Sprintf(".%s", time.Now().Format("2006-01-02 15:04:05")) + msg)

	err := smtp.SendMail(
		s.Host,
		auth,
		s.Username,
		s.RecipientAddresses,
		mailmsg,
	)

	return err
}

func (s *SmtpWriter) Flush() {
	return
}
func (s *SmtpWriter) Destroy() {
	return
}

func init() {
	Register("smtp", NewSmtpWriter)
}
