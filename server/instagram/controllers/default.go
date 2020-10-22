package controllers

import (
	"log"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	log.Println("hhelo")
}
