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

var successCheck struct {
	Status string `json:"status"`
}

func (this *LoginController) Signup() {
	inputName := this.GetString("Name")
	inputEmail := this.GetString("Email")
	inputPassword := this.GetString("Password")
	status := successCheck

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
			status.Status = "failed"
			this.Data["json"] = status
			this.ServeJSON()
			return
		}

		this.SetSession("UserId", id)
		this.SetSession("Name", inputName)

		sessionId := session.SessionID()
		user.SessionId = sessionId

		this.Ctx.SetCookie("instaSessionID", sessionId)
		o.Update(&user, "SessionId")

		status.Status = "success"
		this.Data["json"] = status
		this.ServeJSON()
	}
}

func (this *LoginController) Login() {
	inputEmail := this.GetString("Email")
	inputPassword := this.GetString("Password")

	o := orm.NewOrm()
	user := models.User{Email: inputEmail}
	err := o.Read(&user, "Email")
	status := successCheck

	//ハッシュ値と平文を比較
	passwordError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword))

	if err == orm.ErrNoRows || passwordError != nil {
		status.Status = "failed"
		this.Data["json"] = status
		this.ServeJSON()
		return
	}

	session := this.StartSession()

	userId := session.Get("UserId")
	Name := session.Get("Name")

	if userId != nil || Name != nil {
		status.Status = "failed"
		this.Data["json"] = status
		this.ServeJSON()
		return
	}

	this.SetSession("UserId", user.Id)
	this.SetSession("Name", user.Name)

	//セッションIDをDBに保存しておく
	sessionId := session.SessionID()
	user.SessionId = sessionId
	o.Update(&user, "SessionId")
	this.Ctx.SetCookie("instaSessionID", sessionId)

	status.Status = "success"
	this.Data["json"] = status
	this.ServeJSON()
}

func (this *LoginController) Logout() {
	session := this.StartSession()

	userId := session.Get("UserId")
	name := session.Get("Name")

	if userId != nil || name != nil {
		session.Delete("UserId")
		session.Delete("Name")
	}
}
