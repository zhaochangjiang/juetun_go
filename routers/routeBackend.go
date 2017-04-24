package routers

//后台路由设置文件
import (
	"github.com/astaxie/beego"

	"juetun/controllers/backend/dashboard"
)

func InitRouteBackend() {
	beego.AutoRouter(&dashboard.IndexController{})
}
