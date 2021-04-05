package session

import (
	"net/http"
	"testing"
)

func TestCfgCookieLifeTime(t *testing.T) {
	value := 8754
	c := NewManagerConfig(
		CfgCookieLifeTime(value),
	)

	if c.CookieLifeTime != value {
		t.Error()
	}
}

func TestCfgDomain(t *testing.T) {
	value := `http://domain.com`
	c := NewManagerConfig(
		CfgDomain(value),
	)

	if c.Domain != value {
		t.Error()
	}
}

func TestCfgSameSite(t *testing.T) {
	value := http.SameSiteLaxMode
	c := NewManagerConfig(
		CfgSameSite(value),
	)

	if c.CookieSameSite != value {
		t.Error()
	}
}

func TestCfgSecure(t *testing.T) {
	c := NewManagerConfig(
		CfgSecure(true),
	)

	if c.Secure != true {
		t.Error()
	}
}

func TestCfgSecure1(t *testing.T) {
	c := NewManagerConfig(
		CfgSecure(false),
	)

	if c.Secure != false {
		t.Error()
	}
}

func TestCfgSessionIdPrefix(t *testing.T) {
	value := `sodiausodkljalsd`
	c := NewManagerConfig(
		CfgSessionIdPrefix(value),
	)

	if c.SessionIDPrefix != value {
		t.Error()
	}
}

func TestCfgSetSessionNameInHTTPHeader(t *testing.T) {
	value := `sodiausodkljalsd`
	c := NewManagerConfig(
		CfgSetSessionNameInHTTPHeader(value),
	)

	if c.SessionNameInHTTPHeader != value {
		t.Error()
	}
}

func TestCfgCookieName(t *testing.T) {
	value := `sodiausodkljalsd`
	c := NewManagerConfig(
		CfgCookieName(value),
	)

	if c.CookieName != value {
		t.Error()
	}
}

func TestCfgEnableSidInURLQuery(t *testing.T) {
	c := NewManagerConfig(
		CfgEnableSidInURLQuery(true),
	)

	if c.EnableSidInURLQuery != true {
		t.Error()
	}
}

func TestCfgGcLifeTime(t *testing.T) {
	value := int64(5454)
	c := NewManagerConfig(
		CfgGcLifeTime(value),
	)

	if c.Gclifetime != value {
		t.Error()
	}
}

func TestCfgHTTPOnly(t *testing.T) {
	c := NewManagerConfig(
		CfgHTTPOnly(true),
	)

	if c.DisableHTTPOnly != false {
		t.Error()
	}
}

func TestCfgHTTPOnly2(t *testing.T) {
	c := NewManagerConfig(
		CfgHTTPOnly(false),
	)

	if c.DisableHTTPOnly != true {
		t.Error()
	}
}

func TestCfgMaxLifeTime(t *testing.T) {
	value := int64(5454)
	c := NewManagerConfig(
		CfgMaxLifeTime(value),
	)

	if c.Maxlifetime != value {
		t.Error()
	}
}

func TestCfgProviderConfig(t *testing.T) {
	value := `asodiuasldkj12i39809as`
	c := NewManagerConfig(
		CfgProviderConfig(value),
	)

	if c.ProviderConfig != value {
		t.Error()
	}
}

func TestCfgSessionIdInHTTPHeader(t *testing.T) {
	c := NewManagerConfig(
		CfgSessionIdInHTTPHeader(true),
	)

	if c.EnableSidInHTTPHeader != true {
		t.Error()
	}
}

func TestCfgSessionIdInHTTPHeader1(t *testing.T) {
	c := NewManagerConfig(
		CfgSessionIdInHTTPHeader(false),
	)

	if c.EnableSidInHTTPHeader != false {
		t.Error()
	}
}

func TestCfgSessionIdLength(t *testing.T) {
	value := int64(100)
	c := NewManagerConfig(
		CfgSessionIdLength(value),
	)

	if c.SessionIDLength != value {
		t.Error()
	}
}

func TestCfgSetCookie(t *testing.T) {
	c := NewManagerConfig(
		CfgSetCookie(true),
	)

	if c.EnableSetCookie != true {
		t.Error()
	}
}

func TestCfgSetCookie1(t *testing.T) {
	c := NewManagerConfig(
		CfgSetCookie(false),
	)

	if c.EnableSetCookie != false {
		t.Error()
	}
}

func TestNewManagerConfig(t *testing.T) {
	c := NewManagerConfig()
	if c == nil {
		t.Error()
	}
}

func TestManagerConfig_Opts(t *testing.T) {
	c := NewManagerConfig()
	c.Opts(CfgSetCookie(true))

	if c.EnableSetCookie != true {
		t.Error()
	}
}
