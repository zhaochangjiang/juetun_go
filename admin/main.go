package main

import (
	_ "juetun/admin/routers"
	"juetun/common/general"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	_ "github.com/astaxie/beego/session/redis"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	logs.SetLogger(logs.AdapterFile, `{"filename":"success.log","level":7,"maxlines":10000,"maxsize":0,"daily":true,"maxdays":10}`)
	logs.Async() //为了提升性能, 可以设置异步输出:

	logs.Async(1e3) //异步输出允许设置缓冲 chan 的大小

	//	logs.Debug("my book is bought in the year of ", 2016)
	//	logs.Info("this %s cat is %v years old", "yellow", 3)
	//	logs.Warn("json is a type of kv like", map[string]int{"key": 2016})
	//	logs.Error(1024, "is a very", "good game")

	beego.Info("the server is starting...")

	//初始化语言包
	general.InitLanguage()
	//初始化数据库
	general.InitDatabase()

	general.InitSession()
	//初始化模板函数
	general.InitAddFuncMap()

	beego.SetStaticPath("/assets", "static/assets")
	beego.Run()
}
