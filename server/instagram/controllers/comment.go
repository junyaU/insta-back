package controllers

import (
	"instagram/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type CommentController struct {
	beego.Controller
}

func (this *CommentController) Comment() {
	userId, _ := this.GetInt64("userId")
	postId, _ := this.GetInt64("postId")
	text := this.GetString("text")

	o := orm.NewOrm()

	comment := models.Comment{
		User:    &models.User{Id: userId},
		Post:    &models.Post{Id: postId},
		Comment: text,
	}

	o.Insert(&comment)
}
