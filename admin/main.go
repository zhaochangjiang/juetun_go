package main

import (
	_ "juetun/admin/routers"
	"juetun/general"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	//	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/redis"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logs.SetLogger(logs.AdapterConsole, `{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	beego.Info("the server is starting...")

	//初始化语言包
	general.InitLanguage()
	//初始化数据库
	general.InitDatabase()
	general.InitSession()

	beego.SetStaticPath("/assets", "static/assets")
	beego.Run()
}