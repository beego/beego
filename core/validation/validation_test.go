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

package validation

import (
	"regexp"
	"testing"
	"time"
)

func TestRequired(t *testing.T) {
	valid := Validation{}

	if valid.Required(nil, "nil").Ok {
		t.Error("nil object should be false")
	}
	if !valid.Required(true, "bool").Ok {
		t.Error("Bool value should always return true")
	}
	if !valid.Required(false, "bool").Ok {
		t.Error("Bool value should always return true")
	}
	if valid.Required("", "string").Ok {
		t.Error("\"'\" string should be false")
	}
	if valid.Required(" ", "string").Ok {
		t.Error("\" \" string should be false") // For #2361
	}
	if valid.Required("\n", "string").Ok {
		t.Error("new line string should be false") // For #2361
	}
	if !valid.Required("astaxie", "string").Ok {
		t.Error("string should be true")
	}
	if valid.Required(0, "zero").Ok {
		t.Error("Integer should not be equal 0")
	}
	if !valid.Required(1, "int").Ok {
		t.Error("Integer except 0 should be true")
	}
	if !valid.Required(time.Now(), "time").Ok {
		t.Error("time should be true")
	}
	if valid.Required([]string{}, "emptySlice").Ok {
		t.Error("empty slice should be false")
	}
	if !valid.Required([]interface{}{"ok"}, "slice").Ok {
		t.Error("slice should be true")
	}
}

func TestMin(t *testing.T) {
	valid := Validation{}

	if valid.Min(-1, 0, "min0").Ok {
		t.Error("-1 is less than the minimum value of 0 should be false")
	}
	if !valid.Min(1, 0, "min0").Ok {
		t.Error("1 is greater or equal than the minimum value of 0 should be true")
	}
}

func TestMax(t *testing.T) {
	valid := Validation{}

	if valid.Max(1, 0, "max0").Ok {
		t.Error("1 is greater than the minimum value of 0 should be false")
	}
	if !valid.Max(-1, 0, "max0").Ok {
		t.Error("-1 is less or equal than the maximum value of 0 should be true")
	}
}

func TestRange(t *testing.T) {
	valid := Validation{}

	if valid.Range(-1, 0, 1, "range0_1").Ok {
		t.Error("-1 is between 0 and 1 should be false")
	}
	if !valid.Range(1, 0, 1, "range0_1").Ok {
		t.Error("1 is between 0 and 1 should be true")
	}
}

func TestMinSize(t *testing.T) {
	valid := Validation{}

	if valid.MinSize("", 1, "minSize1").Ok {
		t.Error("the length of \"\" is less than the minimum value of 1 should be false")
	}
	if !valid.MinSize("ok", 1, "minSize1").Ok {
		t.Error("the length of \"ok\" is greater or equal than the minimum value of 1 should be true")
	}
	if valid.MinSize([]string{}, 1, "minSize1").Ok {
		t.Error("the length of empty slice is less than the minimum value of 1 should be false")
	}
	if !valid.MinSize([]interface{}{"ok"}, 1, "minSize1").Ok {
		t.Error("the length of [\"ok\"] is greater or equal than the minimum value of 1 should be true")
	}
}

func TestMaxSize(t *testing.T) {
	valid := Validation{}

	if valid.MaxSize("ok", 1, "maxSize1").Ok {
		t.Error("the length of \"ok\" is greater than the maximum value of 1 should be false")
	}
	if !valid.MaxSize("", 1, "maxSize1").Ok {
		t.Error("the length of \"\" is less or equal than the maximum value of 1 should be true")
	}
	if valid.MaxSize([]interface{}{"ok", false}, 1, "maxSize1").Ok {
		t.Error("the length of [\"ok\", false] is greater than the maximum value of 1 should be false")
	}
	if !valid.MaxSize([]string{}, 1, "maxSize1").Ok {
		t.Error("the length of empty slice is less or equal than the maximum value of 1 should be true")
	}
}

func TestLength(t *testing.T) {
	valid := Validation{}

	if valid.Length("", 1, "length1").Ok {
		t.Error("the length of \"\" must equal 1 should be false")
	}
	if !valid.Length("1", 1, "length1").Ok {
		t.Error("the length of \"1\" must equal 1 should be true")
	}
	if valid.Length([]string{}, 1, "length1").Ok {
		t.Error("the length of empty slice must equal 1 should be false")
	}
	if !valid.Length([]interface{}{"ok"}, 1, "length1").Ok {
		t.Error("the length of [\"ok\"] must equal 1 should be true")
	}
}

func TestAlpha(t *testing.T) {
	valid := Validation{}

	if valid.Alpha("a,1-@ $", "alpha").Ok {
		t.Error("\"a,1-@ $\" are valid alpha characters should be false")
	}
	if !valid.Alpha("abCD", "alpha").Ok {
		t.Error("\"abCD\" are valid alpha characters should be true")
	}
}

func TestNumeric(t *testing.T) {
	valid := Validation{}

	if valid.Numeric("a,1-@ $", "numeric").Ok {
		t.Error("\"a,1-@ $\" are valid numeric characters should be false")
	}
	if !valid.Numeric("1234", "numeric").Ok {
		t.Error("\"1234\" are valid numeric characters should be true")
	}
}

func TestAlphaNumeric(t *testing.T) {
	valid := Validation{}

	if valid.AlphaNumeric("a,1-@ $", "alphaNumeric").Ok {
		t.Error("\"a,1-@ $\" are valid alpha or numeric characters should be false")
	}
	if !valid.AlphaNumeric("1234aB", "alphaNumeric").Ok {
		t.Error("\"1234aB\" are valid alpha or numeric characters should be true")
	}
}

func TestMatch(t *testing.T) {
	valid := Validation{}

	if valid.Match("suchuangji@gmail", regexp.MustCompile(`^\w+@\w+\.\w+$`), "match").Ok {
		t.Error("\"suchuangji@gmail\" match \"^\\w+@\\w+\\.\\w+$\"  should be false")
	}
	if !valid.Match("suchuangji@gmail.com", regexp.MustCompile(`^\w+@\w+\.\w+$`), "match").Ok {
		t.Error("\"suchuangji@gmail\" match \"^\\w+@\\w+\\.\\w+$\"  should be true")
	}
}

func TestNoMatch(t *testing.T) {
	valid := Validation{}

	if valid.NoMatch("123@gmail", regexp.MustCompile(`[^\w\d]`), "nomatch").Ok {
		t.Error("\"123@gmail\" not match \"[^\\w\\d]\"  should be false")
	}
	if !valid.NoMatch("123gmail", regexp.MustCompile(`[^\w\d]`), "match").Ok {
		t.Error("\"123@gmail\" not match \"[^\\w\\d@]\"  should be true")
	}
}

func TestAlphaDash(t *testing.T) {
	valid := Validation{}

	if valid.AlphaDash("a,1-@ $", "alphaDash").Ok {
		t.Error("\"a,1-@ $\" are valid alpha or numeric or dash(-_) characters should be false")
	}
	if !valid.AlphaDash("1234aB-_", "alphaDash").Ok {
		t.Error("\"1234aB\" are valid alpha or numeric or dash(-_) characters should be true")
	}
}

func TestEmail(t *testing.T) {
	valid := Validation{}

	if valid.Email("not@a email", "email").Ok {
		t.Error("\"not@a email\" is a valid email address should be false")
	}
	if !valid.Email("suchuangji@gmail.com", "email").Ok {
		t.Error("\"suchuangji@gmail.com\" is a valid email address should be true")
	}
	if valid.Email("@suchuangji@gmail.com", "email").Ok {
		t.Error("\"@suchuangji@gmail.com\" is a valid email address should be false")
	}
	if valid.Email("suchuangji@gmail.com ok", "email").Ok {
		t.Error("\"suchuangji@gmail.com ok\" is a valid email address should be false")
	}
}

func TestIP(t *testing.T) {
	valid := Validation{}

	if valid.IP("11.255.255.256", "IP").Ok {
		t.Error("\"11.255.255.256\" is a valid ip address should be false")
	}
	if !valid.IP("01.11.11.11", "IP").Ok {
		t.Error("\"suchuangji@gmail.com\" is a valid ip address should be true")
	}
}

func TestBase64(t *testing.T) {
	valid := Validation{}

	if valid.Base64("suchuangji@gmail.com", "base64").Ok {
		t.Error("\"suchuangji@gmail.com\" are a valid base64 characters should be false")
	}
	if !valid.Base64("c3VjaHVhbmdqaUBnbWFpbC5jb20=", "base64").Ok {
		t.Error("\"c3VjaHVhbmdqaUBnbWFpbC5jb20=\" are a valid base64 characters should be true")
	}
}

func TestMobile(t *testing.T) {
	valid := Validation{}

	validMobiles := []string{
		"19800008888",
		"18800008888",
		"18000008888",
		"8618300008888",
		"+8614700008888",
		"17300008888",
		"+8617100008888",
		"8617500008888",
		"8617400008888",
		"16200008888",
		"16500008888",
		"16600008888",
		"16700008888",
		"13300008888",
		"14900008888",
		"15300008888",
		"17300008888",
		"17700008888",
		"18000008888",
		"18900008888",
		"19100008888",
		"19900008888",
		"19300008888",
		"13000008888",
		"13100008888",
		"13200008888",
		"14500008888",
		"15500008888",
		"15600008888",
		"16600008888",
		"17100008888",
		"17500008888",
		"17600008888",
		"18500008888",
		"18600008888",
		"13400008888",
		"13500008888",
		"13600008888",
		"13700008888",
		"13800008888",
		"13900008888",
		"14700008888",
		"15000008888",
		"15100008888",
		"15200008888",
		"15800008888",
		"15900008888",
		"17200008888",
		"17800008888",
		"18200008888",
		"18300008888",
		"18400008888",
		"18700008888",
		"18800008888",
		"19800008888",
	}

	for _, m := range validMobiles {
		if !valid.Mobile(m, "mobile").Ok {
			t.Error(m + " is a valid mobile phone number should be true")
		}
	}
}

func TestTel(t *testing.T) {
	valid := Validation{}

	if valid.Tel("222-00008888", "telephone").Ok {
		t.Error("\"222-00008888\" is a valid telephone number should be false")
	}
	if !valid.Tel("022-70008888", "telephone").Ok {
		t.Error("\"022-70008888\" is a valid telephone number should be true")
	}
	if !valid.Tel("02270008888", "telephone").Ok {
		t.Error("\"02270008888\" is a valid telephone number should be true")
	}
	if !valid.Tel("70008888", "telephone").Ok {
		t.Error("\"70008888\" is a valid telephone number should be true")
	}
}

func TestPhone(t *testing.T) {
	valid := Validation{}

	if valid.Phone("222-00008888", "phone").Ok {
		t.Error("\"222-00008888\" is a valid phone number should be false")
	}
	if !valid.Mobile("+8614700008888", "phone").Ok {
		t.Error("\"+8614700008888\" is a valid phone number should be true")
	}
	if !valid.Tel("02270008888", "phone").Ok {
		t.Error("\"02270008888\" is a valid phone number should be true")
	}
}

func TestZipCode(t *testing.T) {
	valid := Validation{}

	if valid.ZipCode("", "zipcode").Ok {
		t.Error("\"00008888\" is a valid zipcode should be false")
	}
	if !valid.ZipCode("536000", "zipcode").Ok {
		t.Error("\"536000\" is a valid zipcode should be true")
	}
}

func TestEnum(t *testing.T) {
	valid := Validation{}

	if valid.Enum("sms_code", "sms|email|code", "enum").Ok {
		t.Error("\"sms_code\" is in the enum list of \"sms|email|code\" should be false")
	}

	if !valid.Enum("sms", "email|sms|code", "enum").Ok {
		t.Error("\"sms\" is in the enum list of \"email|sms|code\" should be true")
	}

	if valid.Enum(200, "code|email|sms", "enum").Ok {
		t.Error("200 is in the enum list of \"code|email|sms\" should be false")
	}
}

func TestValid(t *testing.T) {
	type user struct {
		ID   int
		Name string `valid:"Required;Match(/^(test)?\\w*@(/test/);com$/)"`
		Age  int    `valid:"Required;Range(1, 140)"`
	}
	valid := Validation{}

	u := user{Name: "test@/test/;com", Age: 40}
	b, err := valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error("validation should be passed")
	}

	uptr := &user{Name: "test", Age: 40}
	valid.Clear()
	b, err = valid.Valid(uptr)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error("validation should not be passed")
	}
	if len(valid.Errors) != 1 {
		t.Fatalf("valid errors len should be 1 but got %d", len(valid.Errors))
	}
	if valid.Errors[0].Key != "Name.Match" {
		t.Errorf("Message key should be `Name.Match` but got %s", valid.Errors[0].Key)
	}

	u = user{Name: "test@/test/;com", Age: 180}
	valid.Clear()
	b, err = valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error("validation should not be passed")
	}
	if len(valid.Errors) != 1 {
		t.Fatalf("valid errors len should be 1 but got %d", len(valid.Errors))
	}
	if valid.Errors[0].Key != "Age.Range." {
		t.Errorf("Message key should be `Age.Range` but got %s", valid.Errors[0].Key)
	}
}

func TestRecursiveValid(t *testing.T) {
	type User struct {
		ID   int
		Name string `valid:"Required;Match(/^(test)?\\w*@(/test/);com$/)"`
		Age  int    `valid:"Required;Range(1, 140)"`
	}

	type AnonymouseUser struct {
		ID2   int
		Name2 string `valid:"Required;Match(/^(test)?\\w*@(/test/);com$/)"`
		Age2  int    `valid:"Required;Range(1, 140)"`
	}

	type Account struct {
		Password string `valid:"Required"`
		U        User
		AnonymouseUser
	}
	valid := Validation{}

	u := Account{Password: "abc123_", U: User{}}
	b, err := valid.RecursiveValid(u)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error("validation should not be passed")
	}
}

func TestSkipValid(t *testing.T) {
	type User struct {
		ID int

		Email    string `valid:"Email"`
		ReqEmail string `valid:"Required;Email"`

		IP    string `valid:"IP"`
		ReqIP string `valid:"Required;IP"`

		Mobile    string `valid:"Mobile"`
		ReqMobile string `valid:"Required;Mobile"`

		Tel    string `valid:"Tel"`
		ReqTel string `valid:"Required;Tel"`

		Phone    string `valid:"Phone"`
		ReqPhone string `valid:"Required;Phone"`

		ZipCode    string `valid:"ZipCode"`
		ReqZipCode string `valid:"Required;ZipCode"`
	}

	u := User{
		ReqEmail:   "a@a.com",
		ReqIP:      "127.0.0.1",
		ReqMobile:  "18888888888",
		ReqTel:     "02088888888",
		ReqPhone:   "02088888888",
		ReqZipCode: "510000",
	}

	valid := Validation{}
	b, err := valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Fatal("validation should not be passed")
	}

	valid = Validation{RequiredFirst: true}
	b, err = valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Fatal("validation should be passed")
	}
}

func TestPointer(t *testing.T) {
	type User struct {
		ID int

		Email    *string `valid:"Email"`
		ReqEmail *string `valid:"Required;Email"`
	}

	u := User{
		ReqEmail: nil,
		Email:    nil,
	}

	valid := Validation{}
	b, err := valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Fatal("validation should not be passed")
	}

	validEmail := "a@a.com"
	u = User{
		ReqEmail: &validEmail,
		Email:    nil,
	}

	valid = Validation{RequiredFirst: true}
	b, err = valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Fatal("validation should be passed")
	}

	u = User{
		ReqEmail: &validEmail,
		Email:    nil,
	}

	valid = Validation{}
	b, err = valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Fatal("validation should not be passed")
	}

	invalidEmail := "a@a"
	u = User{
		ReqEmail: &validEmail,
		Email:    &invalidEmail,
	}

	valid = Validation{RequiredFirst: true}
	b, err = valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Fatal("validation should not be passed")
	}

	u = User{
		ReqEmail: &validEmail,
		Email:    &invalidEmail,
	}

	valid = Validation{}
	b, err = valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Fatal("validation should not be passed")
	}
}

func TestCanSkipAlso(t *testing.T) {
	type User struct {
		ID int

		Email      string `valid:"Email"`
		ReqEmail   string `valid:"Required;Email"`
		MatchRange int    `valid:"Range(10, 20)"`
	}

	u := User{
		ReqEmail:   "a@a.com",
		Email:      "",
		MatchRange: 0,
	}

	valid := Validation{RequiredFirst: true}
	b, err := valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Fatal("validation should not be passed")
	}

	valid = Validation{RequiredFirst: true}
	valid.CanSkipAlso("Range")
	b, err = valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Fatal("validation should be passed")
	}
}

func TestFieldNoEmpty(t *testing.T) {
	type User struct {
		Name string `json:"name" valid:"Match(/^[a-zA-Z][a-zA-Z0-9._-]{0,31}$/)"`
	}
	u := User{
		Name: "*",
	}

	valid := Validation{}
	b, err := valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Fatal("validation should be passed")
	}
	if len(valid.Errors) == 0 {
		t.Fatal("validation should be passed")
	}
	validErr := valid.Errors[0]
	if len(validErr.Field) == 0 {
		t.Fatal("validation should be passed")
	}
}

func TestAddErrorKeyParsing(t *testing.T) {
	two := &Validation{}
	two.AddError("Email.Custom", "invalid")
	if len(two.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(two.Errors))
	}
	e := two.Errors[0]
	if e.Field != "Email" {
		t.Errorf("two-component key: expected Field %q, got %q", "Email", e.Field)
	}
	if e.Name != "Custom" {
		t.Errorf("two-component key: expected Name %q, got %q", "Custom", e.Name)
	}
	if e.Label != "Email" {
		t.Errorf("two-component key: expected Label %q, got %q", "Email", e.Label)
	}
	if e.Message != "Email invalid" {
		t.Errorf("two-component key: expected Message %q, got %q", "Email invalid", e.Message)
	}
	if got := two.ErrorsMap["Email"]; len(got) != 1 || got[0] != e {
		t.Errorf("two-component key: expected error stored under ErrorsMap[\"Email\"]")
	}
	if _, ok := two.ErrorsMap[""]; ok {
		t.Errorf("two-component key: error should not be stored under ErrorsMap[\"\"]")
	}

	three := &Validation{}
	three.AddError("Email.Custom.EmailLabel", "invalid")
	e3 := three.Errors[0]
	if e3.Field != "Email" || e3.Name != "Custom" || e3.Label != "EmailLabel" {
		t.Errorf("three-component key: got Field=%q Name=%q Label=%q", e3.Field, e3.Name, e3.Label)
	}
	if e3.Message != "EmailLabel invalid" {
		t.Errorf("three-component key: expected Message %q, got %q", "EmailLabel invalid", e3.Message)
	}
	if got := three.ErrorsMap["Email"]; len(got) != 1 || got[0] != e3 {
		t.Errorf("three-component key: expected error stored under ErrorsMap[\"Email\"]")
	}

	threeEmptyLabel := &Validation{}
	threeEmptyLabel.AddError("Email.Custom.", "invalid")
	e4 := threeEmptyLabel.Errors[0]
	if e4.Label != "Email" {
		t.Errorf("three-component key with empty label: expected Label %q, got %q", "Email", e4.Label)
	}
	if e4.Message != "Email invalid" {
		t.Errorf("three-component key with empty label: expected Message %q, got %q", "Email invalid", e4.Message)
	}

	single := &Validation{}
	single.AddError("Email", "invalid")
	e1 := single.Errors[0]
	if e1.Field != "" {
		t.Errorf("single-component key: expected empty Field, got %q", e1.Field)
	}
	if e1.Message != "invalid" {
		t.Errorf("single-component key: expected Message %q without leading space, got %q", "invalid", e1.Message)
	}
}

func TestApplyMessageLabelPrefix(t *testing.T) {
	single := &Validation{}
	single.Required("", "Name")
	if len(single.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(single.Errors))
	}
	e := single.Errors[0]
	if e.Label != "" {
		t.Errorf("single-component key: expected empty Label, got %q", e.Label)
	}
	if e.Message != "Can not be empty" {
		t.Errorf("single-component key: expected Message %q without leading space, got %q", "Can not be empty", e.Message)
	}

	labeled := &Validation{}
	labeled.Required("", "Email.Custom.EmailLabel")
	e2 := labeled.Errors[0]
	if e2.Field != "Email" || e2.Name != "Custom" || e2.Label != "EmailLabel" {
		t.Errorf("three-component key: got Field=%q Name=%q Label=%q", e2.Field, e2.Name, e2.Label)
	}
	if e2.Message != "EmailLabel Can not be empty" {
		t.Errorf("three-component key: expected Message %q, got %q", "EmailLabel Can not be empty", e2.Message)
	}
}

func TestAddErrorEdgeCaseKeys(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		field   string
		errName string
		label   string
		message string
		mapKey  string
	}{
		{
			name:    "empty key",
			key:     "",
			field:   "",
			errName: "",
			label:   "",
			message: "invalid",
			mapKey:  "",
		},
		{
			name:    "leading dot",
			key:     ".Custom",
			field:   "",
			errName: "Custom",
			label:   "",
			message: "invalid",
			mapKey:  "",
		},
		{
			name:    "trailing dot",
			key:     "Email.",
			field:   "Email",
			errName: "",
			label:   "Email",
			message: "Email invalid",
			mapKey:  "Email",
		},
		{
			name:    "empty field with label",
			key:     ".Custom.EmailLabel",
			field:   "",
			errName: "Custom",
			label:   "EmailLabel",
			message: "EmailLabel invalid",
			mapKey:  "",
		},
		{
			name:    "more than three components",
			key:     "A.B.C.D",
			field:   "",
			errName: "A.B.C.D",
			label:   "",
			message: "invalid",
			mapKey:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validation{}
			v.AddError(tt.key, "invalid")
			if len(v.Errors) != 1 {
				t.Fatalf("expected 1 error, got %d", len(v.Errors))
			}
			e := v.Errors[0]
			if e.Key != tt.key {
				t.Errorf("expected Key %q, got %q", tt.key, e.Key)
			}
			if e.Field != tt.field {
				t.Errorf("expected Field %q, got %q", tt.field, e.Field)
			}
			if e.Name != tt.errName {
				t.Errorf("expected Name %q, got %q", tt.errName, e.Name)
			}
			if e.Label != tt.label {
				t.Errorf("expected Label %q, got %q", tt.label, e.Label)
			}
			if e.Message != tt.message {
				t.Errorf("expected Message %q, got %q", tt.message, e.Message)
			}
			if got := v.ErrorsMap[tt.mapKey]; len(got) != 1 || got[0] != e {
				t.Errorf("expected error stored under ErrorsMap[%q]", tt.mapKey)
			}
		})
	}
}

func TestApplyEdgeCaseKeys(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		field   string
		errName string
		label   string
		message string
	}{
		{
			name:    "empty key",
			key:     "",
			field:   "",
			errName: "",
			label:   "",
			message: "Can not be empty",
		},
		{
			name:    "leading dot",
			key:     ".Custom",
			field:   "",
			errName: "Custom",
			label:   "",
			message: "Can not be empty",
		},
		{
			name:    "trailing dot",
			key:     "Email.",
			field:   "Email",
			errName: "",
			label:   "Email",
			message: "Email Can not be empty",
		},
		{
			name:    "more than three components",
			key:     "A.B.C.D",
			field:   "",
			errName: "A.B.C.D",
			label:   "",
			message: "Can not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validation{}
			v.Required("", tt.key)
			if len(v.Errors) != 1 {
				t.Fatalf("expected 1 error, got %d", len(v.Errors))
			}
			e := v.Errors[0]
			if e.Key != tt.key {
				t.Errorf("expected Key %q, got %q", tt.key, e.Key)
			}
			if e.Field != tt.field {
				t.Errorf("expected Field %q, got %q", tt.field, e.Field)
			}
			if e.Name != tt.errName {
				t.Errorf("expected Name %q, got %q", tt.errName, e.Name)
			}
			if e.Label != tt.label {
				t.Errorf("expected Label %q, got %q", tt.label, e.Label)
			}
			if e.Message != tt.message {
				t.Errorf("expected Message %q, got %q", tt.message, e.Message)
			}
		})
	}
}

func TestErrorsMapGroupsByField(t *testing.T) {
	v := &Validation{}
	v.AddError("Email.Custom", "invalid")
	v.AddError("Email.Required", "required")
	v.AddError("Name.Custom", "invalid")

	if len(v.Errors) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(v.Errors))
	}
	if got := v.ErrorsMap["Email"]; len(got) != 2 {
		t.Errorf("expected 2 errors under ErrorsMap[\"Email\"], got %d", len(got))
	}
	if got := v.ErrorsMap["Name"]; len(got) != 1 {
		t.Errorf("expected 1 error under ErrorsMap[\"Name\"], got %d", len(got))
	}
	if _, ok := v.ErrorsMap[""]; ok {
		t.Errorf("no error should be stored under ErrorsMap[\"\"]")
	}
}
