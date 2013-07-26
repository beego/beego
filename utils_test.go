package beego

import (
	"net/url"
	"testing"
	"time"
)

func TestWebTime(t *testing.T) {
	ts := "Fri, 26 Jul 2013 12:27:42 CST"
	l, _ := time.LoadLocation("GST")
	tt, _ := time.ParseInLocation(time.RFC1123, ts, l)
	if ts != webTime(tt) {
		t.Error("should be equal")
	}
	if "Fri, 26 Jul 2013 12:27:42 GMT" != webTime(tt.UTC()) {
		t.Error("should be equal")
	}
}

func TestSubstr(t *testing.T) {
	s := `012345`
	if Substr(s, 0, 2) != "01" {
		t.Error("should be equal")
	}
	if Substr(s, 0, 100) != "012345" {
		t.Error("should be equal")
	}
}

func TestHtml2str(t *testing.T) {
	h := `<HTML><style></style><script>x<x</script></HTML><123>  123\n


	\n`
	if Html2str(h) != "123\\n\n\\n" {
		t.Error("should be equal")
	}
}

func TestDateFormat(t *testing.T) {
	ts := "Mon, 01 Jul 2013 13:27:42 CST"
	tt, _ := time.Parse(time.RFC1123, ts)
	if DateFormat(tt, "2006-01-02 15:04:05") != "2013-07-01 13:27:42" {
		t.Error("should be equal")
	}
}

func TestDate(t *testing.T) {
	ts := "Mon, 01 Jul 2013 13:27:42 CST"
	tt, _ := time.Parse(time.RFC1123, ts)
	if Date(tt, "Y-m-d H:i:s") != "2013-07-01 13:27:42" {
		t.Error("should be equal")
	}
	if Date(tt, "y-n-j h:i:s A") != "13-7-1 01:27:42 PM" {
		t.Error("should be equal")
	}
	if Date(tt, "D, d M Y g:i:s a") != "Mon, 01 Jul 2013 1:27:42 pm" {
		t.Error("should be equal")
	}
	if Date(tt, "l, d F Y G:i:s") != "Monday, 01 July 2013 13:27:42" {
		t.Error("should be equal")
	}
}

func TestCompare(t *testing.T) {
	if !Compare("abc", "abc") {
		t.Error("should be equal")
	}
	if Compare("abc", "aBc") {
		t.Error("should be not equal")
	}
	if !Compare("1", 1) {
		t.Error("should be equal")
	}
}

func TestHtmlquote(t *testing.T) {
	h := `&lt;&#39;&nbsp;&rdquo;&ldquo;&amp;&quot;&gt;`
	s := `<' ”“&">`
	if Htmlquote(s) != h {
		t.Error("should be equal")
	}
}

func TestHtmlunquote(t *testing.T) {
	h := `&lt;&#39;&nbsp;&rdquo;&ldquo;&amp;&quot;&gt;`
	s := `<' ”“&">`
	if Htmlunquote(h) != s {
		t.Error("should be equal")
	}
}

func TestInSlice(t *testing.T) {
	sl := []string{"A", "b"}
	if !inSlice("A", sl) {
		t.Error("should be true")
	}
	if inSlice("B", sl) {
		t.Error("should be false")
	}
}

func TestParseForm(t *testing.T) {
	type user struct {
		Id    int
		tag   string      `form:tag`
		Name  interface{} `form:"username"`
		Age   int         `form:"age"`
		Email string
	}

	u := user{}
	form := url.Values{
		"tag":      []string{"no"},
		"username": []string{"test"},
		"age":      []string{"40"},
		"Email":    []string{"test@gmail.com"},
	}
	if err := ParseForm(form, u); err == nil {
		t.Fatal("nothing will be changed")
	}
	if err := ParseForm(form, &u); err != nil {
		t.Fatal(err)
	}
	if u.Id != 0 {
		t.Error("Id should not be changed")
	}
	if len(u.tag) != 0 {
		t.Error("tag should not be changed")
	}
	if u.Name.(string) != "test" {
		t.Error("should be equal")
	}
	if u.Age != 40 {
		t.Error("should be equal")
	}
	if u.Email != "test@gmail.com" {
		t.Error("should be equal")
	}
}
