package routers

import (
	"juetun/admin/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.AutoRouter(&controllers.Passport{})
	beego.AutoRouter(&controllers.Data{})
	beego.Router("/", &controllers.MainController{})

}
