package routers

import (
	"juetun/admin/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.AutoRouter(&controllers.Passport{})
	beego.AutoRouter(&controllers.DataController{})
	beego.AutoRouter(&controllers.MainController{})

	beego.Router("/", &controllers.MainController{})

}
