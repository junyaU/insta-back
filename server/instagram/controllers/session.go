package controllers

import (
	"instagram/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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

	if sessionUserId == nil {
		return
	}

	o := orm.NewOrm()
	sessionUser := models.User{Id: sessionUserId.(int64)}
	o.Read(&sessionUser)
	o.LoadRelated(&sessionUser, "FavoritePosts")

	this.Data["json"] = sessionUser

	this.ServeJSON()
}
