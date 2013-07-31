## Object

对 object 操作的三个方法 Insert / Update / Delete

	o := orm.NewOrm()
	user := NewUser()
	user.UserName = "slene"
	user.Password = "password"
	user.Email = "vslene@gmail.com"
	obj := o.Object(user)
	fmt.Println(obj.Insert())
	user.UserName = "Your"
	fmt.Println(obj.Update())
	fmt.Println(obj.Delete())

### Read

	var user User
	err := o.QueryTable("user").Filter("id", 1).One(&user)
	if err != orm.ErrMultiRows {
		fmt.Println(user.UserName)
	}

### Create

	profile := NewProfile()
	profile.Age = 30
	profile.Money = 9.8

	user := NewUser()
	user.Profile = profile
	user.UserName = "slene"
	user.Password = "password"
	user.Email = "vslene@gmail.com"
	user.IsActive = true

	fmt.Println(o.Object(profile).Insert())
	fmt.Println(o.Object(user).Insert())
	fmt.Println(user.Id)
	
创建后会自动对 auto 的 field 赋值

### Update

	var user User
	err := o.QueryTable("user").Filter("id", 1).One(&user)
	if err != orm.ErrMultiRows {
		fmt.Println(user.UserName)
	}
	user.UserName = "MyName"
	o.Object(&user).Update()

### Delete
	
	o.Object(user).Delete()

Delete 操作会对反向关系进行操作，此例中 Post 拥有一个到 User 的外键。删除 User 的时候。如果 on_delete 设置为默认的级联操作，将删除对应的 Post

删除以后会清除 auto field 的值
