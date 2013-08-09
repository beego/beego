package orm

import (
	"fmt"
	"os"
	"time"

	_ "github.com/bmizerany/pq"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id         int    `orm:"auto"`
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
}

func NewUser() *User {
	obj := new(User)
	return obj
}

type Profile struct {
	Id    int     `orm:"auto"`
	Age   int16   ``
	Money float64 ``
	User  *User   `orm:"reverse(one)" json:"-"`
}

func (u *Profile) TableName() string {
	return "user_profile"
}

func NewProfile() *Profile {
	obj := new(Profile)
	return obj
}

type Post struct {
	Id      int       `orm:"auto"`
	User    *User     `orm:"rel(fk)"` //
	Title   string    `orm:"size(60)"`
	Content string    ``
	Created time.Time `orm:"auto_now_add"`
	Updated time.Time `orm:"auto_now"`
	Tags    []*Tag    `orm:"rel(m2m)"`
}

func NewPost() *Post {
	obj := new(Post)
	return obj
}

type Tag struct {
	Id    int     `orm:"auto"`
	Name  string  `orm:"size(30)"`
	Posts []*Post `orm:"reverse(many)" json:"-"`
}

func NewTag() *Tag {
	obj := new(Tag)
	return obj
}

type Comment struct {
	Id      int       `orm:"auto"`
	Post    *Post     `orm:"rel(fk)"`
	Content string    ``
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

var dORM Ormer

func init() {
	RegisterModel(new(User))
	RegisterModel(new(Profile))
	RegisterModel(new(Post))
	RegisterModel(new(Tag))
	RegisterModel(new(Comment))

	Debug, _ = StrTo(DBARGS.Debug).Bool()

	if DBARGS.Driver == "" || DBARGS.Source == "" {
		fmt.Println(`need driver and source!

Default DB Drivers.

  driver: url
   mysql: https://github.com/go-sql-driver/mysql
 sqlite3: https://github.com/mattn/go-sqlite3
postgres: https://github.com/bmizerany/pq

eg: mysql
ORM_DRIVER=mysql ORM_SOURCE="root:root@/my_db?charset=utf8" go test github.com/astaxie/beego/orm
`)
		os.Exit(2)
	}

	RegisterDataBase("default", DBARGS.Driver, DBARGS.Source, 20)

	BootStrap()

	truncateTables()

	dORM = NewOrm()
}

func truncateTables() {
	logs := "truncate tables for test\n"
	o := NewOrm()
	for _, m := range modelCache.allOrdered() {
		query := fmt.Sprintf("truncate table `%s`", m.table)
		_, err := o.Raw(query).Exec()
		logs += query + "\n"
		if err != nil {
			fmt.Println(logs)
			fmt.Println(err)
			os.Exit(2)
		}
	}
}
