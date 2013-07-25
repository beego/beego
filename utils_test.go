package beego

import (
	"net/url"
	"testing"
)

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
