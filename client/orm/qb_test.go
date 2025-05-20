package orm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm/internal/models"
	"github.com/beego/beego/v2/client/orm/internal/session"
	"github.com/beego/beego/v2/client/orm/internal/utils"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

var (
	dORM session.Ormer
)

var helpinfo = `need driver and source!

	Default DB Drivers.

	  driver: url
	   mysql: https://github.com/go-sql-driver/mysql
	 sqlite3: https://github.com/mattn/go-sqlite3
	postgres: https://github.com/lib/pq
	tidb: https://github.com/pingcap/tidb

	usage:

	go Get -u github.com/beego/beego/v2/client/orm
	go Get -u github.com/go-sql-driver/mysql
	go Get -u github.com/mattn/go-sqlite3
	go Get -u github.com/lib/pq
	go Get -u github.com/pingcap/tidb

	#### MySQL
	mysql -u root -e 'create database orm_test;'
	export ORM_DRIVER=mysql
	export ORM_SOURCE="root:@/orm_test?charset=utf8"
	go test -v github.com/beego/beego/v2/client/orm


	#### Sqlite3
	export ORM_DRIVER=sqlite3
	export ORM_SOURCE='file:memory_test?mode=memory'
	go test -v github.com/beego/beego/v2/client/orm


	#### PostgreSQL
	psql -c 'create database orm_test;' -U postgres
	export ORM_DRIVER=postgres
	export ORM_SOURCE="user=postgres dbname=orm_test sslmode=disable"
	go test -v github.com/beego/beego/v2/client/orm

	#### TiDB
	export ORM_DRIVER=tidb
	export ORM_SOURCE='memory://test/test'
	go test -v github.com/beego/beego/v2/pgk/orm

	`

var DBARGS = struct {
	Driver string
	Source string
	Debug  string
}{
	os.Getenv("ORM_DRIVER"),
	os.Getenv("ORM_SOURCE"),
	os.Getenv("ORM_DEBUG"),
}

func init() {

	if DBARGS.Driver == "" || DBARGS.Source == "" {
		fmt.Println(helpinfo)
		os.Exit(2)
	}

	err := session.RegisterDataBase("default", DBARGS.Driver, DBARGS.Source, session.MaxIdleConnections(20))
	if err != nil {
		panic(fmt.Sprintf("can not Register database: %v", err))
	}

}

func registerAllModel() {
	session.RegisterModel(new(User))
	session.RegisterModel(new(Profile))
	session.RegisterModel(new(Post))
	session.RegisterModel(new(NullValue))
	session.RegisterModel(new(Tag))
}

func TestRegisterModels(_ *testing.T) {
	models.DefaultModelCache.Clean()
	registerAllModel()

	session.BootStrap()

	dORM = session.NewOrm()
}

type User struct {
	ID             int    `orm:"column(id)"`
	UserName       string `orm:"size(30);unique"`
	Email          string `orm:"size(100)"`
	Password       string `orm:"size(100)"`
	Status         int16  `orm:"column(Status)"`
	IsStaff        bool
	IsActive       bool `orm:"default(true)"`
	Unexported     bool `orm:"-"`
	UnexportedBool bool
	Created        time.Time `orm:"auto_now_add;type(date)"`
	Updated        time.Time `orm:"auto_now"`
	Profile        *Profile  `orm:"null;rel(one);on_delete(set_null)"`
	Posts          []*Post   `orm:"reverse(many)" json:"-"`
	ShouldSkip     string    `orm:"-"`
	Nums           int
	Langs          SliceStringField `orm:"size(100)"`
	Extra          JSONFieldTest    `orm:"type(text)"`
}

func (u *User) TableIndex() [][]string {
	return [][]string{
		{"Id", "UserName"},
		{"Id", "Created"},
	}
}

func (u *User) TableUnique() [][]string {
	return [][]string{
		{"UserName", "Email"},
	}
}

func NewUser() *User {
	obj := new(User)
	return obj
}

// user_profile table
type userProfile struct {
	User
	Age   int
	Money float64
}

func TestPSQueryBuilder(t *testing.T) {
	// only test postgres
	if dORM.Driver().Type() != 4 {
		return
	}

	var user User
	var l []userProfile
	o := session.NewOrm()

	qb, err := NewQueryBuilder("postgres")
	if err != nil {
		throwFailNow(t, err)
	}
	qb.Select("user.id", "user.user_name").
		From("user").Where("id = ?").OrderBy("user_name").
		Desc().Limit(1).Offset(0)
	sql := qb.String()
	err = o.Raw(sql, 2).QueryRow(&user)
	if err != nil {
		throwFailNow(t, err)
	}
	throwFail(t, AssertIs(user.UserName, "slene"))

	qb.Select("*").
		From("user_profile").InnerJoin("user").
		On("user_profile.id = user.id")
	sql = qb.String()
	num, err := o.Raw(sql).QueryRows(&l)
	if err != nil {
		throwFailNow(t, err)
	}
	throwFailNow(t, AssertIs(num, 1))
	throwFailNow(t, AssertIs(l[0].UserName, "astaxie"))
	throwFailNow(t, AssertIs(l[0].Age, 30))
}

// deprecated using assert.XXX
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

func getCaller(skip int) string {
	pc, file, line, _ := runtime.Caller(skip)
	fun := runtime.FuncForPC(pc)
	_, fn := filepath.Split(file)
	data, err := os.ReadFile(file)
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
			ls := formatLines(string(lines[o+i]))
			code := fmt.Sprintf(" %s %5d:   %s", flag, cur, ls)
			if code != "" {
				codes = append(codes, code)
			}
		}
	}
	funName := fun.Name()
	if i := strings.LastIndex(funName, "."); i > -1 {
		funName = funName[i+1:]
	}
	return fmt.Sprintf("%s:%s:%d: \n%s", fn, funName, line, strings.Join(codes, "\n"))
}

func formatLines(s string) string {
	return strings.ReplaceAll(s, "\t", "    ")
}

// Deprecated: Using stretchr/testify/assert
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

type argAny []interface{}

// Get interface by index from interface slice
func (a argAny) Get(i int, args ...interface{}) (r interface{}) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

func ValuesCompare(is bool, a interface{}, args ...interface{}) (ok bool, err error) {
	if len(args) == 0 {
		return false, fmt.Errorf("miss args")
	}
	b := args[0]
	arg := argAny(args)

	switch v := a.(type) {
	case reflect.Kind:
		ok = reflect.ValueOf(b).Kind() == v
	case time.Time:
		if v2, vo := b.(time.Time); vo {
			if arg.Get(1) != nil {
				format := utils.ToStr(arg.Get(1))
				a = v.Format(format)
				b = v2.Format(format)
				ok = a == b
			} else {
				err = fmt.Errorf("compare datetime miss format")
				goto wrongArg
			}
		}
	default:
		ok = utils.ToStr(a) == utils.ToStr(b)
	}
	ok = is && ok || !is && !ok
	if !ok {
		err = fmt.Errorf("expected: `%v`, Get `%v`", b, a)
	}

wrongArg:
	if err != nil {
		return false, err
	}

	return true, nil
}

func AssertIs(a interface{}, args ...interface{}) error {
	if ok, err := ValuesCompare(true, a, args...); !ok {
		return err
	}
	return nil
}

type Profile struct {
	ID       int `orm:"column(id)"`
	Age      int16
	Money    float64
	User     *User `orm:"reverse(one)" json:"-"`
	BestPost *Post `orm:"rel(one);null"`
}

func (u *Profile) TableName() string {
	return "user_profile"
}

func NewProfile() *Profile {
	obj := new(Profile)
	return obj
}

type Post struct {
	ID               int       `orm:"column(id)"`
	User             *User     `orm:"rel(fk)"`
	Title            string    `orm:"size(60)"`
	Content          string    `orm:"type(text)"`
	Created          time.Time `orm:"auto_now_add"`
	Updated          time.Time `orm:"auto_now"`
	UpdatedPrecision time.Time `orm:"auto_now;type(datetime);precision(4)"`
	Tags             []*Tag    `orm:"rel(m2m);rel_through(github.com/beego/beego/v2/client/orm.PostTags)"`
}

func (u *Post) TableIndex() [][]string {
	return [][]string{
		{"Id", "Created"},
	}
}

func NewPost() *Post {
	obj := new(Post)
	return obj
}

type NullValue struct {
	ID    int    `orm:"column(id)"`
	Value string `orm:"size(30);null"`
}

type Tag struct {
	ID       int     `orm:"column(id)"`
	Name     string  `orm:"size(30)"`
	BestPost *Post   `orm:"rel(one);null"`
	Posts    []*Post `orm:"reverse(many)" json:"-"`
}

func NewTag() *Tag {
	obj := new(Tag)
	return obj
}

type SliceStringField []string

func (e SliceStringField) Value() []string {
	return []string(e)
}

func (e *SliceStringField) Set(d []string) {
	*e = SliceStringField(d)
}

func (e *SliceStringField) Add(v string) {
	*e = append(*e, v)
}

func (e *SliceStringField) String() string {
	return strings.Join(e.Value(), ",")
}

func (e *SliceStringField) FieldType() int {
	return models.TypeVarCharField
}

func (e *SliceStringField) SetRaw(value interface{}) error {
	f := func(str string) {
		if len(str) > 0 {
			parts := strings.Split(str, ",")
			v := make([]string, 0, len(parts))
			for _, p := range parts {
				v = append(v, strings.TrimSpace(p))
			}
			e.Set(v)
		}
	}

	switch d := value.(type) {
	case []string:
		e.Set(d)
	case string:
		f(d)
	case []byte:
		f(string(d))
	default:
		return fmt.Errorf("<SliceStringField.SetRaw> unknown value `%v`", value)
	}
	return nil
}

func (e *SliceStringField) RawValue() interface{} {
	return e.String()
}

var _ models.Fielder = new(SliceStringField)

// A json field.
type JSONFieldTest struct {
	Name string
	Data string
}

func (e *JSONFieldTest) String() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *JSONFieldTest) FieldType() int {
	return models.TypeTextField
}

func (e *JSONFieldTest) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case string:
		return json.Unmarshal([]byte(d), e)
	case []byte:
		return json.Unmarshal(d, e)
	default:
		return fmt.Errorf("<JSONField.SetRaw> unknown value `%v`", value)
	}
}

func (e *JSONFieldTest) RawValue() interface{} {
	return e.String()
}

var _ models.Fielder = new(JSONFieldTest)
