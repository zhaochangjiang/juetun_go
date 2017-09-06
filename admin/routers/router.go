package routers

import (
	"juetun/admin/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.AutoRouter(&controllers.LocController{})
	beego.AutoRouter(&controllers.Passport{})
	beego.AutoRouter(&controllers.DataController{})
	beego.AutoRouter(&controllers.GroupController{})
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.PermitController{})

	//初始化错误信息
	beego.ErrorController(&controllers.ErrorController{})

}
