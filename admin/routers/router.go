package routers

import (
	"juetun/admin/controllers/data"

	"juetun/admin/controllers/system"

	"github.com/astaxie/beego"
)

func init() {

	beego.AutoRouter(&data.DataController{})
	beego.AutoRouter(&data.UserController{})

	beego.AutoRouter(&system.SwitchController{})
	beego.AutoRouter(&system.LocController{})
	beego.AutoRouter(&system.Passport{})
	beego.AutoRouter(&system.GroupController{})
	beego.AutoRouter(&system.PermitController{})
	beego.AutoRouter(&system.ExportController{})

	beego.Router("/", &system.MainController{})
	//初始化错误信息
	beego.ErrorController(&system.ErrorController{})

}
