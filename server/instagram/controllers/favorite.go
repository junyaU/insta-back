package controllers

import (
	"instagram/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type FavoriteController struct {
	beego.Controller
}

func (this *FavoriteController) Favorite() {
	userId, _ := this.GetInt64("userid")
	postId, _ := this.GetInt64("postid")

	post := models.Post{Id: postId}
	user := models.User{Id: userId}

	o := orm.NewOrm()
	o.Read(&user)
	m2m := o.QueryM2M(&post, "Favorite")
	m2m.Add(user)
}

func (this *FavoriteController) UnFavorite() {
	userId, _ := this.GetInt64("userid")
	postId, _ := this.GetInt64("postid")

	post := models.Post{Id: postId}
	user := models.User{Id: userId}

	o := orm.NewOrm()
	o.Read(&user)
	m2m := o.QueryM2M(&post, "Favorite")
	m2m.Remove(user)
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
