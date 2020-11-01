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

	if _, id, err := o.ReadOrCreate(&user, "Email"); err == nil {
		session := this.StartSession()
		userId := session.Get("UserId")

		if userId != nil {
			this.Redirect("/", 302)
		}

		this.SetSession("UserId", id)
		this.SetSession("Name", inputName)
		this.SetSession("Email", inputEmail)

		this.Redirect("/posthome", 302)
	}

	this.Redirect("/", 302)
}

func (this *LoginController) Login() {
	inputEmail := this.GetString("Email")
	inputPassword := this.GetString("Password")

	o := orm.NewOrm()
	user := models.User{Email: inputEmail}
	err := o.Read(&user, "Email")

	if err == orm.ErrNoRows {
		log.Println("そんなユーザーいないよ")
		this.Redirect("/login", 302)
		return
	} else if inputPassword != user.Password {

		log.Println("パスワードが違います")
		this.Redirect("/login", 302)
		return
	}

	session := this.StartSession()

	userId := session.Get("UserID")
	Name := session.Get("Name")
	Email := session.Get("Email")

	if userId != nil || Name != nil || Email != nil {
		this.Redirect("/posthome", 302)
		return
	}

	this.SetSession("UserId", user.Id)
	this.SetSession("Name", user.Name)
	this.SetSession("Email", user.Email)

	log.Println(inputEmail)
	log.Println(inputPassword)
	log.Println(user.Name)

	this.Redirect("/postform", 302)
}

func (this *LoginController) Logout() {

}
