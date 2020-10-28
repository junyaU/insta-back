package controllers

import (
	"instagram/models"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Signup() {
	inputName := this.GetString("Name")
	inputEmail := this.GetString("Email")
	inputPassword := this.GetString("Password")

	o := orm.NewOrm()

	user := models.User{
		Name:     inputName,
		Email:    inputEmail,
		Password: inputPassword,
	}

	if created, id, err := o.ReadOrCreate(&user, "Email"); err == nil {
		if created {
			log.Println("insert:", id)
		} else {
			log.Println("read:", id)
		}
	}
}

func (this *LoginController) Login() {

}

func (this *LoginController) Logout() {

}
