package common

import (
	"juetun/common/general"

	//"github.com/astaxie/beego"
)

type FrontInterface interface {
	Display(template ...string)
}

//前台公共入口
type FrontendController struct {
	general.BaseController
	userid         int64
	username       string
	nickname       string
	ControllerName string
	actionName     string
	frontend       string
}
