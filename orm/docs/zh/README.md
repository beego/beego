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

type Post struct {
	Id      int    `orm:"auto"`
	User    *User  `orm:"rel(fk)"`
	Title   string `orm:"size(100)"`
	Content string
	orm.Manager
}

func init() {
	// 将表定义注册到 orm 里
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Post))

	// 链接参数设置
	orm.RegisterDataBase("default", "mysql", "root:root@/my_db?charset=utf8", 30)
}

func main() {
	orm.BootStrap() // 确保在所有 RegisterModel 之后执行

	o := orm.NewOrm()

	var user User
	obj := o.Object(&user)

	// 创建
	user.Name = "slene"
	id, err := obj.Insert()
	fmt.Println(id, err)

	// 更新
	user.Name = "astaxie"
	num, err := obj.Update()
	fmt.Println(num, err)

	// 查询单个
	var u User
	err = o.QueryTable("user").Filter("id", &user).One(&u)
	fmt.Println(u.Id, u.Name, err)

	// 创建 post
	var post Post
	post.Title = "beego orm"
	post.Content = "powerful amazing"
	post.User = &u
	id, err = o.Object(&post).Insert()
	fmt.Println(id, err)
	
	// 当然，以 map[string]interface{} 形式的数据返回也是允许的
	var maps []orm.Params
	num, err = o.QueryTable("user").Filter("id", &u).Values(&maps)
	fmt.Println(num, err, maps[0])

	// 删除
	num, err = obj.Delete() // 默认，级联删除 user 以及关系存在的 post
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
