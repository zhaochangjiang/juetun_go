package controllers

import (
	acommon "juetun/admin/common"
)

//权限设置相关功能
type PermitController struct {
	acommon.AdminController
}

func (this *PermitController) List() {

	this.SetListPageMessage()
	this.LoadCommon("permit/list.html")
}
