// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logs

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

const (
	subjectPhrase = "Diagnostic message from server"
)

// smtpWriter implements LoggerInterface and is used to send emails via given SMTP-server.
type SmtpWriter struct {
	Username           string   `json:"Username"`
	Password           string   `json:"password"`
	Host               string   `json:"Host"`
	Subject            string   `json:"subject"`
	FromAddress        string   `json:"fromAddress"`
	RecipientAddresses []string `json:"sendTos"`
	Level              int      `json:"level"`
}

// create smtp writer.
func NewSmtpWriter() LoggerInterface {
	return &SmtpWriter{Level: LevelTrace}
}

// init smtp writer with json config.
// config like:
//	{
//		"Username":"example@gmail.com",
//		"password:"password",
//		"host":"smtp.gmail.com:465",
//		"subject":"email title",
//		"fromAddress":"from@example.com",
//		"sendTos":["email1","email2"],
//		"level":LevelError
//	}
func (s *SmtpWriter) Init(jsonconfig string) error {
	err := json.Unmarshal([]byte(jsonconfig), s)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmtpWriter) GetSmtpAuth(host string) smtp.Auth {
	if len(strings.Trim(s.Username, " ")) == 0 && len(strings.Trim(s.Password, " ")) == 0 {
		return nil
	}
	return smtp.PlainAuth(
		"",
		s.Username,
		s.Password,
		host,
	)
}

func (s *SmtpWriter) sendMail(hostAddressWithPort string, auth smtp.Auth, fromAddress string, recipients []string, msgContent []byte) error {
	client, err := smtp.Dial(hostAddressWithPort)
	if err != nil {
		return err
	}

	host, _, _ := net.SplitHostPort(hostAddressWithPort)
	tlsConn := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	if err = client.StartTLS(tlsConn); err != nil {
		return err
	}

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	if err = client.Mail(fromAddress); err != nil {
		return err
	}

	for _, rec := range recipients {
		if err = client.Rcpt(rec); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(msgContent))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = client.Quit()
	if err != nil {
		return err
	}

	return nil
}

// write message in smtp writer.
// it will send an email with subject and only this message.
func (s *SmtpWriter) WriteMsg(msg string, level int) error {
	if level > s.Level {
		return nil
	}

	hp := strings.Split(s.Host, ":")

	// Set up authentication information.
	auth := s.GetSmtpAuth(hp[0])

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	content_type := "Content-Type: text/plain" + "; charset=UTF-8"
	mailmsg := []byte("To: " + strings.Join(s.RecipientAddresses, ";") + "\r\nFrom: " + s.FromAddress + "<" + s.FromAddress +
		">\r\nSubject: " + s.Subject + "\r\n" + content_type + "\r\n\r\n" + fmt.Sprintf(".%s", time.Now().Format("2006-01-02 15:04:05")) + msg)

	err := s.sendMail(s.Host, auth, s.FromAddress, s.RecipientAddresses, mailmsg)

	return err
}

// implementing method. empty.
func (s *SmtpWriter) Flush() {
	return
}

// implementing method. empty.
func (s *SmtpWriter) Destroy() {
	return
}

func init() {
	Register("smtp", NewSmtpWriter)
}
