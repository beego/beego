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
	orm.BootStrap() // 确保在所有 RegisterModel 之后执行

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

1. [Model Definition](Models.md)
2. Custom Fields
3. [Orm](Orm.md)
4. [Object](Object.md)
5. [Query](Query.md)
6. Condition
7. Raw
8. Transaction
9. Faq
