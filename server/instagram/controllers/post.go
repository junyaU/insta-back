package controllers

import (
	"instagram/models"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type PostController struct {
	beego.Controller
}

func (this *PostController) GetAllPosts() {
	o := orm.NewOrm()
	var allPosts []models.Post
	o.QueryTable(new(models.Post)).All(&allPosts)

	this.Data["json"] = allPosts

	this.ServeJSON()
}

func (this *PostController) Post() {
	inputComment := this.GetString("Comment")
	inputImage := this.GetString("Image")

	o := orm.NewOrm()
	post := models.Post{
		Comment: inputComment,
		User:    &models.User{Id: 1},
		Image:   inputImage,
	}

	id, err := o.Insert(&post)

	if err == nil {
		log.Println(id)
	} else {
		log.Println(err)
	}

	this.Redirect("/posthome", 302)

}
