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

package adapter

import (
	"html/template"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubstr(t *testing.T) {
	s := `012345`
	assert.Equal(t, "01", Substr(s, 0, 2))
	assert.Equal(t, "012345", Substr(s, 0, 100))
	assert.Equal(t, "012345", Substr(s, 12, 100))
}

func TestHtml2str(t *testing.T) {
	h := `<HTML><style></style><script>x<x</script></HTML><123>  123\n


	\n`
	assert.Equal(t, "123\\n\n\\n", HTML2str(h))
}

func TestDateFormat(t *testing.T) {
	ts := "Mon, 01 Jul 2013 13:27:42 CST"
	tt, _ := time.Parse(time.RFC1123, ts)

	assert.Equal(t, "2013-07-01 13:27:42", DateFormat(tt, "2006-01-02 15:04:05"))
}

func TestDate(t *testing.T) {
	ts := "Mon, 01 Jul 2013 13:27:42 CST"
	tt, _ := time.Parse(time.RFC1123, ts)

	assert.Equal(t, "2013-07-01 13:27:42", Date(tt, "Y-m-d H:i:s"))

	assert.Equal(t, "13-7-1 01:27:42 PM", Date(tt, "y-n-j h:i:s A"))
	assert.Equal(t, "Mon, 01 Jul 2013 1:27:42 pm", Date(tt, "D, d M Y g:i:s a"))
	assert.Equal(t, "Monday, 01 July 2013 13:27:42", Date(tt, "l, d F Y G:i:s"))
}

func TestCompareRelated(t *testing.T) {
	assert.True(t, Compare("abc", "abc"))

	assert.False(t, Compare("abc", "aBc"))

	assert.True(t, Compare("1", 1))

	assert.False(t, CompareNot("abc", "abc"))

	assert.True(t, CompareNot("abc", "aBc"))
	assert.True(t, NotNil("a string"))
}

func TestHtmlquote(t *testing.T) {
	h := `&lt;&#39;&nbsp;&rdquo;&ldquo;&amp;&#34;&gt;`
	s := `<' ”“&">`
	assert.Equal(t, h, Htmlquote(s))
}

func TestHtmlunquote(t *testing.T) {
	h := `&lt;&#39;&nbsp;&rdquo;&ldquo;&amp;&#34;&gt;`
	s := `<' ”“&">`
	assert.Equal(t, s, Htmlunquote(h))

}

func TestParseForm(t *testing.T) {
	type ExtendInfo struct {
		Hobby []string `form:"hobby"`
		Memo  string
	}

	type OtherInfo struct {
		Organization string `form:"organization"`
		Title        string `form:"title"`
		ExtendInfo
	}

	type user struct {
		ID      int         `form:"-"`
		tag     string      `form:"tag"`
		Name    interface{} `form:"username"`
		Age     int         `form:"age,text"`
		Email   string
		Intro   string    `form:",textarea"`
		StrBool bool      `form:"strbool"`
		Date    time.Time `form:"date,2006-01-02"`
		OtherInfo
	}

	u := user{}
	form := url.Values{
		"ID":           []string{"1"},
		"-":            []string{"1"},
		"tag":          []string{"no"},
		"username":     []string{"test"},
		"age":          []string{"40"},
		"Email":        []string{"test@gmail.com"},
		"Intro":        []string{"I am an engineer!"},
		"strbool":      []string{"yes"},
		"date":         []string{"2014-11-12"},
		"organization": []string{"beego"},
		"title":        []string{"CXO"},
		"hobby":        []string{"", "Basketball", "Football"},
		"memo":         []string{"nothing"},
	}

	assert.NotNil(t, ParseForm(form, u))

	assert.Nil(t, ParseForm(form, &u))

	assert.Equal(t, 0, u.ID)

	assert.Equal(t, 0, len(u.tag))

	assert.Equal(t, "test", u.Name)

	assert.Equal(t, 40, u.Age)

	assert.Equal(t, "test@gmail.com", u.Email)

	assert.Equal(t, "I am an engineer!", u.Intro)

	assert.True(t, u.StrBool)

	y, m, d := u.Date.Date()

	assert.Equal(t, 2014, y)
	assert.Equal(t, "November", m.String())
	assert.Equal(t, 12, d)

	assert.Equal(t, "beego", u.Organization)

	assert.Equal(t, "CXO", u.Title)

	assert.Equal(t, "", u.Hobby[0])

	assert.Equal(t, "Basketball", u.Hobby[1])

	assert.Equal(t, "Football", u.Hobby[2])

	assert.Equal(t, 0, len(u.Memo))
}

func TestRenderForm(t *testing.T) {
	type user struct {
		ID      int         `form:"-"`
		Name    interface{} `form:"username"`
		Age     int         `form:"age,text,年龄："`
		Sex     string
		Email   []string
		Intro   string `form:",textarea"`
		Ignored string `form:"-"`
	}

	u := user{Name: "test", Intro: "Some Text"}
	output := RenderForm(u)
	assert.Equal(t, template.HTML(""), output)
	output = RenderForm(&u)
	result := template.HTML(
		`Name: <input name="username" type="text" value="test"></br>` +
			`年龄：<input name="age" type="text" value="0"></br>` +
			`Sex: <input name="Sex" type="text" value=""></br>` +
			`Intro: <textarea name="Intro">Some Text</textarea>`)
	assert.Equal(t, result, output)
}

func TestMapGet(t *testing.T) {
	// test one level map
	m1 := map[string]int64{
		"a": 1,
		"1": 2,
	}

	res, err := MapGet(m1, "a")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), res)

	res, err = MapGet(m1, "1")
	assert.Nil(t, err)
	assert.Equal(t, int64(2), res)


	res, err = MapGet(m1, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), res)

	// test 2 level map
	m2 := M{
		"1": map[string]float64{
			"2": 3.5,
		},
	}

	res, err = MapGet(m2, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, 3.5, res)

	// test 5 level map
	m5 := M{
		"1": M{
			"2": M{
				"3": M{
					"4": M{
						"5": 1.2,
					},
				},
			},
		},
	}

	res, err = MapGet(m5, 1, 2, 3, 4, 5)
	assert.Nil(t, err)
	assert.Equal(t, 1.2, res)

	// check whether element not exists in map
	res, err = MapGet(m5, 5, 4, 3, 2, 1)
	assert.Nil(t, err)
	assert.Nil(t, res)

}
