## 模型定义

复杂的模型定义不是必须的，此功能用作数据库数据转换和自动建表

## Struct Tag 设置参数
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
Name `orm:"column(user_name)"`
```
#### default

为字段设置默认值，类型必须符合
```go
type User struct {
	...
	Status int `orm:"default(1)"`
```
仅当进行 orm.Manager 初始化时才会赋值
```go
func NewUser() *User {
	obj := new(User)
	obj.Manager.Init(obj)
	return obj
}

u := NewUser()
fmt.Println(u.Status) // 1
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
## 表关系设置

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


## Struct Field 类型与数据库的对应

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