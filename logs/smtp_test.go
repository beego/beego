package logs

import (
	"testing"
)

func TestSmtp(t *testing.T) {
	log := NewLogger(10000)
	log.SetLogger("smtp", `{"username":"beegotest@gmail.com","password":"xxxxxxxx","host":"smtp.gmail.com:587","sendTos":["xiemengjun@gmail.com"]}`)
	log.Critical("sendmail critical")
}
