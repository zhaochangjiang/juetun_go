package controllers

import (
	acommon "juetun/admin/common"
)

type MainController struct {
	acommon.AdminController
}

func (this *MainController) Logout() {
	this.Ctx.Output.Body([]byte("退出登录"))
}
func (this *MainController) Get() {
	this.LoadCommon("default/index.html")
}
