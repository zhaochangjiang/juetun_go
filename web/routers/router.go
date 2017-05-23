package routers

import (
	web "juetun/web/controllers"

	"github.com/astaxie/beego"
)

func init() {

	//	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
	//		AllowAllOrigins:  true,
	//		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	//		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	//		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	//		AllowCredentials: true,
	//	}))
	beego.Router("/", &web.Default{})
	beego.AutoRouter(&web.Passport{})

}
