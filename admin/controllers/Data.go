package controllers

import (
	acommon "juetun/admin/common"
)

type Data struct {
	acommon.AdminController
}

func (this *Data) List() {

	this.Data["UserId"] = this.GetSession("uid")
	this.Data["PageTitle"] = " 后台管理中心"
	this.Data["Avater"] = "/assets/img/user.jpg"
	this.Data["Username"] = "长江"
	this.LoadCommon("data/list.html")
}
