package routers

import (
	"juetun/admin/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//初始化错误信息
	beego.ErrorController(&controllers.ErrorController{})

	beego.AutoRouter(&controllers.Passport{})
	beego.AutoRouter(&controllers.DataController{})
	beego.AutoRouter(&controllers.MainController{})

	beego.Router("/", &controllers.MainController{})

}
