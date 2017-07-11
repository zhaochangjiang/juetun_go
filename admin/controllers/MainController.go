package controllers

import (
	acommon "juetun/admin/common"
	"log"
)

type MainController struct {
	acommon.AdminController
}

/**
* 301跳转页面
 */
func (this *MainController) Goto() {
	log.Println("----------------")
	log.Println(this.Ctx.Input.Params())
	log.Println("----------------")

	this.Ctx.Output.Body([]byte("您好"))
}
func (this *MainController) Get() {
	this.LoadCommon("default/index.html")
}
