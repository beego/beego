## Object

对 object 操作简单的三个方法 Read / Insert / Update / Delete
```go
o := orm.NewOrm()
user := NewUser()
user.UserName = "slene"
user.Password = "password"
user.Email = "vslene@gmail.com"

fmt.Println(o.Insert(user))

user.UserName = "Your"
fmt.Println(o.Update(user))

fmt.Println(o.Delete(user))
```
### Read
```go
o := orm.NewOrm()
user := User{Id: 1}

o.Read(&user)

if err == sql.ErrNoRows {
	fmt.Println("查询不到")
} else if err == orm.ErrMissPK {
	fmt.Println("找不到主键")
} else {
	fmt.Println(user.Id, user.UserName)
}
```
### Insert
```go
o := orm.NewOrm()
var user User
user.UserName = "slene"
user.Password = "password"
user.Email = "vslene@gmail.com"
user.IsActive = true

fmt.Println(o.Insert(&user))
fmt.Println(user.Id)
```
创建后会自动对 auto 的 field 赋值

### Update
```go
o := orm.NewOrm()
user := User{Id: 1}
if o.Read(&user) == nil {
	user.UserName = "MyName"
	o.Update(&user)
}
```
### Delete
```go
o := orm.NewOrm()
o.Delete(&User{Id: 1})
```
Delete 操作会对反向关系进行操作，此例中 Post 拥有一个到 User 的外键。删除 User 的时候。如果 on_delete 设置为默认的级联操作，将删除对应的 Post

删除以后会清除 auto field 的值
