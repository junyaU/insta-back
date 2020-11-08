package controllers

import (
	"github.com/astaxie/beego"
)

type SessionController struct {
	beego.Controller
}

type SessionData struct {
	UserId int64  `json:"userid"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func (this *SessionController) GetSessionData() {
	session := this.StartSession()

	sessionUserId := session.Get("UserId")
	sessionName := session.Get("Name")
	sessionEmail := session.Get("Email")

	if sessionUserId == nil || sessionName == nil || sessionEmail == nil {
		return
	}

	var data = SessionData{}

	data.UserId = sessionUserId.(int64)
	data.Name = sessionName.(string)
	data.Email = sessionEmail.(string)

	this.Data["json"] = data

	this.ServeJSON()
}
