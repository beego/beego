package beego

import (
	"html/template"
	"net/url"
	"testing"
	"time"
)

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

func TestParseForm(t *testing.T) {
	type user struct {
		Id    int         `form:"-"`
		tag   string      `form:"tag"`
		Name  interface{} `form:"username"`
		Age   int         `form:"age,text"`
		Email string
		Intro string `form:",textarea"`
	}

	u := user{}
	form := url.Values{
		"Id":       []string{"1"},
		"-":        []string{"1"},
		"tag":      []string{"no"},
		"username": []string{"test"},
		"age":      []string{"40"},
		"Email":    []string{"test@gmail.com"},
		"Intro":    []string{"I am an engineer!"},
	}
	if err := ParseForm(form, u); err == nil {
		t.Fatal("nothing will be changed")
	}
	if err := ParseForm(form, &u); err != nil {
		t.Fatal(err)
	}
	if u.Id != 0 {
		t.Errorf("Id should equal 0 but got %v", u.Id)
	}
	if len(u.tag) != 0 {
		t.Errorf("tag's length should equal 0 but got %v", len(u.tag))
	}
	if u.Name.(string) != "test" {
		t.Errorf("Name should equal `test` but got `%v`", u.Name.(string))
	}
	if u.Age != 40 {
		t.Errorf("Age should equal 40 but got %v", u.Age)
	}
	if u.Email != "test@gmail.com" {
		t.Errorf("Email should equal `test@gmail.com` but got `%v`", u.Email)
	}
	if u.Intro != "I am an engineer!" {
		t.Errorf("Intro should equal `I am an engineer!` but got `%v`", u.Intro)
	}
}

func TestRenderForm(t *testing.T) {
	type user struct {
		Id    int         `form:"-"`
		tag   string      `form:"tag"`
		Name  interface{} `form:"username"`
		Age   int         `form:"age,text,年龄："`
		Sex   string
		Email []string
		Intro string `form:",textarea"`
	}

	u := user{Name: "test"}
	output := RenderForm(u)
	if output != template.HTML("") {
		t.Errorf("output should be empty but got %v", output)
	}
	output = RenderForm(&u)
	result := template.HTML(
		`Name: <input name="username" type="text" value="test"></br>` +
			`年龄：<input name="age" type="text" value="0"></br>` +
			`Sex: <input name="Sex" type="text" value=""></br>` +
			`Intro: <input name="Intro" type="textarea" value="">`)
	if output != result {
		t.Errorf("output should equal `%v` but got `%v`", result, output)
	}
}
