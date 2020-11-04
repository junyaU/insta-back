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

	o.QueryTable(new(models.Post)).RelatedSel().All(&allPosts)

	var afterPost []models.Post

	for _, post := range allPosts {

		m2m := o.QueryM2M(&post, "Favorite")

		nums, _ := m2m.Count()

		post.Favonum = nums

		afterPost = append(afterPost, post)
	}

	this.Data["json"] = afterPost

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

	aaa := models.Post{Id: 1}
	o.Read(&aaa)

	log.Println(aaa)

	// log.Println(aaa)

	this.Redirect("/posthome", 302)
}
