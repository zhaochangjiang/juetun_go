package fcommon

import (
	"fmt"
	"juetun/general"
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
	controllerName string
	actionName     string
	frontend       string
}

func (this *FrontendController) Display(template ...string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("出了错：", err)
		}
	}()

	if nil != template {
		if len(template) > 1 {
			panic("参数最多1位或者没有参数!")
		}
		this.TplName = "frontend/" + template[0]
	}

}
