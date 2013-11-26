package orm

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// A slice string field.
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
	return TypeCharField
}

func (e *SliceStringField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case []string:
		e.Set(d)
	case string:
		if len(d) > 0 {
			parts := strings.Split(d, ",")
			v := make([]string, 0, len(parts))
			for _, p := range parts {
				v = append(v, strings.TrimSpace(p))
			}
			e.Set(v)
		}
	default:
		return fmt.Errorf("<SliceStringField.SetRaw> unknown value `%v`", value)
	}
	return nil
}

func (e *SliceStringField) RawValue() interface{} {
	return e.String()
}

var _ Fielder = new(SliceStringField)

// A json field.
type JsonField struct {
	Name string
	Data string
}

func (e *JsonField) String() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *JsonField) FieldType() int {
	return TypeTextField
}

func (e *JsonField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case string:
		return json.Unmarshal([]byte(d), e)
	default:
		return fmt.Errorf("<JsonField.SetRaw> unknown value `%v`", value)
	}
	return nil
}

func (e *JsonField) RawValue() interface{} {
	return e.String()
}

var _ Fielder = new(JsonField)

type Data struct {
	Id       int
	Boolean  bool
	Char     string    `orm:"size(50)"`
	Text     string    `orm:"type(text)"`
	Date     time.Time `orm:"type(date)"`
	DateTime time.Time `orm:"column(datetime)"`
	Byte     byte
	Rune     rune
	Int      int
	Int8     int8
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
	Decimal  float64 `orm:"digits(8);decimals(4)"`
}

type DataNull struct {
	Id       int
	Boolean  bool      `orm:"null"`
	Char     string    `orm:"null;size(50)"`
	Text     string    `orm:"null;type(text)"`
	Date     time.Time `orm:"null;type(date)"`
	DateTime time.Time `orm:"null;column(datetime)""`
	Byte     byte      `orm:"null"`
	Rune     rune      `orm:"null"`
	Int      int       `orm:"null"`
	Int8     int8      `orm:"null"`
	Int16    int16     `orm:"null"`
	Int32    int32     `orm:"null"`
	Int64    int64     `orm:"null"`
	Uint     uint      `orm:"null"`
	Uint8    uint8     `orm:"null"`
	Uint16   uint16    `orm:"null"`
	Uint32   uint32    `orm:"null"`
	Uint64   uint64    `orm:"null"`
	Float32  float32   `orm:"null"`
	Float64  float64   `orm:"null"`
	Decimal  float64   `orm:"digits(8);decimals(4);null"`
}

// only for mysql
type UserBig struct {
	Id   uint64
	Name string
}

type User struct {
	Id         int
	UserName   string `orm:"size(30);unique"`
	Email      string `orm:"size(100)"`
	Password   string `orm:"size(100)"`
	Status     int16
	IsStaff    bool
	IsActive   bool      `orm:"default(1)"`
	Created    time.Time `orm:"auto_now_add;type(date)"`
	Updated    time.Time `orm:"auto_now"`
	Profile    *Profile  `orm:"null;rel(one);on_delete(set_null)"`
	Posts      []*Post   `orm:"reverse(many)" json:"-"`
	ShouldSkip string    `orm:"-"`
	Nums       int
	Langs      SliceStringField `orm:"size(100)"`
	Extra      JsonField        `orm:"type(text)"`
}

func (u *User) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "UserName"},
		[]string{"Id", "Created"},
	}
}

func (u *User) TableUnique() [][]string {
	return [][]string{
		[]string{"UserName", "Email"},
	}
}

func NewUser() *User {
	obj := new(User)
	return obj
}

type Profile struct {
	Id       int
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
	Id      int
	User    *User     `orm:"rel(fk)"`
	Title   string    `orm:"size(60)"`
	Content string    `orm:"type(text)"`
	Created time.Time `orm:"auto_now_add"`
	Updated time.Time `orm:"auto_now"`
	Tags    []*Tag    `orm:"rel(m2m);rel_through(github.com/astaxie/beego/orm.PostTags)"`
}

func (u *Post) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "Created"},
	}
}

func NewPost() *Post {
	obj := new(Post)
	return obj
}

type Tag struct {
	Id       int
	Name     string  `orm:"size(30)"`
	BestPost *Post   `orm:"rel(one);null"`
	Posts    []*Post `orm:"reverse(many)" json:"-"`
}

func NewTag() *Tag {
	obj := new(Tag)
	return obj
}

type PostTags struct {
	Id   int
	Post *Post `orm:"rel(fk)"`
	Tag  *Tag  `orm:"rel(fk)"`
}

func (m *PostTags) TableName() string {
	return "prefix_post_tags"
}

type Comment struct {
	Id      int
	Post    *Post     `orm:"rel(fk);column(post)"`
	Content string    `orm:"type(text)"`
	Parent  *Comment  `orm:"null;rel(fk)"`
	Created time.Time `orm:"auto_now_add"`
}

func NewComment() *Comment {
	obj := new(Comment)
	return obj
}

var DBARGS = struct {
	Driver string
	Source string
	Debug  string
}{
	os.Getenv("ORM_DRIVER"),
	os.Getenv("ORM_SOURCE"),
	os.Getenv("ORM_DEBUG"),
}

var (
	IsMysql    = DBARGS.Driver == "mysql"
	IsSqlite   = DBARGS.Driver == "sqlite3"
	IsPostgres = DBARGS.Driver == "postgres"
)

var (
	dORM     Ormer
	dDbBaser dbBaser
)

func init() {
	Debug, _ = StrTo(DBARGS.Debug).Bool()

	if DBARGS.Driver == "" || DBARGS.Source == "" {
		fmt.Println(`need driver and source!

Default DB Drivers.

  driver: url
   mysql: https://github.com/go-sql-driver/mysql
 sqlite3: https://github.com/mattn/go-sqlite3
postgres: https://github.com/lib/pq

usage:

go get -u github.com/astaxie/beego/orm
go get -u github.com/go-sql-driver/mysql
go get -u github.com/mattn/go-sqlite3
go get -u github.com/lib/pq

#### MySQL
mysql -u root -e 'create database orm_test;'
export ORM_DRIVER=mysql
export ORM_SOURCE="root:@/orm_test?charset=utf8"
go test -v github.com/astaxie/beego/orm


#### Sqlite3
touch /path/to/orm_test.db
export ORM_DRIVER=sqlite3
export ORM_SOURCE=/path/to/orm_test.db
go test -v github.com/astaxie/beego/orm


#### PostgreSQL
psql -c 'create database orm_test;' -U postgres
export ORM_DRIVER=postgres
export ORM_SOURCE="user=postgres dbname=orm_test sslmode=disable"
go test -v github.com/astaxie/beego/orm
`)
		os.Exit(2)
	}

	RegisterDataBase("default", DBARGS.Driver, DBARGS.Source, 20)

	alias := getDbAlias("default")
	if alias.Driver == DR_MySQL {
		alias.Engine = "INNODB"
	}

}
