package controllers

import (
	acommon "juetun/admin/common"
)

type UserController struct {
	acommon.AdminController
}

func (this *UserController) List() {
	this.Debug("UserController_List")
	this.LoadCommon("layout/list.html")
}
