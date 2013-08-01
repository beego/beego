## Orm

beego/orm 的使用方法
```go
package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// 这个用来设置 driverName 对应的数据库类型
	// mysql / sqlite3 / postgres 这三种是默认已经注册过的，所以可以无需设置
	orm.RegisterDriver("mysql", orm.DR_MySQL)

	// 参数1   自定义的数据库名称，用来在orm中切换数据库使用
	// 参数2   driverName
	// 参数3   对应的链接字符串
	// 参数4   设置最大的空闲连接数，使用 golang 自己的连接池
	orm.RegisterDataBase("default", "mysql", "root:root@/orm_test?charset=utf8", 30)
}

func main() {
	// 请确保在所有 RegisterModel 之前执行
	orm.BootStrap() // 强制在 main 函数里调用，检查 Model 关系，检测数据库参数，调用 orm 提供的 Command

	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库

	profile := NewProfile()
	profile.Age = 30
	profile.Money = 9.8

	user := NewUser()
	user.Profile = profile
	user.UserName = "slene"
	user.Password = "password"
	user.Email = "vslene@gmail.com"
	user.IsActive = true

	fmt.Println(o.Insert(profile))
	fmt.Println(o.Insert(user))

	var params []orm.Params
	if cnt, err := o.QueryTable("user").RelatedSel().Limit(3).OrderBy("-id").Values(&params); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(cnt)
		for _, p := range params {
			fmt.Println(p)
		}
	}

	var users []*User
	if cnt, err := o.QueryTable("user").RelatedSel().Limit(3).OrderBy("-id").All(&users); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(cnt)
		for _, u := range users {
			fmt.Println(u.Id, u.Profile)
		}
	}
}
```