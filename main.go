package main

import (
	"juetun/general"
	"juetun/models"
	_ "juetun/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logs.SetLogger(logs.AdapterConsole, `{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	beego.Info("the server is starting...")
	//beego.SetLogFuncCall(true)
	//beego.SetStaticPath("/assets/web", "public")
	//beego.SetLogger("file", `{"filename":"logs/web.log"}`)

	//初始化语言包
	general.InitLanguage()

	//初始化数据库
	general.InitDatabase()

	//初始化Model
	models.InitModel()

	beego.SetStaticPath("/assets", "static/assets")

	beego.Run()
}
