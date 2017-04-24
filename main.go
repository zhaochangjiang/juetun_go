package main

import (
	"juetun/common"
	_ "juetun/routers"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	//beego.SetLogFuncCall(true)
	//beego.SetStaticPath("/assets/web", "public")
	//beego.SetLogger("file", `{"filename":"logs/web.log"}`)

	//初始化语言包
	common.InitLanguage()

	//初始化数据库
	common.InitDatabase()

	beego.SetStaticPath("/assets", "static/assets")
	beego.Info("the server is starting...")
	beego.Run()
}
