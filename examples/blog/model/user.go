package model
import (
	"./beego"
)

type Users struct {
	username string
	password string
	beego.M
}

func NewUsers() (a *Users) {
	return &Users{}
}