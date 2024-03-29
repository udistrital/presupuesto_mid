package main

import (
	"github.com/astaxie/beego/plugins/cors"
	"github.com/udistrital/presupuesto_mid/mongoProcess"
	_ "github.com/udistrital/presupuesto_mid/routers"

	"github.com/astaxie/beego"
)

func init() {
}

func main() {

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	mongoProcess.PresupuestoMongoJobInit()
	beego.Run()

}
