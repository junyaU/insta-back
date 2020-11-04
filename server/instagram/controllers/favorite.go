package controllers

import (
	"instagram/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type FavoriteController struct {
	beego.Controller
}

func (this *FavoriteController) Favorite() {
	session := this.StartSession()
	sessionUserId := session.Get("UserId")
	name := session.Get("Name")
	email := session.Get("Email")

	if sessionUserId == nil || name == nil || email == nil {
		this.Redirect("/", 302)
		return
	}

	id, _ := this.GetInt64("postid")

	o := orm.NewOrm()
	post := models.Post{Id: id}

	m2m := o.QueryM2M(&post, "Favorite")

	userId := sessionUserId.(int64)
	user := models.User{Id: userId}

	o.Read(&user)

	m2m.Add(user)
}
