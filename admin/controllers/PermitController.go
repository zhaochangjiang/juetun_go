package controllers

import (
	acommon "juetun/admin/common"
	modelsAdmin "juetun/common/models/admin"
)

//权限设置相关功能
type PermitController struct {
	acommon.AdminController
}

//列表界面
func (this *PermitController) List() {
	var id = this.GetString("pid")
	var nowChidList = new([]modelsAdmin.PermitAdmin)
	this.Data["PList"], nowChidList = this.getPermitList(id)
	this.Data["NowChidList"] = this.leftJoinUponPermit(nowChidList)

	//设置有表格属性
	this.ConContext.IncludePageProperty.HaveTable = true

	//设置有checkbox属性
	this.ConContext.IncludePageProperty.HaveCheckbox = true

	this.Data["PageSmallTitle"], this.Data["TableTitle"] = "权限管理", "权限管理"
	this.Data["Pid"] = id
	this.LoadCommon("permit/list.html")
}

//编辑界面
func (this *PermitController) Edit() {
	this.LoadCommon("permit/edit.html")
}

//编辑界面
func (this *PermitController) IframeEdit() {
	request := this.Ctx.Request
	request.ParseForm()
	this.Debug(request.Form)

}

//处理上级权限名称
func (this *PermitController) leftJoinUponPermit(list *[]modelsAdmin.PermitAdmin) *[]modelsAdmin.PermitAdmin {

	var ids = this.getUponIdList(list)
	permit := new(modelsAdmin.Permit)

	var args = make(map[string]interface{})
	args["id__in"] = *ids
	permitList, _ := permit.FetchPermit(args)
	var permitListMap = make(map[string]modelsAdmin.Permit)
	for _, v := range *permitList {
		permitListMap[v.Id] = v
	}
	per := permit.GetModuleDefaultPermit(*permit)
	for k, v := range *list {
		//如果上级权限信息查询到了
		if _, ok := permitListMap[v.UppermitId]; ok {
			(*list)[k].UppermitId = permitListMap[v.UppermitId].Name

		}
		if v.Controller == (*per).Controller {
			(*list)[k].Controller = ""
		}
		if v.Action == (*per).Action {
			(*list)[k].Action = ""
		}
	}
	return list
}
func (this *PermitController) getUponIdList(list *[]modelsAdmin.PermitAdmin) *[]string {
	var ids = make([]string, 0)
	for _, v := range *list {
		ids = append(ids, v.UppermitId)
	}
	return &ids
}

func (this *PermitController) getPermitList(id string) (*[][]modelsAdmin.PermitAdmin, *[]modelsAdmin.PermitAdmin) {
	var args = make(map[string]interface{})

	args["IsSuperAdmin"] = this.ConContext.IsSuperAdmin
	args["Pid"] = id
	args["GroupIds"] = this.ConContext.GroupIds
	permit := new(modelsAdmin.Permit)
	return permit.GetList(args)
}
