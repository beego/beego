package orm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

type T_Code int

const (
	// =
	T_Equal T_Code = iota
	// <
	T_Less
	// >
	T_Large
	// elment in slice/array
	// T_In
	// key exists in map
	// T_KeyExist
	// index != -1
	// T_Contain
	// index == 0
	// T_StartWith
	// index == len(x) - 1
	// T_EndWith
)

func ValuesCompare(is bool, a interface{}, o T_Code, args ...interface{}) (err error, ok bool) {
	if len(args) == 0 {
		return fmt.Errorf("miss args"), false
	}
	b := args[0]
	arg := argAny(args)
	switch o {
	case T_Equal:
		switch v := a.(type) {
		case reflect.Kind:
			ok = reflect.ValueOf(b).Kind() == v
		case time.Time:
			if v2, vo := b.(time.Time); vo {
				if arg.Get(1) != nil {
					format := ToStr(arg.Get(1))
					ok = v.Format(format) == v2.Format(format)
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
				err = fmt.Errorf("should: a == b, a = `%v`, b = `%v`", a, b)
			} else {
				err = fmt.Errorf("should: a != b, a = `%v`, b = `%v`", a, b)
			}
		}
	case T_Less, T_Large:
		as := ToStr(a)
		bs := ToStr(b)
		f1, er := StrTo(as).Float64()
		if er != nil {
			err = fmt.Errorf("wrong type need numeric: `%v`", a)
			goto wrongArg
		}
		f2, er := StrTo(bs).Float64()
		if er != nil {
			err = fmt.Errorf("wrong type need numeric: `%v`", b)
			goto wrongArg
		}
		var opts []string
		if o == T_Less {
			opts = []string{"<", ">="}
			ok = f1 < f2
		} else {
			opts = []string{">", "<="}
			ok = f1 > f2
		}
		ok = is && ok || !is && !ok
		if !ok {
			if is {
				err = fmt.Errorf("should: a %s b, a = `%v`, b = `%v`", opts[0], f1, f2)
			} else {
				err = fmt.Errorf("should: a %s b, a = `%v`, b = `%v`", opts[1], f1, f2)
			}
		}
	}
wrongArg:
	if err != nil {
		return err, false
	}

	return nil, true
}

func AssertIs(a interface{}, o T_Code, args ...interface{}) error {
	if err, ok := ValuesCompare(true, a, o, args...); ok == false {
		return err
	}
	return nil
}

func AssertNot(a interface{}, o T_Code, args ...interface{}) error {
	if err, ok := ValuesCompare(false, a, o, args...); ok == false {
		return err
	}
	return nil
}

func getCaller(skip int) string {
	pc, file, line, _ := runtime.Caller(skip)
	fun := runtime.FuncForPC(pc)
	_, fn := filepath.Split(file)
	data, err := ioutil.ReadFile(file)
	code := ""
	if err == nil {
		lines := bytes.Split(data, []byte{'\n'})
		code = strings.TrimSpace(string(lines[line-1]))
	}
	funName := fun.Name()
	if i := strings.LastIndex(funName, "."); i > -1 {
		funName = funName[i+1:]
	}
	return fmt.Sprintf("%s:%d: %s: %s", fn, line, funName, code)
}

func throwFail(t *testing.T, err error, args ...interface{}) {
	if err != nil {
		params := []interface{}{"\n", getCaller(2), "\n", err, "\n"}
		params = append(params, args...)
		t.Error(params...)
		t.Fail()
	}
}

func throwFailNow(t *testing.T, err error, args ...interface{}) {
	if err != nil {
		params := []interface{}{"\n", getCaller(2), "\n", err, "\n"}
		params = append(params, args...)
		t.Error(params...)
		t.FailNow()
	}
}

func TestModelSyntax(t *testing.T) {
	mi, ok := modelCache.get("user")
	throwFail(t, AssertIs(ok, T_Equal, true))
	if ok {
		throwFail(t, AssertIs(mi.fields.GetByName("ShouldSkip") == nil, T_Equal, true))
	}
}

func TestCRUD(t *testing.T) {
	profile := NewProfile()
	profile.Age = 30
	profile.Money = 1234.12
	id, err := dORM.Insert(profile)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(id, T_Large, 0))

	user := NewUser()
	user.UserName = "slene"
	user.Email = "vslene@gmail.com"
	user.Password = "pass"
	user.Status = 3
	user.IsStaff = true
	user.IsActive = true

	id, err = dORM.Insert(user)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(id, T_Large, 0))

	u := &User{Id: user.Id}
	err = dORM.Read(u)
	throwFailNow(t, err)

	throwFailNow(t, AssertIs(u.UserName, T_Equal, "slene"))
	throwFailNow(t, AssertIs(u.Email, T_Equal, "vslene@gmail.com"))
	throwFailNow(t, AssertIs(u.Password, T_Equal, "pass"))
	throwFailNow(t, AssertIs(u.Status, T_Equal, 3))
	throwFailNow(t, AssertIs(u.IsStaff, T_Equal, true))
	throwFailNow(t, AssertIs(u.IsActive, T_Equal, true))
	throwFailNow(t, AssertIs(u.Created, T_Equal, user.Created, format_Date))
	throwFailNow(t, AssertIs(u.Updated, T_Equal, user.Updated, format_DateTime))

	user.UserName = "astaxie"
	user.Profile = profile
	num, err := dORM.Update(user)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, T_Equal, 1))

	u = &User{Id: user.Id}
	err = dORM.Read(u)
	throwFailNow(t, err)

	throwFailNow(t, AssertIs(u.UserName, T_Equal, "astaxie"))
	throwFailNow(t, AssertIs(u.Profile.Id, T_Equal, profile.Id))

	num, err = dORM.Delete(profile)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, T_Equal, 1))

	u = &User{Id: user.Id}
	err = dORM.Read(u)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(true, T_Equal, u.Profile == nil))

	num, err = dORM.Delete(user)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(num, T_Equal, 1))

	u = &User{Id: 100}
	err = dORM.Read(u)
	throwFailNow(t, AssertIs(err, T_Equal, ErrNoRows))
}

func TestInsertTestData(t *testing.T) {
	var users []*User

	profile := NewProfile()
	profile.Age = 28
	profile.Money = 1234.12

	id, err := dORM.Insert(profile)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(id, T_Large, 0))

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
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(id, T_Large, 0))

	profile = NewProfile()
	profile.Age = 30
	profile.Money = 4321.09

	id, err = dORM.Insert(profile)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(id, T_Large, 0))

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
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(id, T_Large, 0))

	user = NewUser()
	user.UserName = "nobody"
	user.Email = "nobody@gmail.com"
	user.Password = "nobody"
	user.Status = 3
	user.IsStaff = false
	user.IsActive = false

	users = append(users, user)

	id, err = dORM.Insert(user)
	throwFailNow(t, err)
	throwFailNow(t, AssertIs(id, T_Large, 0))

	tags := []*Tag{
		&Tag{Name: "golang"},
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
		throwFailNow(t, err)
		throwFailNow(t, AssertIs(id, T_Large, 0))
	}

	for _, post := range posts {
		id, err := dORM.Insert(post)
		throwFailNow(t, err)
		throwFailNow(t, AssertIs(id, T_Large, 0))
		// dORM.M2mAdd(post, "tags", post.Tags)
	}

	for _, comment := range comments {
		id, err := dORM.Insert(comment)
		throwFailNow(t, err)
		throwFailNow(t, AssertIs(id, T_Large, 0))
	}
}

func TestExpr(t *testing.T) {
	qs := dORM.QueryTable("User")
	qs = dORM.QueryTable("user")
	num, err := qs.Filter("UserName", "slene").Filter("user_name", "slene").Filter("profile__Age", 28).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))
}

func TestOperators(t *testing.T) {
	qs := dORM.QueryTable("user")
	num, err := qs.Filter("user_name", "slene").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.Filter("user_name__exact", "slene").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.Filter("user_name__iexact", "Slene").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.Filter("user_name__contains", "e").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))

	num, err = qs.Filter("user_name__contains", "E").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 0))

	num, err = qs.Filter("user_name__icontains", "E").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))

	num, err = qs.Filter("user_name__icontains", "E").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))

	num, err = qs.Filter("status__gt", 1).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))

	num, err = qs.Filter("status__gte", 1).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 3))

	num, err = qs.Filter("status__lt", 3).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))

	num, err = qs.Filter("status__lte", 3).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 3))

	num, err = qs.Filter("user_name__startswith", "s").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.Filter("user_name__startswith", "S").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 0))

	num, err = qs.Filter("user_name__istartswith", "S").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.Filter("user_name__endswith", "e").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))

	num, err = qs.Filter("user_name__endswith", "E").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 0))

	num, err = qs.Filter("user_name__iendswith", "E").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))

	num, err = qs.Filter("profile__isnull", true).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.Filter("status__in", 1, 2).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))

	num, err = qs.Filter("status__in", []int{1, 2}).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))

	n1, n2 := 1, 2
	num, err = qs.Filter("status__in", []*int{&n1}, &n2).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))
}

func TestAll(t *testing.T) {
	var users []*User
	qs := dORM.QueryTable("user")
	num, err := qs.All(&users)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 3))

	qs = dORM.QueryTable("user")
	num, err = qs.Filter("user_name", "nothing").All(&users)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 0))
}

func TestOne(t *testing.T) {
	var user User
	qs := dORM.QueryTable("user")
	err := qs.One(&user)
	throwFail(t, AssertIs(err, T_Equal, ErrMultiRows))

	err = qs.Filter("user_name", "nothing").One(&user)
	throwFail(t, AssertIs(err, T_Equal, ErrNoRows))
}

func TestValues(t *testing.T) {
	var maps []Params
	qs := dORM.QueryTable("user")

	num, err := qs.Values(&maps)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 3))
	if num == 3 {
		throwFail(t, AssertIs(maps[0]["UserName"], T_Equal, "slene"))
		throwFail(t, AssertIs(maps[2]["Profile"], T_Equal, nil))
	}

	num, err = qs.Values(&maps, "UserName", "Profile__Age")
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 3))
	if num == 3 {
		throwFail(t, AssertIs(maps[0]["UserName"], T_Equal, "slene"))
		throwFail(t, AssertIs(maps[0]["Profile__Age"], T_Equal, 28))
		throwFail(t, AssertIs(maps[2]["Profile__Age"], T_Equal, nil))
	}
}

func TestValuesList(t *testing.T) {
	var list []ParamsList
	qs := dORM.QueryTable("user")

	num, err := qs.ValuesList(&list)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 3))
	if num == 3 {
		throwFail(t, AssertIs(list[0][1], T_Equal, "slene"))
		throwFail(t, AssertIs(list[2][9], T_Equal, nil))
	}

	num, err = qs.ValuesList(&list, "UserName", "Profile__Age")
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 3))
	if num == 3 {
		throwFail(t, AssertIs(list[0][0], T_Equal, "slene"))
		throwFail(t, AssertIs(list[0][1], T_Equal, 28))
		throwFail(t, AssertIs(list[2][1], T_Equal, nil))
	}
}

func TestValuesFlat(t *testing.T) {
	var list ParamsList
	qs := dORM.QueryTable("user")

	num, err := qs.OrderBy("id").ValuesFlat(&list, "UserName")
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 3))
	if num == 3 {
		throwFail(t, AssertIs(list[0], T_Equal, "slene"))
		throwFail(t, AssertIs(list[1], T_Equal, "astaxie"))
		throwFail(t, AssertIs(list[2], T_Equal, "nobody"))
	}
}

func TestRelatedSel(t *testing.T) {
	qs := dORM.QueryTable("user")
	num, err := qs.Filter("profile__age", 28).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.Filter("profile__age__gt", 28).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.Filter("profile__user__profile__age__gt", 28).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	var user User
	err = qs.Filter("user_name", "slene").RelatedSel("profile").One(&user)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))
	throwFail(t, AssertNot(user.Profile, T_Equal, nil))
	if user.Profile != nil {
		throwFail(t, AssertIs(user.Profile.Age, T_Equal, 28))
	}

	err = qs.Filter("user_name", "slene").RelatedSel().One(&user)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))
	throwFail(t, AssertNot(user.Profile, T_Equal, nil))
	throwFail(t, AssertIs(user.Profile.Age, T_Equal, 28))
	if user.Profile != nil {
		throwFail(t, AssertIs(user.Profile.Age, T_Equal, 28))
	}

	err = qs.Filter("user_name", "nobody").RelatedSel("profile").One(&user)
	throwFail(t, AssertIs(num, T_Equal, 1))
	throwFail(t, AssertIs(user.Profile, T_Equal, nil))

	qs = dORM.QueryTable("user_profile")
	num, err = qs.Filter("user__username", "slene").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))
}

func TestSetCond(t *testing.T) {
	cond := NewCondition()
	cond1 := cond.And("profile__isnull", false).AndNot("status__in", 1).Or("profile__age__gt", 2000)

	qs := dORM.QueryTable("user")
	num, err := qs.SetCond(cond1).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	cond2 := cond.AndCond(cond1).OrCond(cond.And("user_name", "slene"))
	num, err = qs.SetCond(cond2).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))
}

func TestLimit(t *testing.T) {
	var posts []*Post
	qs := dORM.QueryTable("post")
	num, err := qs.Limit(1).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.Limit(-1).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 4))

	num, err = qs.Limit(-1, 2).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))

	num, err = qs.Limit(0, 2).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))
}

func TestOffset(t *testing.T) {
	var posts []*Post
	qs := dORM.QueryTable("post")
	num, err := qs.Limit(1).Offset(2).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.Offset(2).All(&posts)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))
}

func TestOrderBy(t *testing.T) {
	qs := dORM.QueryTable("user")
	num, err := qs.OrderBy("-status").Filter("user_name", "nobody").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.OrderBy("status").Filter("user_name", "slene").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	num, err = qs.OrderBy("-profile__age").Filter("user_name", "astaxie").Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))
}

func TestPrepareInsert(t *testing.T) {
	qs := dORM.QueryTable("user")
	i, err := qs.PrepareInsert()
	throwFail(t, err)

	var user User
	user.UserName = "testing1"
	num, err := i.Insert(&user)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Large, 0))

	user.UserName = "testing2"
	num, err = i.Insert(&user)
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Large, 0))

	num, err = qs.Filter("user_name__in", "testing1", "testing2").Delete()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 2))

	err = i.Close()
	throwFail(t, err)
	err = i.Close()
	throwFail(t, AssertIs(err, T_Equal, ErrStmtClosed))
}

func TestRaw(t *testing.T) {
	switch dORM.Driver().Type() {
	case DR_MySQL:
		num, err := dORM.Raw("UPDATE user SET user_name = ? WHERE user_name = ?", "testing", "slene").Exec()
		throwFail(t, err)
		throwFail(t, AssertIs(num, T_Equal, 1))

		num, err = dORM.Raw("UPDATE user SET user_name = ? WHERE user_name = ?", "slene", "testing").Exec()
		throwFail(t, err)
		throwFail(t, AssertIs(num, T_Equal, 1))

		var maps []Params
		num, err = dORM.Raw("SELECT user_name FROM user WHERE status = ?", 1).Values(&maps)
		throwFail(t, err)
		throwFail(t, AssertIs(num, T_Equal, 1))
		if num == 1 {
			throwFail(t, AssertIs(maps[0]["user_name"], T_Equal, "slene"))
		}

		var lists []ParamsList
		num, err = dORM.Raw("SELECT user_name FROM user WHERE status = ?", 1).ValuesList(&lists)
		throwFail(t, err)
		throwFail(t, AssertIs(num, T_Equal, 1))
		if num == 1 {
			throwFail(t, AssertIs(lists[0][0], T_Equal, "slene"))
		}

		var list ParamsList
		num, err = dORM.Raw("SELECT profile_id FROM user ORDER BY id ASC").ValuesFlat(&list)
		throwFail(t, err)
		throwFail(t, AssertIs(num, T_Equal, 3))
		if num == 3 {
			throwFail(t, AssertIs(list[0], T_Equal, "2"))
			throwFail(t, AssertIs(list[1], T_Equal, "3"))
			throwFail(t, AssertIs(list[2], T_Equal, ""))
		}
	}
}

func TestUpdate(t *testing.T) {
	qs := dORM.QueryTable("user")
	num, err := qs.Filter("user_name", "slene").Update(Params{
		"is_staff": true,
	})
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))
}

func TestDelete(t *testing.T) {
	qs := dORM.QueryTable("user_profile")
	num, err := qs.Filter("user__user_name", "slene").Delete()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

	qs = dORM.QueryTable("user")
	num, err = qs.Filter("user_name", "slene").Filter("profile__isnull", true).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))
}

func TestTransaction(t *testing.T) {
	o := NewOrm()
	err := o.Begin()
	throwFail(t, err)

	var names = []string{"1", "2", "3"}

	var user User
	user.UserName = names[0]
	id, err := o.Insert(&user)
	throwFail(t, err)
	throwFail(t, AssertIs(id, T_Large, 0))

	num, err := o.QueryTable("user").Filter("user_name", "slene").Update(Params{"user_name": names[1]})
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Large, 0))

	switch o.Driver().Type() {
	case DR_MySQL:
		id, err := o.Raw("INSERT INTO user (user_name) VALUES (?)", names[2]).Exec()
		throwFail(t, err)
		throwFail(t, AssertIs(id, T_Large, 0))
	}

	err = o.Rollback()
	throwFail(t, err)

	num, err = o.QueryTable("user").Filter("user_name__in", &user).Count()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 0))

	err = o.Begin()
	throwFail(t, err)

	user.UserName = "commit"
	id, err = o.Insert(&user)
	throwFail(t, err)
	throwFail(t, AssertIs(id, T_Large, 0))

	o.Commit()
	throwFail(t, err)

	num, err = o.QueryTable("user").Filter("user_name", "commit").Delete()
	throwFail(t, err)
	throwFail(t, AssertIs(num, T_Equal, 1))

}
