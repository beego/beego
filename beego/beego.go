package beego

import "./core"

type C struct {
	core.Content
}

type M struct{
	core.Model
}

type D struct{
	core.Config
}

type U struct{
	core.URL
}

type A struct{
	core.Controller
}

type V struct{
	core.View
}