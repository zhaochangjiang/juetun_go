package controllers

import (
	acommon "juetun/admin/common"
	modelsAdmin "juetun/common/models/admin"
	"juetun/common/utils"
)

//权限设置相关功能
//@author karl.zhao<zhaocj2009@hotmail.com>
//@date 2017/09/12
type PermitController struct {
	acommon.AdminController
}

/**
* @author karl.zhao <zhaocj2009@hotmail.com>
* @date 2017/08/14
* 实现本结构体的基本加载，本文件中所有的界面不需要验证登录
*
 */
func (this *PermitController) Prepare() {
	_, a := this.GetControllerAndAction()
	//设置不需要登录的Action
	var notNeedLogin = []interface{}{"GetChild"}
	if utils.InArrayOrSlice(a, notNeedLogin) {
		//设置本页面不需要验证权限
		this.ConContext.NotNeedValidatePermit = true
	}
	this.AdminController.Prepare()

}

//列表界面
//@author karl.zhao<zhaocj2009@hotmail.com>
//@date 2017/09/12
func (this *PermitController) List() {
	var id = this.GetString("pid")
	var nowChidList = new([]modelsAdmin.PermitAdmin)
	this.Data["PList"], nowChidList = this.getPermitList(id)
	this.Data["NowChidList"] = this.leftJoinUponPermit(nowChidList)

	this.ConContext.IncludePageProperty.HaveTable = true    //设置有表格属性
	this.ConContext.IncludePageProperty.HaveCheckbox = true //设置有checkbox属性

	this.Data["PageSmallTitle"], this.Data["TableTitle"] = "权限管理", "权限管理"
	this.Data["Pid"] = id
	this.Data["Currenturl"] = this.Ctx.Request.RequestURI
	this.LoadCommon("permit/list.html")
}

//删除权限
func (this *PermitController) Del() {
	var id = this.GetString("pid")
	permitAdmin := new(modelsAdmin.PermitAdmin)
	var pid = []string{id}
	res, err := permitAdmin.DeletePermit(pid)
	if err != nil {
		this.ConContext.OutputResult.Code = 100
		this.ConContext.OutputResult.Message = "操作失败"
		this.ConContext.OutputResult.Data = err.Error()
	} else {
		this.ConContext.OutputResult.Data = res
	}
	this.Data["json"] = this.ConContext.OutputResult
	this.ServeJSON()

}
func (this *PermitController) GetChild() {
	var id = this.GetString("pid")
	if id == "" {
		this.ConContext.OutputResult.Code = 100
		this.ConContext.OutputResult.Message = "您输入的id为空!"
	} else {

		permitAdmin := new(modelsAdmin.PermitAdmin)
		this.ConContext.OutputResult.Data = permitAdmin.GetAllChildByPids(&[]string{id},
			//this.ConContext.IsSuperAdmin,
			true,
			&this.ConContext.GroupIds)
	}
	this.Data["json"] = this.ConContext.OutputResult
	this.ServeJSON()
}

//编辑界面
//@author karl.zhao<zhaocj2009@hotmail.com>
//@date 2017/09/12
func (this *PermitController) Edit() {
	var id = this.GetString("pid")
	var act = this.GetString("act")
	var gotoUrl = this.GetString("goto")

	var nowChidList = new([]modelsAdmin.PermitAdmin)

	this.Data["PList"], nowChidList = this.getPermitList(id)
	this.Data["NowChidList"] = this.leftJoinUponPermit(nowChidList)
	this.Data["PageTitle"] = "编辑"
	this.Data["DoAct"] = act
	this.Data["Goto"] = gotoUrl
	this.Data["Pid"] = id
	this.LoadCommon("permit/edit.html")
}

//编辑界面
//@author karl.zhao<zhaocj2009@hotmail.com>
//@date 2017/09/12
func (this *PermitController) IframeEdit() {
	request := this.Ctx.Request
	request.ParseForm()
	this.Debug(request.Form)

}

//处理上级权限名称
//@author karl.zhao<zhaocj2009@hotmail.com>
//@date 2017/09/12
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

//@author karl.zhao<zhaocj2009@hotmail.com>
//@date 2017/09/12
func (this *PermitController) getUponIdList(list *[]modelsAdmin.PermitAdmin) *[]string {
	var ids = make([]string, 0)
	for _, v := range *list {
		ids = append(ids, v.UppermitId)
	}
	return &ids
}

//@author karl.zhao<zhaocj2009@hotmail.com>
//@date 2017/09/12
func (this *PermitController) getPermitList(id string) (*[][]modelsAdmin.PermitAdmin, *[]modelsAdmin.PermitAdmin) {
	var args = make(map[string]interface{})

	args["IsSuperAdmin"] = this.ConContext.IsSuperAdmin
	args["Pid"] = id
	args["GroupIds"] = this.ConContext.GroupIds
	permit := new(modelsAdmin.Permit)
	return permit.GetList(args)
}
