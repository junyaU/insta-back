package controllers

import (
	"instagram/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type PostController struct {
	beego.Controller
}

func (this *PostController) GetAllPosts() {
	o := orm.NewOrm()
	var allPosts []models.Post
	o.QueryTable(new(models.Post)).RelatedSel().All(&allPosts)

	this.Data["json"] = allPosts

	this.ServeJSON()
}

func (this *PostController) Post() {
	session := this.StartSession()
	userId := session.Get("UserId")
	Name := session.Get("Name")
	Email := session.Get("Email")

	if userId == nil || Name == nil || Email == nil {
		this.Redirect("/", 302)
		return
	}

	id := userId.(int64)

	inputComment := this.GetString("Comment")
	inputImage := this.GetString("Image")

	post := models.Post{
		Comment: inputComment,
		User:    &models.User{Id: id},
		Image:   inputImage,
	}

	o := orm.NewOrm()
	o.Insert(&post)

	this.Redirect("/posthome", 302)
}
