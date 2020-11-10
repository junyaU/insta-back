package controllers

import (
	"instagram/models"
	"log"
	"strconv"

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
		log.Println("戻ります！")
		return
	}

	id, _ := this.GetInt64("postid")

	o := orm.NewOrm()
	post := models.Post{Id: id}

	m2m := o.QueryM2M(&post, "Favorite")

	userId := sessionUserId.(int64)
	user := models.User{Id: userId}

	o.Read(&user)

	num, err := m2m.Add(user)

	if err != nil {
		log.Println(err)
		log.Println("失敗だよ！")
	} else {
		log.Println(num)
		log.Println("だよ")
	}
}

func (this *FavoriteController) GetFavoriteUser() {
	inputId := this.Ctx.Input.Param(":id")
	postId, _ := strconv.ParseInt(inputId, 10, 64)
	post := models.Post{Id: postId}

	o := orm.NewOrm()
	o.Read(&post)
	o.LoadRelated(&post, "Favorite")

	this.Data["json"] = post

	this.ServeJSON()

}
