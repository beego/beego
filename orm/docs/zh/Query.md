## Query

orm 以 **QuerySeter** 来组织查询，每个返回 **QuerySeter** 的方法都会获得一个新的 **QuerySeter** 对象。

基本使用方法:
```go
o := orm.NewOrm()

// 获取 QuerySeter 对象，user 为表名
qs := o.QueryTable("user")

// 也可以直接使用对象作为表名
user := NewUser()
qs = o.QueryTable(user) // 返回 QuerySeter
```
## expr

QuerySeter 中用于描述字段和 sql 操作符使用简单的 expr 查询方法

字段组合的前后顺序依照表的关系，比如 User 表拥有 Profile 的外键，那么对 User 表查询对应的 Profile.Age 为条件，则使用 `Profile__Age` 注意，字段的分隔符号使用双下划线 `__`，除了描述字段， expr 的尾部可以增加操作符以执行对应的 sql 操作。比如 `Profile__Age__gt` 代表 Profile.Age > 18 的条件查询。

注释后面将描述对应的 sql 语句，仅仅是描述 expr 的类似结果，并不代表实际生成的语句。
```go
qs.Filter("id", 1) // WHERE id = 1
qs.Filter("profile__age", 18) // WHERE profile.age = 18
qs.Filter("Profile__Age", 18) // 使用字段名和Field名都是允许的
qs.Filter("profile__age", 18) // WHERE profile.age = 18
qs.Filter("profile__age__gt", 18) // WHERE profile.age > 18
qs.Filter("profile__age__gte", 18) // WHERE profile.age >= 18
qs.Filter("profile__age__in", 18, 20) // WHERE profile.age IN (18, 20)

qs.Filter("profile__age__in", 18, 20).Exclude("profile__money__lt", 1000)
// WHERE profile.age IN (18, 20) AND NOT profile.money < 1000
```
## Operators

当前支持的操作符号：

* [exact](#exact) / [iexact](#iexact) 等于
* [contains](#contains) / [icontains](#icontains) 包含
* [gt / gte](#gt / gte) 大于 / 大于等于
* [lt / lte](#lt / lte) 小于 / 小于等于
* [startswith](#startswith) / [istartswith](#istartswith) 以...起始
* [endswith](#endswith) / [iendswith](#iendswith) 以...结束
* [in](#in)
* [isnull](#isnull)

后面以 `i` 开头的表示：大小写不敏感

#### exact

Filter / Exclude / Condition expr 的默认值
```go
qs.Filter("user_name", "slene") // WHERE user_name = 'slene'
qs.Filter("user_name__exact", "slene") // WHERE user_name = 'slene'
// 使用 = 匹配，大小写是否敏感取决于数据表使用的 collation
qs.Filter("profile", nil) // WHERE profile_id IS NULL
```
#### iexact
```go
qs.Filter("user_name__iexact", "slene")
// WHERE user_name LIKE 'slene'
// 大小写不敏感，匹配任意 'Slene' 'sLENE'
```
#### contains
```go
qs.Filter("user_name__contains", "slene")
// WHERE user_name LIKE BINARY '%slene%'
// 大小写敏感, 匹配包含 slene 的字符
```
#### icontains
```go
qs.Filter("user_name__icontains", "slene")
// WHERE user_name LIKE '%slene%'
// 大小写不敏感, 匹配任意 'im Slene', 'im sLENE'
```
#### in
```go
qs.Filter("profile__age__in", 17, 18, 19, 20)
// WHERE profile.age IN (17, 18, 19, 20)
```
#### gt / gte
```go
qs.Filter("profile__age__gt", 17)
// WHERE profile.age > 17

qs.Filter("profile__age__gte", 18)
// WHERE profile.age >= 18
```
#### lt / lte
```go
qs.Filter("profile__age__lt", 17)
// WHERE profile.age < 17

qs.Filter("profile__age__lte", 18)
// WHERE profile.age <= 18
```
#### startswith
```go
qs.Filter("user_name__startswith", "slene")
// WHERE user_name LIKE BINARY 'slene%'
// 大小写敏感, 匹配以 'slene' 起始的字符串
```
#### istartswith
```go
qs.Filter("user_name__istartswith", "slene")
// WHERE user_name LIKE 'slene%'
// 大小写不敏感, 匹配任意以 'slene', 'Slene' 起始的字符串
```
#### endswith
```go
qs.Filter("user_name__endswith", "slene")
// WHERE user_name LIKE BINARY '%slene'
// 大小写敏感, 匹配以 'slene' 结束的字符串
```
#### iendswith
```go
qs.Filter("user_name__startswith", "slene")
// WHERE user_name LIKE '%slene'
// 大小写不敏感, 匹配任意以 'slene', 'Slene' 结束的字符串
```
#### isnull
```go
qs.Filter("profile__isnull", true)
qs.Filter("profile_id__isnull", true)
// WHERE profile_id IS NULL

qs.Filter("profile__isnull", false)
// WHERE profile_id IS NOT NULL
```
## QuerySeter

QuerySeter 当前支持的方法

* type QuerySeter interface {
	* [Filter(string, ...interface{}) QuerySeter](#Filter)
	* [Exclude(string, ...interface{}) QuerySeter](#Exclude)
	* [SetCond(*Condition) QuerySeter](#SetCond)
	* [Limit(int, ...int64) QuerySeter](#Limit)
	* [Offset(int64) QuerySeter](#Offset)
	* [OrderBy(...string) QuerySeter](#OrderBy)
	* [RelatedSel(...interface{}) QuerySeter](#RelatedSel)
	* [Count() (int64, error)](#Count)
	* [Update(Params) (int64, error)](#Update)
	* [Delete() (int64, error)](#Delete)
	* [PrepareInsert() (Inserter, error)](#PrepareInsert)
	* [All(interface{}) (int64, error)](#All)
	* [One(Modeler) error](#One)
	* [Values(*[]Params, ...string) (int64, error)](#Values)
	* [ValuesList(*[]ParamsList, ...string) (int64, error)](#ValuesList)
	* [ValuesFlat(*ParamsList, string) (int64, error)](#ValuesFlat)
* }

每个返回 QuerySeter 的 api 调用时都会新建一个 QuerySeter，不影响之前创建的。

#### Filter

多个 Filter 之间使用 `AND` 连接
```go
qs.Filter("profile__isnull", true).Filter("user_name", "slene")
// WHERE profile_id IS NULL AND user_name = 'slene'
```

#### Exclude

使用 `NOT` 排除条件

多个 Exclude 之间使用 `AND` 连接
```go
qs.Exclude("profile__isnull", true).Filter("user_name", "slene")
// WHERE NOT profile_id IS NULL AND user_name = 'slene'
```

#### SetCond

自定义条件表达式

```go
cond := NewCondition()
cond1 := cond.And("profile__isnull", false).AndNot("status__in", 1).Or("profile__age__gt", 2000)

qs := orm.QueryTable("user")
qs = qs.SetCond(cond1)
// WHERE ... AND ... AND NOT ... OR ...

cond2 := cond.AndCond(cond1).OrCond(cond.And("user_name", "slene"))
qs = qs.SetCond(cond2).Count()
// WHERE (... AND ... AND NOT ... OR ...) OR ( ... )
```

#### Limit

限制最大返回数据行数，第二个参数可以设置 `Offset`
```go
var DefaultRowsLimit = 1000 // orm 默认的 limit 值为 1000

// 默认情况下 select 查询的最大行数为 1000
// LIMIT 1000

qs.Limit(10)
// LIMIT 10

qs.Limit(10, 20)
// LIMIT 10 OFFSET 20

qs.Limit(-1)
// no limit

qs.Limit(-1, 100)
// LIMIT 18446744073709551615 OFFSET 100
// 18446744073709551615 是 1<<64 - 1 用来指定无 limit 限制 但有 offset 偏移的情况
```

#### Offset
	
设置 偏移行数
```go
qs.Offset(20)
// LIMIT 1000 OFFSET 20
```

#### OrderBy

参数使用 **expr**

在 expr 前使用减号 `-` 表示 `DESC` 的排列
```go
qs.OrderBy("id", "-profile__age")
// ORDER BY id ASC, profile.age DESC

qs.OrderBy("-profile__money", "profile")
// ORDER BY profile.money DESC, profile_id ASC
```

#### RelatedSel

关系查询，参数使用 **expr**
```go
var DefaultRelsDepth = 5 // 默认情况下直接调用 RelatedSel 将进行最大 5 层的关系查询

qs := o.QueryTable("post")

qs.RelateSel()
// INNER JOIN user ... LEFT OUTER JOIN profile ...

qs.RelateSel("user")
// INNER JOIN user ... 
// 设置 expr 只对设置的字段进行关系查询

// 对设置 null 属性的 Field 将使用 LEFT OUTER JOIN
```

#### Count
依据当前的查询条件，返回结果行数
```go
cnt, err := o.QueryTable("user").Count() // SELECT COUNT(*) FROM USER
fmt.Printf("Count Num: %s, %s", cnt, err)
```

#### Update
依据当前查询条件，进行批量更新操作
```go
num, err := o.QueryTable("user").Filter("user_name", "slene").Update(orm.Params{
	"user_name": "astaxie",
})
fmt.Printf("Affected Num: %s, %s", num, err)
// SET user_name = "astaixe" WHERE user_name = "slene"
```

#### Delete
依据当前查询条件，进行批量删除操作
```go
num, err := o.QueryTable("user").Filter("user_name", "slene").Delete()
fmt.Printf("Affected Num: %s, %s", num, err)
// DELETE FROM user WHERE user_name = "slene"
```

#### PrepareInsert

用于一次 prepare 多次 insert 插入，以提高批量插入的速度。

```go
var users []*User
...
qs := dORM.QueryTable("user")
i, _ := qs.PrepareInsert()
for _, user := range users {
	id, err := i.Insert(user)
	if err != nil {
		...
	}
}
// PREPARE INSERT INTO user (`user_name`, ...) VALUES (?, ...)
// EXECUTE INSERT INTO user (`user_name`, ...) VALUES ("slene", ...)
// EXECUTE ...
// ...
i.Close() // 别忘记关闭 statement
```

#### All
返回对应的结果集对象
```go
var users []*User
num, err := o.QueryTable("user").Filter("user_name", "slene").All(&users)
fmt.Printf("Returned Rows Num: %s, %s", num, err)
```

#### One

尝试返回单条记录

```go
var user *User
err := o.QueryTable("user").Filter("user_name", "slene").One(&user)
if err == orm.ErrMultiRows {
	// 多条的时候报错
	fmt.Printf("Returned Multi Rows Not One")
}
if err == orm.ErrNoRows {
	// 没有找到记录
	fmt.Printf("Not row found")
}
```

#### Values
返回结果集的 key => value 值

key 为 Model 里的 Field name，value 的值 以 string 保存

```go
var maps []orm.Params
num, err := o.QueryTable("user").Values(&maps)
if err != nil {
	fmt.Printf("Result Nums: %d\n", num)
	for _, m := range maps {
		fmt.Println(m["Id"], m["UserName"])
	}
}
```

返回指定的 Field 数据

**TODO**: 暂不支持级联查询 **RelatedSel** 直接返回 Values

但可以直接指定 expr 级联返回需要的数据

```go
var maps []orm.Params
num, err := o.QueryTable("user").Values(&maps, "id", "user_name", "profile", "profile__age")
if err != nil {
	fmt.Printf("Result Nums: %d\n", num)
	for _, m := range maps {
		fmt.Println(m["Id"], m["UserName"], m["Profile"], m["Profile__Age"])
		// map 中的数据都是展开的，没有复杂的嵌套
	}
}
```

#### ValuesList

顾名思义，返回的结果集以slice存储

结果的排列与 Model 中定义的 Field 顺序一致

返回的每个元素值以 string 保存

```go
var lists []orm.ParamsList
num, err := o.QueryTable("user").ValuesList(&lists)
if err != nil {
	fmt.Printf("Result Nums: %d\n", num)
	for _, row := range lists {
		fmt.Println(row)
	}
}
```

当然也可以指定 expr 返回指定的 Field

```go
var lists []orm.ParamsList
num, err := o.QueryTable("user").ValuesList(&lists, "user_name", "profile__age")
if err != nil {
	fmt.Printf("Result Nums: %d\n", num)
	for _, row := range lists {
		fmt.Printf("UserName: %s, Age: %s\m", row[0], row[1])
	}
}
```

#### ValuesFlat

只返回特定的 Field 值，讲结果集展开到单个 slice 里

```go
var list orm.ParamsList
num, err := o.QueryTable("user").ValuesFlat(&list, "user_name")
if err != nil {
	fmt.Printf("Result Nums: %d\n", num)
	fmt.Printf("All User Names: %s", strings.Join(list, ", ")
}
```



