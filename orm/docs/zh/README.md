## beego orm 介绍

## 快速入门

```go
package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// 最简单的两个表的定义
type User struct {
	Id   int    `orm:"auto"`
	Name string `orm:"size(100)"`
	orm.Manager
}

func init() {
	// 将表定义注册到 orm 里
	orm.RegisterModel(new(User))

	// 链接参数设置
	orm.RegisterDataBase("default", "mysql", "root:root@/my_db?charset=utf8", 30)
}

func main() {
	o := orm.NewOrm()

	user := User{Name: "slene"}

	// 创建
	id, err := o.Insert(&user)
	fmt.Println(id, err)

	// 更新
	user.Name = "astaxie"
	num, err := o.Update(&user)
	fmt.Println(num, err)

	// 查询单个
	u := User{Id: user.Id}
	err = o.Read(&u)
	fmt.Println(u.Name, err)

	// 删除
	num, err = o.Delete(&u)
	fmt.Println(num, err)
}
```

## 详细文档

1. [模型定义](Models.md)
	- [支持的 Field 类型](Models.md#Field Type)
	- [Field 设置参数](Models.md#Field Options)
	- [关系型 Field 设置](Models.md#Relation Field Options)
2. Custom Fields
3. [Orm 使用方法](Orm.md)
	- [Ormer 接口](Orm.md#Ormer)
4. [对象操作](Object.md)
5. [复杂查询](Query.md)
	- [查询使用的表达式语法](Query.md#expr)
	- [查询支持的操作符号](Query.md#Operators)
	- [QuerySeter 接口](Query.md#QuerySeter)
6. Raw
7. Transaction
8. Faq
