package routers

//前台路由设置文件
import (
	"juetun/controllers/frontend/web"

	"github.com/astaxie/beego"
)

func InitRouteFrontend() {
	beego.AutoRouter(&web.PassportController{})
}
