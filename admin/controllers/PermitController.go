package controllers

import (
	acommon "juetun/admin/common"
	modelsAdmin "juetun/common/models/admin"
)

//权限设置相关功能
type PermitController struct {
	acommon.AdminController
}

func (this *PermitController) List() {
	var id = this.GetString("pid")

	this.Data["PList"] = this.getPermitList(id)
	this.SetListPageMessage()
	this.LoadCommon("permit/list.html")
}

func (this *PermitController) getPermitList(id string) *[][]modelsAdmin.PermitAdmin {
	var args = make(map[string]interface{})

	args["IsSuperAdmin"] = this.ConContext.IsSuperAdmin
	args["Pid"] = id
	args["GroupIds"] = this.ConContext.GroupIds
	permit := new(modelsAdmin.Permit)
	return permit.GetList(args)
}
