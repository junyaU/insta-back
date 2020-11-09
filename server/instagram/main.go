package main

import (
	"fmt"
	"instagram/models"
	_ "instagram/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/session"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	setupDB()

	sessionconf := &session.ManagerConfig{
		CookieName:      "instasessionId",
		Gclifetime:      3600,
		EnableSetCookie: true,
		Maxlifetime:     3600,
		CookieLifeTime:  3600,
		ProviderConfig:  "",
	}
	beego.GlobalSessions, _ = session.NewManager("memory", sessionconf)
	go beego.GlobalSessions.GC()

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

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
		new(models.User), new(models.Post), new(models.Imageprofile),
	)
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Println(err)
	}
}
