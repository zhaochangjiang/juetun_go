package controllers

import (
	"fmt"
	acommon "juetun/admin/common"
)

type MainController struct {
	acommon.AdminController
}

func (this *MainController) Logout() {
	this.Ctx.Output.Body([]byte("退出登录"))
}
func (this *MainController) Get() {

	this.Data["UserId"] = this.GetSession("uid")
	this.Data["PageTitle"] = " 后台管理中心"
	this.Data["Avater"] = "/assets/img/user.jpg"
	this.Data["Username"] = "长江"
	this.Data["Items"] = this.GetPermitItem()
	this.LoadCommon("default/index.html")
}
