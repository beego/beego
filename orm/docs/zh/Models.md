## Model Definition

比较全面的 Model 定义例子，后文所有的例子如无特殊说明都以这个为基础。

当前还没有完成自动创建表的功能，所以提供一个 [Models.sql](Models.sql) 测试

note: 根据文档的更新，随时都可能更新这个 Model

##### models.go:

```go
package main

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id          int        `orm:"auto"`            // 设置为auto主键
	UserName    string     `orm:"size(30);unique"` // 设置字段为unique
	Email       string     `orm:"size(100)"`       // 设置string字段长度时,会使用varchar类型
	Password    string     `orm:"size(100)"`
	Status      int16      `orm:"choices(0,1,2,3);defalut(0)"` // choices设置可选值
	IsStaff     bool       `orm:"default(false)"`              // default设置默认值
	IsActive    bool       `orm:"default(0)"`
	Created     time.Time  `orm:"auto_now_add;type(date)"`           // 创建时自动设置时间
	Updated     time.Time  `orm:"auto_now"`                          // 每次更新时自动设置时间
	Profile     *Profile   `orm:"null;rel(one);on_delete(set_null)"` // OneToOne relation, 级联删除时设置为NULL
	Posts       []*Post `orm:"reverse(many)" json:"-"` // fk 的反向关系
	orm.Manager `json:"-"` // 每个model都需要定义orm.Manager
}

// 定义NewModel进行orm.Manager的初始化(必须)
func NewUser() *User {
	obj := new(User)
	obj.Manager.Init(obj)
	return obj
}

type Profile struct {
	Id          int     `orm:"auto"`
	Age         int16   ``
	Money       float64 ``
	User        *User   `orm:"reverse(one)" json:"-"` // 设置反向关系(字段可选)
	orm.Manager `json:"-"`
}

func (u *Profile) TableName() string {
	return "profile" // 自定义表名
}

func NewProfile() *Profile {
	obj := new(Profile)
	obj.Manager.Init(obj)
	return obj
}

type Post struct {
	Id          int       `orm:"auto"`
	User        *User     `orm:"rel(fk)"` // RelForeignKey relation
	Title       string    `orm:"size(60)"`
	Content     string    ``
	Created     time.Time ``
	Updated     time.Time ``
	Tags        []*Tag    `orm:"rel(m2m)"` // ManyToMany relation
	orm.Manager `json:"-"`
}

func NewPost() *Post {
	obj := new(Post)
	obj.Manager.Init(obj)
	return obj
}

type Tag struct {
	Id          int     `orm:"auto"`
	Name        string  `orm:"size(30)"`
	Status      int16   `orm:"choices(0,1,2);default(0)"`
	Posts       []*Post `orm:"reverse(many)" json:"-"`
	orm.Manager `json:"-"`
}

func NewTag() *Tag {
	obj := new(Tag)
	obj.Manager.Init(obj)
	return obj
}

type Comment struct {
	Id          int       `orm:"auto"`
	Post        *Post     `orm:"rel(fk)"`
	Content     string    ``
	Parent      *Comment  `orm:"null;rel(fk)"` // null设置allow NULL
	Status      int16     `orm:"choices(0,1,2);default(0)"`
	Created     time.Time `orm:"auto_now_add"`
	orm.Manager `json:"-"`
}

func NewComment() *Comment {
	obj := new(Comment)
	obj.Manager.Init(obj)
	return obj
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Profile))
	orm.RegisterModel(new(Post))
	orm.RegisterModel(new(Tag))
	orm.RegisterModel(new(Comment))
}
```

## Field Type

现在 orm 支持下面的字段形式

| go type		   | field type  | mysql type
| :---   	   | :---        | :---
| bool | TypeBooleanField | tinyint
| string | TypeCharField | varchar
| string | TypeTextField | longtext
| time.Time | TypeDateField | date
| time.TIme | TypeDateTimeField | datetime
|  int16 |TypeSmallIntegerField | int(4)
|  int, int32 |TypeIntegerField | int(11)
|  int64 |TypeBigIntegerField | bigint(20)
|  uint, uint16 |TypePositiveSmallIntegerField | int(4) unsigned
|  uint32 |TypePositiveIntegerField | int(11) unsigned
|  uint64 |TypePositiveBigIntegerField | bigint(20) unsigned
| float32, float64 | TypeFloatField | double
| float32, float64 | TypeDecimalField | double(digits, decimals)

关系型的字段，其字段类型取决于对应的主键。

* RelForeignKey
* RelOneToOne
* RelManyToMany
* RelReverseOne
* RelReverseMany

## Field Options
```go
orm:"null;rel(fk)"
```
	
通常每个 Field 的 StructTag 里包含两种类型的设置，类似 null 的 bool 型设置，还有 类似 rel(fk) 的指定值设置，bool 型默认为 false，指定以后即表示为 true

多个设置间使用 `;` 分隔，设置的值如果是多个，使用 `,` 分隔。

#### auto

设置为 Autoincrement Primary Key

#### pk

设置为 Primary Key

#### null

数据库表默认为 `NOT NULL`，设置 null 代表 `ALLOW NULL`

#### blank

设置 string 类型的字段允许为空，否则 clean 会返回错误

#### index

为字段增加索引

#### unique

为字段增加 unique 键

#### column

为字段设置 db 字段的名称
```go
UserName `orm:"column(db_user_name)"`
```
#### default

为字段设置默认值，类型必须符合
```go
Status int `orm:"default(1)"`
```
#### choices

为字段设置一组可选的值，类型必须符合。其他值 clean 会返回错误
```go
Status int `orm:"choices(1,2,3,4)"`
```
#### size (string)

string 类型字段设置 size 以后，db type 将使用 varchar
```go
Title string `orm:"size(60)"`
```
#### digits / decimals

设置 float32, float64 类型的浮点精度
```go
Money float64 `orm:"digits(12);decimals(4)"`
```
总长度 12 小数点后 4 位 eg: `99999999.9999`

#### auto_now / auto_now_add
```go
Created     time.Time `auto_now_add`
Updated     time.Time `auto_now`
```
* auto_now 每次 model 保存时都会对时间自动更新
* auto_now_add 第一次保存时才设置时间

对于批量的 update 此设置是不生效的

#### type

设置为 date, time.Time 字段的对应 db 类型使用 date
```go
Created time.Time `orm:"auto_now_add;type(date)"`
```
## Relation Field Options

#### rel / reverse

**RelOneToOne**:
```go
type User struct {
	...
	Profile *Profile `orm:"null;rel(one);on_delete(set_null)"`
```
对应的反向关系 **RelReverseOne**:
```go
type Profile struct {
	...
	User *User `orm:"reverse(one)" json:"-"`
```
**RelForeignKey**:
```go
type Post struct {
	...
	User*User `orm:"rel(fk)"` // RelForeignKey relation
```
对应的反向关系 **RelReverseMany**:
```go
type User struct {
	...
	Posts []*Post `orm:"reverse(many)" json:"-"` // fk 的反向关系
```
**RelManyToMany**:
```go
type Post struct {
	...
	Tags []*Tag `orm:"rel(m2m)"` // ManyToMany relation
```
对应的反向关系 **RelReverseMany**:
```go
type Tag struct {
	...
	Posts []*Post `orm:"reverse(many)" json:"-"`
```
#### rel_table / rel_through

此设置针对 `orm:"rel(m2m)"` 的关系字段

	rel_table       设置自动生成的 m2m 关系表的名称
	rel_through     如果要在 m2m 关系中使用自定义的 m2m 关系表
	                通过这个设置其名称，格式为 pkg.path.ModelName
	                eg: app.models.PostTagRel
	                PostTagRel 表需要有到 Post 和 Tag 的关系

当设置 rel_table 时会忽略 rel_through

#### on_delete

设置对应的 rel 关系删除时，如何处理关系字段。

	cascade        级联删除(默认值)
	set_null       设置为 NULL，需要设置 null = true
	set_default    设置为默认值，需要设置 default 值
	do_nothing     什么也不做，忽略

```go
type User struct {
	...
	Profile *Profile `orm:"null;rel(one);on_delete(set_null)"`
...
type Profile struct {
	...
	User *User `orm:"reverse(one)" json:"-"`

// 删除 Profile 时将设置 User.Profile 的数据库字段为 NULL
```