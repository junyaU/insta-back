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
	Id   int64  `json:"Id"`
	Name string `json:"Name"`
}

func (this *SessionController) GetSessionData() {
	session := this.StartSession()
	sessionUserId := session.Get("UserId")
	if sessionUserId == nil {
		return
	}

	sessionUser := models.User{Id: sessionUserId.(int64)}
	o := orm.NewOrm()
	o.Read(&sessionUser)
	sessionId := session.SessionID()
	//DBに保存された値とsessionIdが一致しなければリターン
	if sessionId != sessionUser.SessionId {
		return
	}

	o.LoadRelated(&sessionUser, "FavoritePosts")

	sessData := SessionData{}
	sessData.Id = sessionUser.Id
	sessData.Name = sessionUser.Name

	this.Data["json"] = sessData
	this.ServeJSON()
}

func (this *SessionController) AuthCheck() {
	o := orm.NewOrm()
	session := this.StartSession()
	sessionUserId := session.Get("UserId")
	if sessionUserId == nil {
		return
	}

	sessionId := session.SessionID()
	sessionUser := models.User{Id: sessionUserId.(int64)}
	o.Read(&sessionUser)
	if sessionUser.SessionId != sessionId {
		return
	}
}
