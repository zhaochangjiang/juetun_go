package routers

import (
	web "juetun/web/controllers"

	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/", &web.Default{})
	beego.AutoRouter(&web.Passport{})

}
