package main

import (
	//初始化功能设置
	_ "juetun/admin/routers"
	_ "juetun/common/general"

	//引入缓存包
	_ "github.com/astaxie/beego/cache/redis"

	"github.com/astaxie/beego"
)

func main() {
	beego.Info("the server is starting...")
	//引入静态文件路径
	beego.SetStaticPath("/assets", "static/assets")
	beego.Run()
}
