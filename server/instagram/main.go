package main

import (
	"fmt"
	"instagram/models"
	_ "instagram/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	setupDB()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

func setupDB() {
	db := beego.AppConfig.String("driver")

	orm.RegisterDriver(db, orm.DRMySQL)
	orm.RegisterDataBase("default", db, beego.AppConfig.String("sqlconn")+"?charset=utf8")
	orm.RegisterModel(
		new(models.User), new(models.Test), new(models.InstaUser), new(models.Post),
	)
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Println(err)
	}
}
