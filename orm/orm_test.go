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

package orm

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

var _ = os.PathSeparator

var (
	test_Date     = format_Date + " -0700"
	test_DateTime = format_DateTime + " -0700"
)

func ValuesCompare(is bool, a interface{}, args ...interface{}) (err error, ok bool) {
	if len(args) == 0 {
		return fmt.Errorf("miss args"), false
	}
	b := args[0]
	arg := argAny(args)

	switch v := a.(type) {
	case reflect.Kind:
		ok = reflect.ValueOf(b).Kind() == v
	case time.Time:
		if v2, vo := b.(time.Time); vo {
			if arg.Get(1) != nil {
				format := ToStr(arg.Get(1))
				a = v.Format(format)
				b = v2.Format(format)
				ok = a == b
			} else {
				err = fmt.Errorf("compare datetime miss format")
				goto wrongArg
			}
		}
	default:
		ok = ToStr(a) == ToStr(b)
	}
	ok = is && ok || !is && !ok
	if !ok {
		if is {
			err = fmt.Errorf("expected: `%v`, get `%v`", b, a)
		} else {
			err = fmt.Errorf("expected: `%v`, get `%v`", b, a)
		}
	}

wrongArg:
	if err != nil {
		return err, false
	}

	return nil, true
}

func AssertIs(a interface{}, args ...interface{}) error {
	if err, ok := ValuesCompare(true, a, args...); ok == false {
		return err
	}
	return nil
}

func AssertNot(a interface{}, args ...interface{}) error {
	if err, ok := ValuesCompare(false, a, args...); ok == false {
		return err
	}
	return nil
}

func getCaller(skip int) string {
	pc, file, line, _ := runtime.Caller(skip)
	fun := runtime.FuncForPC(pc)
	_, fn := filepath.Split(file)
	data, err := ioutil.ReadFile(file)
	var codes []string
	if err == nil {
		lines := bytes.Split(data, []byte{'\n'})
		n := 10
		for i := 0; i < n; i++ {
			o := line - n
			if o < 0 {
				continue
			}
			cur := o + i + 1
			flag := "  "
			if cur == line {
				flag = ">>"
			}
			code := fmt.Sprintf(" %s %5d:   %s", flag, cur, strings.Replace(string(lines[o+i]), "\t", "    ", -1))
			if code != "" {
				codes = append(codes, code)
			}
		}
	}
	funName := fun.Name()
	if i := strings.LastIndex(funName, "."); i > -1 {
		funName = funName[i+1:]
	}
	return fmt.Sprintf("%s:%d: \n%s", fn, line, strings.Join(codes, "\n"))
}

func throwFail(t *testing.T, err error, args ...interface{}) {
	if err != nil {
		con := fmt.Sprintf("\t\nError: %s\n%s\n", err.Error(), getCaller(2))
		if len(args) > 0 {
			parts := make([]string, 0, len(args))
			for _, arg := range args {
				parts = append(parts, fmt.Sprintf("%v", arg))
			}
			con += " " + strings.Join(parts, ", ")
		}
		t.Error(con)
		t.Fail()
	}
}

func throwFailNow(t *testing.T, err error, args ...interface{}) {
	if err != nil {
		con := fmt.Sprintf("\t\nError: %s\n%s\n", err.Error(), getCaller(2))
		if len(args) > 0 {
			parts := make([]string, 0, len(args))
			for _, arg := range args {
				parts = append(parts, fmt.Sprintf("%v", arg))
			}
			con += " " + strings.Join(parts, ", ")
		}
		t.Error(con)
		t.FailNow()
	}
}

func TestGetDB(t *testing.T) {
	if db, err := GetDB(); err != nil {
		throwFailNow(t, err)
	} else {
		err = db.Ping()
		throwFailNow(t, err)
	}
}

func TestSyncDb(t *testing.T) {
	RegisterModel(new(Data), new(DataNull), new(DataCustom))
	RegisterModel(new(User))
	RegisterModel(new(Profile))
	RegisterModel(new(Post))
	RegisterModel(new(Tag))
	RegisterModel(new(Comment))
	RegisterModel(new(UserBig))
	RegisterModel(new(PostTags))

	err := RunSyncdb("default", true, false)
	throwFail(t, err)

	modelCache.clean()
}

func TestRegisterModels(t *testing.T) {
	RegisterModel(new(Data), new(DataNull), new(DataCustom))
	RegisterModel(new(User))
	RegisterModel(new(Profile))
	RegisterModel(new(Post))
	RegisterModel(new(Tag))
	RegisterModel(new(Comment))
	RegisterModel(new(UserBig))
	RegisterModel(new(PostTags))

	BootStrap()

	dORM = NewOrm()
	dDbBaser = getDbAlias("default").DbBaser
}

func TestModelSyntax(t *testing.T) {
	user := &User{}
	ind := reflect.ValueOf(user).Elem()
	fn := getFullName(ind.Type())
	mi, ok := modelCache.getByFN(fn)
	throwFail(t, AssertIs(ok, true))

	mi, ok = modelCache.get("user")
	throwFail(t, AssertIs(ok, true))
	if ok {
		throwFail(t, AssertIs(mi.fields.GetByName("ShouldSkip") == nil, true))
	}
}

var Data_Values = map[string]interface{}{
	"Boolean":  true,
	"Char":     "char",
	"Text":     "text",
	"Date":     time.Now(),
	"DateTime": time.Now(),
	"Byte":     byte(1<<8 - 1),
	"Rune":     rune(1<<31 - 1),
	"Int":      int(1<<31 - 1),
	"Int8":     int8(1<<7 - 1),
	"Int16":    int16(1<<15 - 1),
	"Int32":    int32(1<<31 - 1),
	"Int64":    int64(1<<63 - 1),
	"Uint":     uint(1<<32 - 1),
	"Uint8":    uint8(1<<8 - 1),
	"Uint16":   uint16(1<<16 - 1),
	"Uint32":   uint32(1<<32 - 1),
	"Uint64":   uint64(1<<63 - 1), // uint64 values with high bit set are not supported
	"Float32":  float32(100.1234),
	"Float64":  float64(100.1234),
	"Decimal":  float64(100.1234),
}

func TestDataTypes(t *testing.T) {
	d := Data{}
	ind := reflect.Indirect(reflect.ValueOf(&d))

	for name, value := range Data_Values {
		e := ind.FieldByName(name)
		e.Set(reflect.ValueOf(value))
	}

	id, err := dORM.Insert(&d)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 1))

	d = Data{Id: 1}
	err = dORM.Read(&d)
	throwFail(t, err)

	ind = reflect.Indirect(reflect.ValueOf(&d))

	for name, value := range Data_Values {
		e := ind.FieldByName(name)
		vu := e.Interface()
		switch name {
		case "Date":
			vu = vu.(time.Time).In(DefaultTimeLoc).Format(test_Date)
			value = value.(time.Time).In(DefaultTimeLoc).Format(test_Date)
		case "DateTime":
			vu = vu.(time.Time).In(DefaultTimeLoc).Format(test_DateTime)
			value = value.(time.Time).In(DefaultTimeLoc).Format(test_DateTime)
		}
		throwFail(t, AssertIs(vu == value, true), value, vu)
	}
}

func TestNullDataTypes(t *testing.T) {
	d := DataNull{}

	if IsPostgres {
		// can removed when this fixed
		// https://github.com/lib/pq/pull/125
		d.DateTime = time.Now()
	}

	id, err := dORM.Insert(&d)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 1))

	d = DataNull{Id: 1}
	err = dORM.Read(&d)
	throwFail(t, err)

	throwFail(t, AssertIs(d.NullBool.Valid, false))
	throwFail(t, AssertIs(d.NullString.Valid, false))
	throwFail(t, AssertIs(d.NullInt64.Valid, false))
	throwFail(t, AssertIs(d.NullFloat64.Valid, false))

	throwFail(t, AssertIs(d.BooleanPtr, nil))
	throwFail(t, AssertIs(d.CharPtr, nil))
	throwFail(t, AssertIs(d.TextPtr, nil))
	throwFail(t, AssertIs(d.BytePtr, nil))
	throwFail(t, AssertIs(d.RunePtr, nil))
	throwFail(t, AssertIs(d.IntPtr, nil))
	throwFail(t, AssertIs(d.Int8Ptr, nil))
	throwFail(t, AssertIs(d.Int16Ptr, nil))
	throwFail(t, AssertIs(d.Int32Ptr, nil))
	throwFail(t, AssertIs(d.Int64Ptr, nil))
	throwFail(t, AssertIs(d.UintPtr, nil))
	throwFail(t, AssertIs(d.Uint8Ptr, nil))
	throwFail(t, AssertIs(d.Uint16Ptr, nil))
	throwFail(t, AssertIs(d.Uint32Ptr, nil))
	throwFail(t, AssertIs(d.Uint64Ptr, nil))
	throwFail(t, AssertIs(d.Float32Ptr, nil))
	throwFail(t, AssertIs(d.Float64Ptr, nil))
	throwFail(t, AssertIs(d.DecimalPtr, nil))

	_, err = dORM.Raw(`INSERT INTO data_null (boolean) VALUES (?)`, nil).Exec()
	throwFail(t, err)

	d = DataNull{Id: 2}
	err = dORM.Read(&d)
	throwFail(t, err)

	booleanPtr := true
	charPtr := string("test")
	textPtr := string("test")
	bytePtr := byte('t')
	runePtr := rune('t')
	intPtr := int(42)
	int8Ptr := int8(42)
	int16Ptr := int16(42)
	int32Ptr := int32(42)
	int64Ptr := int64(42)
	uintPtr := uint(42)
	uint8Ptr := uint8(42)
	uint16Ptr := uint16(42)
	uint32Ptr := uint32(42)
	uint64Ptr := uint64(42)
	float32Ptr := float32(42.0)
	float64Ptr := float64(42.0)
	decimalPtr := float64(42.0)

	d = DataNull{
		DateTime:    time.Now(),
		NullString:  sql.NullString{String: "test", Valid: true},
		NullBool:    sql.NullBool{Bool: true, Valid: true},
		NullInt64:   sql.NullInt64{Int64: 42, Valid: true},
		NullFloat64: sql.NullFloat64{Float64: 42.42, Valid: true},
		BooleanPtr:  &booleanPtr,
		CharPtr:     &charPtr,
		TextPtr:     &textPtr,
		BytePtr:     &bytePtr,
		RunePtr:     &runePtr,
		IntPtr:      &intPtr,
		Int8Ptr:     &int8Ptr,
		Int16Ptr:    &int16Ptr,
		Int32Ptr:    &int32Ptr,
		Int64Ptr:    &int64Ptr,
		UintPtr:     &uintPtr,
		Uint8Ptr:    &uint8Ptr,
		Uint16Ptr:   &uint16Ptr,
		Uint32Ptr:   &uint32Ptr,
		Uint64Ptr:   &uint64Ptr,
		Float32Ptr:  &float32Ptr,
		Float64Ptr:  &float64Ptr,
		DecimalPtr:  &decimalPtr,
	}

	id, err = dORM.Insert(&d)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 3))

	d = DataNull{Id: 3}
	err = dORM.Read(&d)
	throwFail(t, err)

	throwFail(t, AssertIs(d.NullBool.Valid, true))
	throwFail(t, AssertIs(d.NullBool.Bool, true))

	throwFail(t, AssertIs(d.NullString.Valid, true))
	throwFail(t, AssertIs(d.NullString.String, "test"))

	throwFail(t, AssertIs(d.NullInt64.Valid, true))
	throwFail(t, AssertIs(d.NullInt64.Int64, 42))

	throwFail(t, AssertIs(d.NullFloat64.Valid, true))
	throwFail(t, AssertIs(d.NullFloat64.Float64, 42.42))

	throwFail(t, AssertIs(*d.BooleanPtr, booleanPtr))
	throwFail(t, AssertIs(*d.CharPtr, charPtr))
	throwFail(t, AssertIs(*d.TextPtr, textPtr))
	throwFail(t, AssertIs(*d.BytePtr, bytePtr))
	throwFail(t, AssertIs(*d.RunePtr, runePtr))
	throwFail(t, AssertIs(*d.IntPtr, intPtr))
	throwFail(t, AssertIs(*d.Int8Ptr, int8Ptr))
	throwFail(t, AssertIs(*d.Int16Ptr, int16Ptr))
	throwFail(t, AssertIs(*d.Int32Ptr, int32Ptr))
	throwFail(t, AssertIs(*d.Int64Ptr, int64Ptr))
	throwFail(t, AssertIs(*d.UintPtr, uintPtr))
	throwFail(t, AssertIs(*d.Uint8Ptr, uint8Ptr))
	throwFail(t, AssertIs(*d.Uint16Ptr, uint16Ptr))
	throwFail(t, AssertIs(*d.Uint32Ptr, uint32Ptr))
	throwFail(t, AssertIs(*d.Uint64Ptr, uint64Ptr))
	throwFail(t, AssertIs(*d.Float32Ptr, float32Ptr))
	throwFail(t, AssertIs(*d.Float64Ptr, float64Ptr))
	throwFail(t, AssertIs(*d.DecimalPtr, decimalPtr))
}

func TestDataCustomTypes(t *testing.T) {
	d := DataCustom{}
	ind := reflect.Indirect(reflect.ValueOf(&d))

	for name, value := range Data_Values {
		e := ind.FieldByName(name)
		if !e.IsValid() {
			continue
		}
		e.Set(reflect.ValueOf(value).Convert(e.Type()))
	}

	id, err := dORM.Insert(&d)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 1))

	d = DataCustom{Id: 1}
	err = dORM.Read(&d)
	throwFail(t, err)

	ind = reflect.Indirect(reflect.ValueOf(&d))

	for name, value := range Data_Values {
		e := ind.FieldByName(name)
		if !e.IsValid() {
			continue
		}
		vu := e.Interface()
		value = reflect.ValueOf(value).Convert(e.Type()).Interface()
		throwFail(t, AssertIs(vu == value, true), value, vu)
	}
}

func TestCRUD(t *testing.T) {
	profile := NewProfile()
	profile.Age = 30
	profile.Money = 1234.12
	id, err := dORM.Insert(profile)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 1))

	user := NewUser()
	user.UserName = "slene"
	user.Email = "vslene@gmail.com"
	user.Password = "pass"
	user.Status = 3
	user.IsStaff = true
	user.IsActive = true

	id, err = dORM.Insert(user)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 1))

	u := &User{Id: user.Id}
	err = dORM.Read(u)
	throwFail(t, err)

	throwFail(t, AssertIs(u.UserName, "slene"))
	throwFail(t, AssertIs(u.Email, "vslene@gmail.com"))
	throwFail(t, AssertIs(u.Password, "pass"))
	throwFail(t, AssertIs(u.Status, 3))
	throwFail(t, AssertIs(u.IsStaff, true))
	throwFail(t, AssertIs(u.IsActive, true))
	throwFail(t, AssertIs(u.Created.In(DefaultTimeLoc), user.Created.In(DefaultTimeLoc), test_Date))
	throwFail(t, AssertIs(u.Updated.In(DefaultTimeLoc), user.Updated.In(DefaultTimeLoc), test_DateTime))

	user.UserName = "astaxie"
	user.Profile = profile
	num, err := dORM.Update(user)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	u = &User{Id: user.Id}
	err = dORM.Read(u)
	throwFailNow(t, err)
	throwFail(t, AssertIs(u.UserName, "astaxie"))
	throwFail(t, AssertIs(u.Profile.Id, profile.Id))

	u = &User{UserName: "astaxie", Password: "pass"}
	err = dORM.Read(u, "UserName")
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(id, 1))

	u.UserName = "QQ"
	u.Password = "111"
	num, err = dORM.Update(u, "UserName")
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	u = &User{Id: user.Id}
	err = dORM.Read(u)
	throwFailNow(t, err)
	throwFail(t, AssertIs(u.UserName, "QQ"))
	throwFail(t, AssertIs(u.Password, "pass"))

	num, err = dORM.Delete(profile)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	u = &User{Id: user.Id}
	err = dORM.Read(u)
	throwFail(t, err)
	throwFail(t, AssertIs(true, u.Profile == nil))

	num, err = dORM.Delete(user)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	u = &User{Id: 100}
	err = dORM.Read(u)
	throwFail(t, AssertIs(err, ErrNoRows))

	ub := UserBig{}
	ub.Name = "name"
	id, err = dORM.Insert(&ub)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 1))

	ub = UserBig{Id: 1}
	err = dORM.Read(&ub)
	throwFail(t, err)
	throwFail(t, AssertIs(ub.Name, "name"))
}

func TestInsertTestData(t *testing.T) {
	var users []*User

	profile := NewProfile()
	profile.Age = 28
	profile.Money = 1234.12

	id, err := dORM.Insert(profile)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 2))

	user := NewUser()
	user.UserName = "slene"
	user.Email = "vslene@gmail.com"
	user.Password = "pass"
	user.Status = 1
	user.IsStaff = false
	user.IsActive = true
	user.Profile = profile

	users = append(users, user)

	id, err = dORM.Insert(user)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 2))

	profile = NewProfile()
	profile.Age = 30
	profile.Money = 4321.09

	id, err = dORM.Insert(profile)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 3))

	user = NewUser()
	user.UserName = "astaxie"
	user.Email = "astaxie@gmail.com"
	user.Password = "password"
	user.Status = 2
	user.IsStaff = true
	user.IsActive = false
	user.Profile = profile

	users = append(users, user)

	id, err = dORM.Insert(user)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 3))

	user = NewUser()
	user.UserName = "nobody"
	user.Email = "nobody@gmail.com"
	user.Password = "nobody"
	user.Status = 3
	user.IsStaff = false
	user.IsActive = false

	users = append(users, user)

	id, err = dORM.Insert(user)
	throwFail(t, err)
	throwFail(t, AssertIs(id, 4))

	tags := []*Tag{
		&Tag{Name: "golang", BestPost: &Post{Id: 2}},
		&Tag{Name: "example"},
		&Tag{Name: "format"},
		&Tag{Name: "c++"},
	}

	posts := []*Post{
		&Post{User: users[0], Tags: []*Tag{tags[0]}, Title: "Introduction", Content: `Go is a new language. Although it borrows ideas from existing languages, it has unusual properties that make effective Go programs different in character from programs written in its relatives. A straightforward translation of a C++ or Java program into Go is unlikely to produce a satisfactory result—Java programs are written in Java, not Go. On the other hand, thinking about the problem from a Go perspective could produce a successful but quite different program. In other words, to write Go well, it's important to understand its properties and idioms. It's also important to know the established conventions for programming in Go, such as naming, formatting, program construction, and so on, so that programs you write will be easy for other Go programmers to understand.
This document gives tips for writing clear, idiomatic Go code. It augments the language specification, the Tour of Go, and How to Write Go Code, all of which you should read first.`},
		&Post{User: users[1], Tags: []*Tag{tags[0], tags[1]}, Title: "Examples", Content: `The Go package sources are intended to serve not only as the core library but also as examples of how to use the language. Moreover, many of the packages contain working, self-contained executable examples you can run directly from the golang.org web site, such as this one (click on the word "Example" to open it up). If you have a question about how to approach a problem or how something might be implemented, the documentation, code and examples in the library can provide answers, ideas and background.`},
		&Post{User: users[1], Tags: []*Tag{tags[0], tags[2]}, Title: "Formatting", Content: `Formatting issues are the most contentious but the least consequential. People can adapt to different formatting styles but it's better if they don't have to, and less time is devoted to the topic if everyone adheres to the same style. The problem is how to approach this Utopia without a long prescriptive style guide.
With Go we take an unusual approach and let the machine take care of most formatting issues. The gofmt program (also available as go fmt, which operates at the package level rather than source file level) reads a Go program and emits the source in a standard style of indentation and vertical alignment, retaining and if necessary reformatting comments. If you want to know how to handle some new layout situation, run gofmt; if the answer doesn't seem right, rearrange your program (or file a bug about gofmt), don't work around it.`},
		&Post{User: users[2], Tags: []*Tag{tags[3]}, Title: "Commentary", Content: `Go provides C-style /* */ block comments and C++-style // line comments. Line comments are the norm; block comments appear mostly as package comments, but are useful within an expression or to disable large swaths of code.
The program—and web server—godoc processes Go source files to extract documentation about the contents of the package. Comments that appear before top-level declarations, with no intervening newlines, are extracted along with the declaration to serve as explanatory text for the item. The nature and style of these comments determines the quality of the documentation godoc produces.`},
	}

	comments := []*Comment{
		&Comment{Post: posts[0], Content: "a comment"},
		&Comment{Post: posts[1], Content: "yes"},
		&Comment{Post: posts[1]},
		&Comment{Post: posts[1]},
		&Comment{Post: posts[2]},
		&Comment{Post: posts[2]},
	}

	for _, tag := range tags {
		id, err := dORM.Insert(tag)
		throwFail(t, err)
		throwFail(t, AssertIs(id > 0, true))
	}

	for _, post := range posts {
		id, err := dORM.Insert(post)
		throwFail(t, err)
		throwFail(t, AssertIs(id > 0, true))

		num := len(post.Tags)
		if num > 0 {
			nums, err := dORM.QueryM2M(post, "tags").Add(post.Tags)
			throwFailNow(t, err)
			throwFailNow(t, AssertIs(nums, num))
		}
	}

	for _, comment := range comments {
		id, err := dORM.Insert(comment)
		throwFail(t, err)
		throwFail(t, AssertIs(id > 0, true))
	}
}

func TestCustomField(t *testing.T) {
	user := User{Id: 2}
	err := dORM.Read(&user)
	throwFailNow(t, err)

	user.Langs = append(user.Langs, "zh-CN", "en-US")
	user.Extra.Name = "beego"
	user.Extra.Data = "orm"
	_, err = dORM.Update(&user, "Langs", "Extra")
	throwFailNow(t, err)

	user = User{Id: 2}
	err = dORM.Read(&user)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(len(user.Langs), 2))
	throwFailNow(t, AssertIs(user.Langs[0], "zh-CN"))
	throwFailNow(t, AssertIs(user.Langs[1], "en-US"))

	throwFailNow(t, AssertIs(user.Extra.Name, "beego"))
	throwFailNow(t, AssertIs(user.Extra.Data, "orm"))
}

func TestExpr(t *testing.T) {
	user := &User{}
	qs := dORM.QueryTable(user)
	qs = dORM.QueryTable((*User)(nil))
	qs = dORM.QueryTable("User")
	qs = dORM.QueryTable("user")
	num, err := qs.Filter("UserName", "slene").Filter("user_name", "slene").Filter("profile__Age", 28).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("created", time.Now()).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 3))

	// num, err = qs.Filter("created", time.Now().Format(format_Date)).Count()
	// throwFail(t, err)
	// throwFail(t, AssertIs(num, 3))
}

func TestOperators(t *testing.T) {
	qs := dORM.QueryTable("user")
	num, err := qs.Filter("user_name", "slene").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("user_name__exact", String("slene")).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("user_name__exact", "slene").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("user_name__iexact", "Slene").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("user_name__contains", "e").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	var shouldNum int

	if IsSqlite {
		shouldNum = 2
	} else {
		shouldNum = 0
	}

	num, err = qs.Filter("user_name__contains", "E").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, shouldNum))

	num, err = qs.Filter("user_name__icontains", "E").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	num, err = qs.Filter("user_name__icontains", "E").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	num, err = qs.Filter("status__gt", 1).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	num, err = qs.Filter("status__gte", 1).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 3))

	num, err = qs.Filter("status__lt", Uint(3)).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	num, err = qs.Filter("status__lte", Int(3)).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 3))

	num, err = qs.Filter("user_name__startswith", "s").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	if IsSqlite {
		shouldNum = 1
	} else {
		shouldNum = 0
	}

	num, err = qs.Filter("user_name__startswith", "S").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, shouldNum))

	num, err = qs.Filter("user_name__istartswith", "S").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("user_name__endswith", "e").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	if IsSqlite {
		shouldNum = 2
	} else {
		shouldNum = 0
	}

	num, err = qs.Filter("user_name__endswith", "E").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, shouldNum))

	num, err = qs.Filter("user_name__iendswith", "E").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	num, err = qs.Filter("profile__isnull", true).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("status__in", 1, 2).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	num, err = qs.Filter("status__in", []int{1, 2}).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	n1, n2 := 1, 2
	num, err = qs.Filter("status__in", []*int{&n1}, &n2).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	num, err = qs.Filter("id__between", 2, 3).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	num, err = qs.Filter("id__between", []int{2, 3}).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))
}

func TestSetCond(t *testing.T) {
	cond := NewCondition()
	cond1 := cond.And("profile__isnull", false).AndNot("status__in", 1).Or("profile__age__gt", 2000)

	qs := dORM.QueryTable("user")
	num, err := qs.SetCond(cond1).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	cond2 := cond.AndCond(cond1).OrCond(cond.And("user_name", "slene"))
	num, err = qs.SetCond(cond2).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))
}

func TestLimit(t *testing.T) {
	var posts []*Post
	qs := dORM.QueryTable("post")
	num, err := qs.Limit(1).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Limit(-1).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 4))

	num, err = qs.Limit(-1, 2).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	num, err = qs.Limit(0, 2).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))
}

func TestOffset(t *testing.T) {
	var posts []*Post
	qs := dORM.QueryTable("post")
	num, err := qs.Limit(1).Offset(2).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Offset(2).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))
}

func TestOrderBy(t *testing.T) {
	qs := dORM.QueryTable("user")
	num, err := qs.OrderBy("-status").Filter("user_name", "nobody").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.OrderBy("status").Filter("user_name", "slene").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.OrderBy("-profile__age").Filter("user_name", "astaxie").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))
}

func TestAll(t *testing.T) {
	var users []*User
	qs := dORM.QueryTable("user")
	num, err := qs.OrderBy("Id").All(&users)
	throwFail(t, err)
	throwFailNow(t, AssertIs(num, 3))

	throwFail(t, AssertIs(users[0].UserName, "slene"))
	throwFail(t, AssertIs(users[1].UserName, "astaxie"))
	throwFail(t, AssertIs(users[2].UserName, "nobody"))

	var users2 []User
	qs = dORM.QueryTable("user")
	num, err = qs.OrderBy("Id").All(&users2)
	throwFail(t, err)
	throwFailNow(t, AssertIs(num, 3))

	throwFailNow(t, AssertIs(users2[0].UserName, "slene"))
	throwFailNow(t, AssertIs(users2[1].UserName, "astaxie"))
	throwFailNow(t, AssertIs(users2[2].UserName, "nobody"))

	qs = dORM.QueryTable("user")
	num, err = qs.OrderBy("Id").RelatedSel().All(&users2, "UserName")
	throwFail(t, err)
	throwFailNow(t, AssertIs(num, 3))
	throwFailNow(t, AssertIs(len(users2), 3))
	throwFailNow(t, AssertIs(users2[0].UserName, "slene"))
	throwFailNow(t, AssertIs(users2[1].UserName, "astaxie"))
	throwFailNow(t, AssertIs(users2[2].UserName, "nobody"))
	throwFailNow(t, AssertIs(users2[0].Id, 0))
	throwFailNow(t, AssertIs(users2[1].Id, 0))
	throwFailNow(t, AssertIs(users2[2].Id, 0))
	throwFailNow(t, AssertIs(users2[0].Profile == nil, false))
	throwFailNow(t, AssertIs(users2[1].Profile == nil, false))
	throwFailNow(t, AssertIs(users2[2].Profile == nil, true))

	qs = dORM.QueryTable("user")
	num, err = qs.Filter("user_name", "nothing").All(&users)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 0))

	var users3 []*User
	qs = dORM.QueryTable("user")
	num, err = qs.Filter("user_name", "nothing").All(&users3)
	throwFailNow(t, AssertIs(users3 == nil, false))
}

func TestOne(t *testing.T) {
	var user User
	qs := dORM.QueryTable("user")
	err := qs.One(&user)
	throwFail(t, AssertIs(err, ErrMultiRows))

	user = User{}
	err = qs.OrderBy("Id").Limit(1).One(&user)
	throwFailNow(t, err)
	throwFail(t, AssertIs(user.UserName, "slene"))

	err = qs.Filter("user_name", "nothing").One(&user)
	throwFail(t, AssertIs(err, ErrNoRows))

}

func TestValues(t *testing.T) {
	var maps []Params
	qs := dORM.QueryTable("user")

	num, err := qs.OrderBy("Id").Values(&maps)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 3))
	if num == 3 {
		throwFail(t, AssertIs(maps[0]["UserName"], "slene"))
		throwFail(t, AssertIs(maps[2]["Profile"], nil))
	}

	num, err = qs.OrderBy("Id").Values(&maps, "UserName", "Profile__Age")
	throwFail(t, err)
	throwFail(t, AssertIs(num, 3))
	if num == 3 {
		throwFail(t, AssertIs(maps[0]["UserName"], "slene"))
		throwFail(t, AssertIs(maps[0]["Profile__Age"], 28))
		throwFail(t, AssertIs(maps[2]["Profile__Age"], nil))
	}

	num, err = qs.Filter("UserName", "slene").Values(&maps)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))
}

func TestValuesList(t *testing.T) {
	var list []ParamsList
	qs := dORM.QueryTable("user")

	num, err := qs.OrderBy("Id").ValuesList(&list)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 3))
	if num == 3 {
		throwFail(t, AssertIs(list[0][1], "slene"))
		throwFail(t, AssertIs(list[2][9], nil))
	}

	num, err = qs.OrderBy("Id").ValuesList(&list, "UserName", "Profile__Age")
	throwFail(t, err)
	throwFail(t, AssertIs(num, 3))
	if num == 3 {
		throwFail(t, AssertIs(list[0][0], "slene"))
		throwFail(t, AssertIs(list[0][1], 28))
		throwFail(t, AssertIs(list[2][1], nil))
	}
}

func TestValuesFlat(t *testing.T) {
	var list ParamsList
	qs := dORM.QueryTable("user")

	num, err := qs.OrderBy("id").ValuesFlat(&list, "UserName")
	throwFail(t, err)
	throwFail(t, AssertIs(num, 3))
	if num == 3 {
		throwFail(t, AssertIs(list[0], "slene"))
		throwFail(t, AssertIs(list[1], "astaxie"))
		throwFail(t, AssertIs(list[2], "nobody"))
	}
}

func TestRelatedSel(t *testing.T) {
	qs := dORM.QueryTable("user")
	num, err := qs.Filter("profile__age", 28).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("profile__age__gt", 28).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("profile__user__profile__age__gt", 28).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	var user User
	err = qs.Filter("user_name", "slene").RelatedSel("profile").One(&user)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))
	throwFail(t, AssertNot(user.Profile, nil))
	if user.Profile != nil {
		throwFail(t, AssertIs(user.Profile.Age, 28))
	}

	err = qs.Filter("user_name", "slene").RelatedSel().One(&user)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))
	throwFail(t, AssertNot(user.Profile, nil))
	if user.Profile != nil {
		throwFail(t, AssertIs(user.Profile.Age, 28))
	}

	err = qs.Filter("user_name", "nobody").RelatedSel("profile").One(&user)
	throwFail(t, AssertIs(num, 1))
	throwFail(t, AssertIs(user.Profile, nil))

	qs = dORM.QueryTable("user_profile")
	num, err = qs.Filter("user__username", "slene").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	var posts []*Post
	qs = dORM.QueryTable("post")
	num, err = qs.RelatedSel().All(&posts)
	throwFail(t, err)
	throwFailNow(t, AssertIs(num, 4))

	throwFailNow(t, AssertIs(posts[0].User.UserName, "slene"))
	throwFailNow(t, AssertIs(posts[1].User.UserName, "astaxie"))
	throwFailNow(t, AssertIs(posts[2].User.UserName, "astaxie"))
	throwFailNow(t, AssertIs(posts[3].User.UserName, "nobody"))
}

func TestReverseQuery(t *testing.T) {
	var profile Profile
	err := dORM.QueryTable("user_profile").Filter("User", 3).One(&profile)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(profile.Age, 30))

	profile = Profile{}
	err = dORM.QueryTable("user_profile").Filter("User__UserName", "astaxie").One(&profile)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(profile.Age, 30))

	var user User
	err = dORM.QueryTable("user").Filter("Posts__Title", "Examples").One(&user)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(user.UserName, "astaxie"))

	user = User{}
	err = dORM.QueryTable("user").Filter("Posts__User__UserName", "astaxie").Limit(1).One(&user)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(user.UserName, "astaxie"))

	user = User{}
	err = dORM.QueryTable("user").Filter("Posts__User__UserName", "astaxie").RelatedSel().Limit(1).One(&user)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(user.UserName, "astaxie"))
	throwFailNow(t, AssertIs(user.Profile == nil, false))
	throwFailNow(t, AssertIs(user.Profile.Age, 30))

	var posts []*Post
	num, err := dORM.QueryTable("post").Filter("Tags__Tag__Name", "golang").All(&posts)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 3))
	throwFailNow(t, AssertIs(posts[0].Title, "Introduction"))

	posts = []*Post{}
	num, err = dORM.QueryTable("post").Filter("Tags__Tag__Name", "golang").Filter("User__UserName", "slene").All(&posts)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(posts[0].Title, "Introduction"))

	posts = []*Post{}
	num, err = dORM.QueryTable("post").Filter("Tags__Tag__Name", "golang").
		Filter("User__UserName", "slene").RelatedSel().All(&posts)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(posts[0].User == nil, false))
	throwFailNow(t, AssertIs(posts[0].User.UserName, "slene"))

	var tags []*Tag
	num, err = dORM.QueryTable("tag").Filter("Posts__Post__Title", "Introduction").All(&tags)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(tags[0].Name, "golang"))

	tags = []*Tag{}
	num, err = dORM.QueryTable("tag").Filter("Posts__Post__Title", "Introduction").
		Filter("BestPost__User__UserName", "astaxie").All(&tags)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(tags[0].Name, "golang"))

	tags = []*Tag{}
	num, err = dORM.QueryTable("tag").Filter("Posts__Post__Title", "Introduction").
		Filter("BestPost__User__UserName", "astaxie").RelatedSel().All(&tags)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(tags[0].Name, "golang"))
	throwFailNow(t, AssertIs(tags[0].BestPost == nil, false))
	throwFailNow(t, AssertIs(tags[0].BestPost.Title, "Examples"))
	throwFailNow(t, AssertIs(tags[0].BestPost.User == nil, false))
	throwFailNow(t, AssertIs(tags[0].BestPost.User.UserName, "astaxie"))
}

func TestLoadRelated(t *testing.T) {
	// load reverse foreign key
	user := User{Id: 3}

	err := dORM.Read(&user)
	throwFailNow(t, err)

	num, err := dORM.LoadRelated(&user, "Posts")
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 2))
	throwFailNow(t, AssertIs(len(user.Posts), 2))
	throwFailNow(t, AssertIs(user.Posts[0].User.Id, 3))

	num, err = dORM.LoadRelated(&user, "Posts", true)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(len(user.Posts), 2))
	throwFailNow(t, AssertIs(user.Posts[0].User.UserName, "astaxie"))

	num, err = dORM.LoadRelated(&user, "Posts", true, 1)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(len(user.Posts), 1))

	num, err = dORM.LoadRelated(&user, "Posts", true, 0, 0, "-Id")
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(len(user.Posts), 2))
	throwFailNow(t, AssertIs(user.Posts[0].Title, "Formatting"))

	num, err = dORM.LoadRelated(&user, "Posts", true, 1, 1, "Id")
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(len(user.Posts), 1))
	throwFailNow(t, AssertIs(user.Posts[0].Title, "Formatting"))

	// load reverse one to one
	profile := Profile{Id: 3}
	profile.BestPost = &Post{Id: 2}
	num, err = dORM.Update(&profile, "BestPost")
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))

	err = dORM.Read(&profile)
	throwFailNow(t, err)

	num, err = dORM.LoadRelated(&profile, "User")
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(profile.User == nil, false))
	throwFailNow(t, AssertIs(profile.User.UserName, "astaxie"))

	num, err = dORM.LoadRelated(&profile, "User", true)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(profile.User == nil, false))
	throwFailNow(t, AssertIs(profile.User.UserName, "astaxie"))
	throwFailNow(t, AssertIs(profile.User.Profile.Age, profile.Age))

	// load rel one to one
	err = dORM.Read(&user)
	throwFailNow(t, err)

	num, err = dORM.LoadRelated(&user, "Profile")
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(user.Profile == nil, false))
	throwFailNow(t, AssertIs(user.Profile.Age, 30))

	num, err = dORM.LoadRelated(&user, "Profile", true)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(user.Profile == nil, false))
	throwFailNow(t, AssertIs(user.Profile.Age, 30))
	throwFailNow(t, AssertIs(user.Profile.BestPost == nil, false))
	throwFailNow(t, AssertIs(user.Profile.BestPost.Title, "Examples"))

	post := Post{Id: 2}

	// load rel foreign key
	err = dORM.Read(&post)
	throwFailNow(t, err)

	num, err = dORM.LoadRelated(&post, "User")
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(post.User == nil, false))
	throwFailNow(t, AssertIs(post.User.UserName, "astaxie"))

	num, err = dORM.LoadRelated(&post, "User", true)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(post.User == nil, false))
	throwFailNow(t, AssertIs(post.User.UserName, "astaxie"))
	throwFailNow(t, AssertIs(post.User.Profile == nil, false))
	throwFailNow(t, AssertIs(post.User.Profile.Age, 30))

	// load rel m2m
	post = Post{Id: 2}

	err = dORM.Read(&post)
	throwFailNow(t, err)

	num, err = dORM.LoadRelated(&post, "Tags")
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 2))
	throwFailNow(t, AssertIs(len(post.Tags), 2))
	throwFailNow(t, AssertIs(post.Tags[0].Name, "golang"))

	num, err = dORM.LoadRelated(&post, "Tags", true)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 2))
	throwFailNow(t, AssertIs(len(post.Tags), 2))
	throwFailNow(t, AssertIs(post.Tags[0].Name, "golang"))
	throwFailNow(t, AssertIs(post.Tags[0].BestPost == nil, false))
	throwFailNow(t, AssertIs(post.Tags[0].BestPost.User.UserName, "astaxie"))

	// load reverse m2m
	tag := Tag{Id: 1}

	err = dORM.Read(&tag)
	throwFailNow(t, err)

	num, err = dORM.LoadRelated(&tag, "Posts")
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 3))
	throwFailNow(t, AssertIs(tag.Posts[0].Title, "Introduction"))
	throwFailNow(t, AssertIs(tag.Posts[0].User.Id, 2))
	throwFailNow(t, AssertIs(tag.Posts[0].User.Profile == nil, true))

	num, err = dORM.LoadRelated(&tag, "Posts", true)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 3))
	throwFailNow(t, AssertIs(tag.Posts[0].Title, "Introduction"))
	throwFailNow(t, AssertIs(tag.Posts[0].User.Id, 2))
	throwFailNow(t, AssertIs(tag.Posts[0].User.UserName, "slene"))
}

func TestQueryM2M(t *testing.T) {
	post := Post{Id: 4}
	m2m := dORM.QueryM2M(&post, "Tags")

	tag1 := []*Tag{&Tag{Name: "TestTag1"}, &Tag{Name: "TestTag2"}}
	tag2 := &Tag{Name: "TestTag3"}
	tag3 := []interface{}{&Tag{Name: "TestTag4"}}

	tags := []interface{}{tag1[0], tag1[1], tag2, tag3[0]}

	for _, tag := range tags {
		_, err := dORM.Insert(tag)
		throwFailNow(t, err)
	}

	num, err := m2m.Add(tag1)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 2))

	num, err = m2m.Add(tag2)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))

	num, err = m2m.Add(tag3)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))

	num, err = m2m.Count()
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 5))

	num, err = m2m.Remove(tag3)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))

	num, err = m2m.Count()
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 4))

	exist := m2m.Exist(tag2)
	throwFailNow(t, AssertIs(exist, true))

	num, err = m2m.Remove(tag2)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))

	exist = m2m.Exist(tag2)
	throwFailNow(t, AssertIs(exist, false))

	num, err = m2m.Count()
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 3))

	num, err = m2m.Clear()
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 3))

	num, err = m2m.Count()
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 0))

	tag := Tag{Name: "test"}
	_, err = dORM.Insert(&tag)
	throwFailNow(t, err)

	m2m = dORM.QueryM2M(&tag, "Posts")

	post1 := []*Post{&Post{Title: "TestPost1"}, &Post{Title: "TestPost2"}}
	post2 := &Post{Title: "TestPost3"}
	post3 := []interface{}{&Post{Title: "TestPost4"}}

	posts := []interface{}{post1[0], post1[1], post2, post3[0]}

	for _, post := range posts {
		p := post.(*Post)
		p.User = &User{Id: 1}
		_, err := dORM.Insert(post)
		throwFailNow(t, err)
	}

	num, err = m2m.Add(post1)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 2))

	num, err = m2m.Add(post2)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))

	num, err = m2m.Add(post3)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))

	num, err = m2m.Count()
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 4))

	num, err = m2m.Remove(post3)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))

	num, err = m2m.Count()
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 3))

	exist = m2m.Exist(post2)
	throwFailNow(t, AssertIs(exist, true))

	num, err = m2m.Remove(post2)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))

	exist = m2m.Exist(post2)
	throwFailNow(t, AssertIs(exist, false))

	num, err = m2m.Count()
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 2))

	num, err = m2m.Clear()
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 2))

	num, err = m2m.Count()
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 0))

	num, err = dORM.Delete(&tag)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
}

func TestQueryRelate(t *testing.T) {
	// post := &Post{Id: 2}

	// qs := dORM.QueryRelate(post, "Tags")
	// num, err := qs.Count()
	// throwFailNow(t, err)
	// throwFailNow(t, AssertIs(num, 2))

	// var tags []*Tag
	// num, err = qs.All(&tags)
	// throwFailNow(t, err)
	// throwFailNow(t, AssertIs(num, 2))
	// throwFailNow(t, AssertIs(tags[0].Name, "golang"))

	// num, err = dORM.QueryTable("Tag").Filter("Posts__Post", 2).Count()
	// throwFailNow(t, err)
	// throwFailNow(t, AssertIs(num, 2))
}

func TestPrepareInsert(t *testing.T) {
	qs := dORM.QueryTable("user")
	i, err := qs.PrepareInsert()
	throwFailNow(t, err)

	var user User
	user.UserName = "testing1"
	num, err := i.Insert(&user)
	throwFail(t, err)
	throwFail(t, AssertIs(num > 0, true))

	user.UserName = "testing2"
	num, err = i.Insert(&user)
	throwFail(t, err)
	throwFail(t, AssertIs(num > 0, true))

	num, err = qs.Filter("user_name__in", "testing1", "testing2").Delete()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 2))

	err = i.Close()
	throwFail(t, err)
	err = i.Close()
	throwFail(t, AssertIs(err, ErrStmtClosed))
}

func TestRawExec(t *testing.T) {
	Q := dDbBaser.TableQuote()

	query := fmt.Sprintf("UPDATE %suser%s SET %suser_name%s = ? WHERE %suser_name%s = ?", Q, Q, Q, Q, Q, Q)
	res, err := dORM.Raw(query, "testing", "slene").Exec()
	throwFail(t, err)
	num, err := res.RowsAffected()
	throwFail(t, AssertIs(num, 1), err)

	res, err = dORM.Raw(query, "slene", "testing").Exec()
	throwFail(t, err)
	num, err = res.RowsAffected()
	throwFail(t, AssertIs(num, 1), err)
}

func TestRawQueryRow(t *testing.T) {
	var (
		Boolean  bool
		Char     string
		Text     string
		Date     time.Time
		DateTime time.Time
		Byte     byte
		Rune     rune
		Int      int
		Int8     int
		Int16    int16
		Int32    int32
		Int64    int64
		Uint     uint
		Uint8    uint8
		Uint16   uint16
		Uint32   uint32
		Uint64   uint64
		Float32  float32
		Float64  float64
		Decimal  float64
	)

	data_values := make(map[string]interface{}, len(Data_Values))

	for k, v := range Data_Values {
		data_values[strings.ToLower(k)] = v
	}

	Q := dDbBaser.TableQuote()

	cols := []string{
		"id", "boolean", "char", "text", "date", "datetime", "byte", "rune", "int", "int8", "int16", "int32",
		"int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64", "decimal",
	}
	sep := fmt.Sprintf("%s, %s", Q, Q)
	query := fmt.Sprintf("SELECT %s%s%s FROM data WHERE id = ?", Q, strings.Join(cols, sep), Q)
	var id int
	values := []interface{}{
		&id, &Boolean, &Char, &Text, &Date, &DateTime, &Byte, &Rune, &Int, &Int8, &Int16, &Int32,
		&Int64, &Uint, &Uint8, &Uint16, &Uint32, &Uint64, &Float32, &Float64, &Decimal,
	}
	err := dORM.Raw(query, 1).QueryRow(values...)
	throwFailNow(t, err)
	for i, col := range cols {
		vu := values[i]
		v := reflect.ValueOf(vu).Elem().Interface()
		switch col {
		case "id":
			throwFail(t, AssertIs(id, 1))
		case "date":
			v = v.(time.Time).In(DefaultTimeLoc)
			value := data_values[col].(time.Time).In(DefaultTimeLoc)
			throwFail(t, AssertIs(v, value, test_Date))
		case "datetime":
			v = v.(time.Time).In(DefaultTimeLoc)
			value := data_values[col].(time.Time).In(DefaultTimeLoc)
			throwFail(t, AssertIs(v, value, test_DateTime))
		default:
			throwFail(t, AssertIs(v, data_values[col]))
		}
	}

	var (
		uid    int
		status *int
		pid    *int
	)

	cols = []string{
		"id", "Status", "profile_id",
	}
	query = fmt.Sprintf("SELECT %s%s%s FROM %suser%s WHERE id = ?", Q, strings.Join(cols, sep), Q, Q, Q)
	err = dORM.Raw(query, 4).QueryRow(&uid, &status, &pid)
	throwFail(t, err)
	throwFail(t, AssertIs(uid, 4))
	throwFail(t, AssertIs(*status, 3))
	throwFail(t, AssertIs(pid, nil))
}

func TestQueryRows(t *testing.T) {
	Q := dDbBaser.TableQuote()

	var datas []*Data

	query := fmt.Sprintf("SELECT * FROM %sdata%s", Q, Q)
	num, err := dORM.Raw(query).QueryRows(&datas)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(len(datas), 1))

	ind := reflect.Indirect(reflect.ValueOf(datas[0]))

	for name, value := range Data_Values {
		e := ind.FieldByName(name)
		vu := e.Interface()
		switch name {
		case "Date":
			vu = vu.(time.Time).In(DefaultTimeLoc).Format(test_Date)
			value = value.(time.Time).In(DefaultTimeLoc).Format(test_Date)
		case "DateTime":
			vu = vu.(time.Time).In(DefaultTimeLoc).Format(test_DateTime)
			value = value.(time.Time).In(DefaultTimeLoc).Format(test_DateTime)
		}
		throwFail(t, AssertIs(vu == value, true), value, vu)
	}

	var datas2 []Data

	query = fmt.Sprintf("SELECT * FROM %sdata%s", Q, Q)
	num, err = dORM.Raw(query).QueryRows(&datas2)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(len(datas2), 1))

	ind = reflect.Indirect(reflect.ValueOf(datas2[0]))

	for name, value := range Data_Values {
		e := ind.FieldByName(name)
		vu := e.Interface()
		switch name {
		case "Date":
			vu = vu.(time.Time).In(DefaultTimeLoc).Format(test_Date)
			value = value.(time.Time).In(DefaultTimeLoc).Format(test_Date)
		case "DateTime":
			vu = vu.(time.Time).In(DefaultTimeLoc).Format(test_DateTime)
			value = value.(time.Time).In(DefaultTimeLoc).Format(test_DateTime)
		}
		throwFail(t, AssertIs(vu == value, true), value, vu)
	}

	var ids []int
	var usernames []string
	query = fmt.Sprintf("SELECT %sid%s, %suser_name%s FROM %suser%s ORDER BY %sid%s ASC", Q, Q, Q, Q, Q, Q, Q, Q)
	num, err = dORM.Raw(query).QueryRows(&ids, &usernames)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, 3))
	throwFailNow(t, AssertIs(len(ids), 3))
	throwFailNow(t, AssertIs(ids[0], 2))
	throwFailNow(t, AssertIs(usernames[0], "slene"))
	throwFailNow(t, AssertIs(ids[1], 3))
	throwFailNow(t, AssertIs(usernames[1], "astaxie"))
	throwFailNow(t, AssertIs(ids[2], 4))
	throwFailNow(t, AssertIs(usernames[2], "nobody"))
}

func TestRawValues(t *testing.T) {
	Q := dDbBaser.TableQuote()

	var maps []Params
	query := fmt.Sprintf("SELECT %suser_name%s FROM %suser%s WHERE %sStatus%s = ?", Q, Q, Q, Q, Q, Q)
	num, err := dORM.Raw(query, 1).Values(&maps)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))
	if num == 1 {
		throwFail(t, AssertIs(maps[0]["user_name"], "slene"))
	}

	var lists []ParamsList
	num, err = dORM.Raw(query, 1).ValuesList(&lists)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))
	if num == 1 {
		throwFail(t, AssertIs(lists[0][0], "slene"))
	}

	query = fmt.Sprintf("SELECT %sprofile_id%s FROM %suser%s ORDER BY %sid%s ASC", Q, Q, Q, Q, Q, Q)
	var list ParamsList
	num, err = dORM.Raw(query).ValuesFlat(&list)
	throwFail(t, err)
	throwFail(t, AssertIs(num, 3))
	if num == 3 {
		throwFail(t, AssertIs(list[0], "2"))
		throwFail(t, AssertIs(list[1], "3"))
		throwFail(t, AssertIs(list[2], nil))
	}
}

func TestRawPrepare(t *testing.T) {
	switch {
	case IsMysql || IsSqlite:

		pre, err := dORM.Raw("INSERT INTO tag (name) VALUES (?)").Prepare()
		throwFail(t, err)
		if pre != nil {
			r, err := pre.Exec("name1")
			throwFail(t, err)

			tid, err := r.LastInsertId()
			throwFail(t, err)
			throwFail(t, AssertIs(tid > 0, true))

			r, err = pre.Exec("name2")
			throwFail(t, err)

			id, err := r.LastInsertId()
			throwFail(t, err)
			throwFail(t, AssertIs(id, tid+1))

			r, err = pre.Exec("name3")
			throwFail(t, err)

			id, err = r.LastInsertId()
			throwFail(t, err)
			throwFail(t, AssertIs(id, tid+2))

			err = pre.Close()
			throwFail(t, err)

			res, err := dORM.Raw("DELETE FROM tag WHERE name IN (?, ?, ?)", []string{"name1", "name2", "name3"}).Exec()
			throwFail(t, err)

			num, err := res.RowsAffected()
			throwFail(t, err)
			throwFail(t, AssertIs(num, 3))
		}

	case IsPostgres:

		pre, err := dORM.Raw(`INSERT INTO "tag" ("name") VALUES (?) RETURNING "id"`).Prepare()
		throwFail(t, err)
		if pre != nil {
			_, err := pre.Exec("name1")
			throwFail(t, err)

			_, err = pre.Exec("name2")
			throwFail(t, err)

			_, err = pre.Exec("name3")
			throwFail(t, err)

			err = pre.Close()
			throwFail(t, err)

			res, err := dORM.Raw(`DELETE FROM "tag" WHERE "name" IN (?, ?, ?)`, []string{"name1", "name2", "name3"}).Exec()
			throwFail(t, err)

			if err == nil {
				num, err := res.RowsAffected()
				throwFail(t, err)
				throwFail(t, AssertIs(num, 3))
			}
		}
	}
}

func TestUpdate(t *testing.T) {
	qs := dORM.QueryTable("user")
	num, err := qs.Filter("user_name", "slene").Filter("is_staff", false).Update(Params{
		"is_staff":  true,
		"is_active": true,
	})
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	// with join
	num, err = qs.Filter("user_name", "slene").Filter("profile__age", 28).Filter("is_staff", true).Update(Params{
		"is_staff": false,
	})
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("user_name", "slene").Update(Params{
		"Nums": ColValue(Col_Add, 100),
	})
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("user_name", "slene").Update(Params{
		"Nums": ColValue(Col_Minus, 50),
	})
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("user_name", "slene").Update(Params{
		"Nums": ColValue(Col_Multiply, 3),
	})
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	num, err = qs.Filter("user_name", "slene").Update(Params{
		"Nums": ColValue(Col_Except, 5),
	})
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	user := User{UserName: "slene"}
	err = dORM.Read(&user, "UserName")
	throwFail(t, err)
	throwFail(t, AssertIs(user.Nums, 30))
}

func TestDelete(t *testing.T) {
	qs := dORM.QueryTable("user_profile")
	num, err := qs.Filter("user__user_name", "slene").Delete()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	qs = dORM.QueryTable("user")
	num, err = qs.Filter("user_name", "slene").Filter("profile__isnull", true).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	qs = dORM.QueryTable("comment")
	num, err = qs.Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 6))

	qs = dORM.QueryTable("post")
	num, err = qs.Filter("Id", 3).Delete()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	qs = dORM.QueryTable("comment")
	num, err = qs.Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 4))

	qs = dORM.QueryTable("comment")
	num, err = qs.Filter("Post__User", 3).Delete()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 3))

	qs = dORM.QueryTable("comment")
	num, err = qs.Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))
}

func TestTransaction(t *testing.T) {
	// this test worked when database support transaction

	o := NewOrm()
	err := o.Begin()
	throwFail(t, err)

	var names = []string{"1", "2", "3"}

	var tag Tag
	tag.Name = names[0]
	id, err := o.Insert(&tag)
	throwFail(t, err)
	throwFail(t, AssertIs(id > 0, true))

	num, err := o.QueryTable("tag").Filter("name", "golang").Update(Params{"name": names[1]})
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

	switch {
	case IsMysql || IsSqlite:
		res, err := o.Raw("INSERT INTO tag (name) VALUES (?)", names[2]).Exec()
		throwFail(t, err)
		if err == nil {
			id, err = res.LastInsertId()
			throwFail(t, err)
			throwFail(t, AssertIs(id > 0, true))
		}
	}

	err = o.Rollback()
	throwFail(t, err)

	num, err = o.QueryTable("tag").Filter("name__in", names).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 0))

	err = o.Begin()
	throwFail(t, err)

	tag.Name = "commit"
	id, err = o.Insert(&tag)
	throwFail(t, err)
	throwFail(t, AssertIs(id > 0, true))

	o.Commit()
	throwFail(t, err)

	num, err = o.QueryTable("tag").Filter("name", "commit").Delete()
	throwFail(t, err)
	throwFail(t, AssertIs(num, 1))

}

func TestReadOrCreate(t *testing.T) {
	u := &User{
		UserName: "Kyle",
		Email:    "kylemcc@gmail.com",
		Password: "other_pass",
		Status:   7,
		IsStaff:  false,
		IsActive: true,
	}

	created, pk, err := dORM.ReadOrCreate(u, "UserName")
	throwFail(t, err)
	throwFail(t, AssertIs(created, true))
	throwFail(t, AssertIs(u.UserName, "Kyle"))
	throwFail(t, AssertIs(u.Email, "kylemcc@gmail.com"))
	throwFail(t, AssertIs(u.Password, "other_pass"))
	throwFail(t, AssertIs(u.Status, 7))
	throwFail(t, AssertIs(u.IsStaff, false))
	throwFail(t, AssertIs(u.IsActive, true))
	throwFail(t, AssertIs(u.Created.In(DefaultTimeLoc), u.Created.In(DefaultTimeLoc), test_Date))
	throwFail(t, AssertIs(u.Updated.In(DefaultTimeLoc), u.Updated.In(DefaultTimeLoc), test_DateTime))

	nu := &User{UserName: u.UserName, Email: "someotheremail@gmail.com"}
	created, pk, err = dORM.ReadOrCreate(nu, "UserName")
	throwFail(t, err)
	throwFail(t, AssertIs(created, false))
	throwFail(t, AssertIs(nu.Id, u.Id))
	throwFail(t, AssertIs(pk, u.Id))
	throwFail(t, AssertIs(nu.UserName, u.UserName))
	throwFail(t, AssertIs(nu.Email, u.Email)) // should contain the value in the table, not the one specified above
	throwFail(t, AssertIs(nu.Password, u.Password))
	throwFail(t, AssertIs(nu.Status, u.Status))
	throwFail(t, AssertIs(nu.IsStaff, u.IsStaff))
	throwFail(t, AssertIs(nu.IsActive, u.IsActive))

	dORM.Delete(u)
}
