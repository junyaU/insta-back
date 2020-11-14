package controllers

import (
	"instagram/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Signup() {
	inputName := this.GetString("Name")
	inputEmail := this.GetString("Email")
	inputPassword := this.GetString("Password")

	hash, _ := bcrypt.GenerateFromPassword([]byte(inputPassword), 10)

	hashPassword := string(hash)

	o := orm.NewOrm()

	user := models.User{
		Name:     inputName,
		Email:    inputEmail,
		Password: hashPassword,
	}

	if _, id, err := o.ReadOrCreate(&user, "Email"); err == nil {
		session := this.StartSession()
		userId := session.Get("UserId")

		if userId != nil {
			this.Redirect("/", 302)
		}

		this.SetSession("UserId", id)
		this.SetSession("Name", inputName)

		sessionId := session.SessionID()
		user.SessionId = sessionId
		o.Update(&user, "SessionId")

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
	//ハッシュ値と平文を比較
	passwordError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword))

	if err == orm.ErrNoRows || passwordError != nil {
		this.Redirect("/login", 302)
		return
	}

	session := this.StartSession()

	userId := session.Get("UserId")
	Name := session.Get("Name")

	if userId != nil || Name != nil {
		this.Redirect("/posthome", 302)
		return
	}

	this.SetSession("UserId", user.Id)
	this.SetSession("Name", user.Name)

	//セッションIDをDBに保存しておく
	sessionId := session.SessionID()
	user.SessionId = sessionId
	o.Update(&user, "SessionId")

	this.Redirect("/postform", 302)
}

func (this *LoginController) Logout() {
	session := this.StartSession()

	userId := session.Get("UserId")
	name := session.Get("Name")

	if userId != nil || name != nil {
		session.Delete("UserId")
		session.Delete("Name")
	}

	this.Redirect("/", 302)
}
