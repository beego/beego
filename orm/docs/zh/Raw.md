## 使用SQL语句进行查询

* 使用 Raw SQL 查询，无需使用 ORM 表定义
* 多数据库，都可直接使用占位符号 `?`，自动转换
* 查询时的参数，支持使用 Model Struct 和 Slice, Array

```go
ids := []int{1, 2, 3}
p.Raw("SELECT name FROM user WHERE id IN (?, ?, ?)", ids)
```

创建一个 **RawSeter**

```go
o := NewOrm()
var r RawSeter
r = o.Raw("UPDATE user SET name = ? WHERE name = ?", "testing", "slene")
```

* type RawSeter interface {
	* [Exec() (int64, error)](#exec)
	* [QueryRow(...interface{}) error](#queryrow)
	* [QueryRows(...interface{}) (int64, error)](#queryrows)
	* [SetArgs(...interface{}) RawSeter](#setargs)
	* [Values(*[]Params) (int64, error)](#values)
	* [ValuesList(*[]ParamsList) (int64, error)](#valueslist)
	* [ValuesFlat(*ParamsList) (int64, error)](#valuesflat)
	* [Prepare() (RawPreparer, error)](#prepare)
* }

#### Exec

执行sql语句

```go
num, err := r.Exec()
```

#### QueryRow

TODO

#### QueryRows

TODO

#### SetArgs

改变 Raw(sql, args...) 中的 args 参数，返回一个新的 RawSeter

用于单条 sql 语句，重复利用，替换参数然后执行。

```go
num, err := r.SetArgs("arg1", "arg2").Exec()
num, err := r.SetArgs("arg1", "arg2").Exec()
...
```
#### Values / ValuesList / ValuesFlat

Raw SQL 查询获得的结果集 Value 为 `string` 类型，NULL 字段的值为空 ``

#### Values


返回结果集的 key => value 值

```go
var maps []orm.Params
num, err = o.Raw("SELECT user_name FROM user WHERE status = ?", 1).Values(&maps)
if err == nil && num > 0 {
	fmt.Println(maps[0]["user_name"]) // slene
}
```

#### ValuesList

返回结果集 slice

```go
var lists []orm.ParamsList
num, err = o.Raw("SELECT user_name FROM user WHERE status = ?", 1).ValuesList(&lists)
if err == nil && num > 0 {
	fmt.Println(lists[0][0]) // slene
}
```

#### ValuesFlat

返回单一字段的平铺 slice 数据

```go
var list orm.ParamsList
num, err = o.Raw("SELECT id FROM user WHERE id < ?", 10).ValuesList(&list)
if err == nil && num > 0 {
	fmt.Println(list) // []{"1","2","3",...}
}
```

#### Prepare

用于一次 prepare 多次 exec，以提高批量执行的速度。

```go
p, err := o.Raw("UPDATE user SET name = ? WHERE name = ?").Prepare()
num, err := p.Exec("testing", "slene")
num, err  = p.Exec("testing", "astaxie")
...
...
p.Close() // 别忘记关闭 statement
```

## FAQ

1. 我的 app 需要支持多类型数据库，如何在使用 Raw SQL 的时候判断当前使用的数据库类型。

使用 Ormer 的 [Driver方法](Orm.md#driver) 可以进行判断
