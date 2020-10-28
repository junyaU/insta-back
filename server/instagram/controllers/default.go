package controllers

import (
	"instagram/models"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	my := models.Test{Id: 1, Content: "hello"}
	this.Data["json"] = &my
	this.ServeJSON()
}
