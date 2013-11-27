package logs

import (
	"encoding/json"
	"errors"
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
	username           string
	password           string
	host               string
	subject            string
	recipientAddresses []string
	level              int
}

func NewSmtpWriter() LoggerInterface {
	return &SmtpWriter{level: LevelTrace}
}

func (s *SmtpWriter) Init(jsonconfig string) error {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonconfig), &m)
	if err != nil {
		return err
	}
	if username, ok := m["username"]; !ok {
		return errors.New("smtp config must have auth username")
	} else if password, ok := m["password"]; !ok {
		return errors.New("smtp config must have auth password")
	} else if hostname, ok := m["host"]; !ok {
		return errors.New("smtp config must have host like 'mail.example.com:25'")
	} else if sendTos, ok := m["sendTos"]; !ok {
		return errors.New("smtp config must have sendTos")
	} else {
		s.username = username.(string)
		s.password = password.(string)
		s.host = hostname.(string)
		for _, v := range sendTos.([]interface{}) {
			s.recipientAddresses = append(s.recipientAddresses, v.(string))
		}
	}

	if subject, ok := m["subject"]; ok {
		s.subject = subject.(string)
	} else {
		s.subject = subjectPhrase
	}
	if lv, ok := m["level"]; ok {
		s.level = int(lv.(float64))
	}
	return nil
}

func (s *SmtpWriter) WriteMsg(msg string, level int) error {
	if level < s.level {
		return nil
	}

	hp := strings.Split(s.host, ":")

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		s.username,
		s.password,
		hp[0],
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	content_type := "Content-Type: text/plain" + "; charset=UTF-8"
	mailmsg := []byte("To: " + strings.Join(s.recipientAddresses, ";") + "\r\nFrom: " + s.username + "<" + s.username +
		">\r\nSubject: " + s.subject + "\r\n" + content_type + "\r\n\r\n" + fmt.Sprintf(".%s", time.Now().Format("2006-01-02 15:04:05")) + msg)

	err := smtp.SendMail(
		s.host,
		auth,
		s.username,
		s.recipientAddresses,
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
