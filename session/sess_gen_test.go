package session

import (
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"strconv"
	"testing"
)

func testBefor() *JwtTokenGenerator {
	return &JwtTokenGenerator{
		config: &ManagerConfig{
			CookieLifeTime: 3600,
			JwtSecKey:      "abc",
		},
		hmac: jwt.SigningMethodHS256,
	}
}

func TestJwtTokenGenerator_NewRawSessionID(t *testing.T) {
	gen := testBefor()
	rid, e := gen.UpgradeRawSessionID("sidabc")
	if nil != e {
		t.Error(e.Error())
	} else {
		t.Log(rid, "test pass")
	}
}

func TestJwtTokenGenerator_Valid(t *testing.T) {
	gen := testBefor()
	id := rand.Int()
	rid, e := gen.UpgradeRawSessionID(strconv.Itoa(id))
	if nil != e {
		t.Error(e.Error())
	}
	sid, e := gen.GetSessionID(rid)
	if nil != e {
		t.Error(e.Error())
	} else {
		t.Log("valid")
	}
	oid, e := strconv.Atoi(sid)
	if nil != e {
		t.Error(e.Error())
	}
	if oid == id {
		t.Log("test pass")
	}
}

func BenchmarkJwtTokenGenerator_GetSessionID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gen := testBefor()
		id := rand.Int()
		rid, e := gen.UpgradeRawSessionID(strconv.Itoa(id))
		if nil != e {
			b.Error(e.Error())
		}
		sid, e := gen.GetSessionID(rid)
		if nil != e {
			b.Error(e.Error())
		} else {
			//b.Log("valid")
		}
		oid, e := strconv.Atoi(sid)
		if nil != e {
			b.Error(e.Error())
		}
		if oid == id {
			//b.Log("test pass")
		}
	}
}

func BenchmarkJwtTokenGenerator_NewRawSessionID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gen := testBefor()

		id := make([]byte, 32)
		rand.Read(id)
		_, e := gen.UpgradeRawSessionID(string(id))
		if nil != e {
			b.Error(e.Error())
		}
	}
}
