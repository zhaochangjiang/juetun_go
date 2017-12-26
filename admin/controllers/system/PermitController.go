package system

import (
	acommon "juetun/admin/common"
	"juetun/common/general"
	modelsAdmin "juetun/common/models/admin"
	"juetun/common/utils"
	"net/url"
	//	"strconv"
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
	this.Data["NowChidList"] = this.formatDomain(this.leftJoinUponPermit(nowChidList))

	this.ConContext.IncludePageProperty.HaveTable = true    //设置有表格属性
	this.ConContext.IncludePageProperty.HaveCheckbox = true //设置有checkbox属性
	this.initCommonData(id)
	this.LoadCommon("permit/list.html")
}

//
func (this *PermitController) formatDomain(list *[]modelsAdmin.PermitAdmin) *[]modelsAdmin.PermitAdmin {

	var domainConfig = this.getDomainCache()
	var domainMap = make(map[string]string)
	for _, v := range *domainConfig {
		domainMap[v.Key] = v.Name
	}
	for k, v := range *list {
		if v.DomainMap != "" {
			(*list)[k].DomainMap = domainMap[v.DomainMap]
		}
	}
	return list
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

//获得当前权限的子权限
//@author karl.zhao<zhaochangjiang@huoyuren.com>
//@date 2017/10/16
//@params  pid string
//@return json字符串
func (this *PermitController) GetChild() {
	var id = this.GetString("pid")
	if id == "" {
		this.ConContext.OutputResult.Code = 100
		this.ConContext.OutputResult.Message = "您输入的id为空!"
	} else {
		permitAdmin := new(modelsAdmin.PermitAdmin)
		this.ConContext.OutputResult.Data = permitAdmin.GetAllChildByPids(&[]string{id},
			this.ConContext.IsSuperAdmin,
			&this.ConContext.GroupIds)
	}
	this.Data["json"] = this.ConContext.OutputResult
	this.ServeJSON()
}

//根据id获得权限信息
func (this *PermitController) getPermitById(id string) modelsAdmin.Permit {
	var fetchParams = make(map[string]interface{})
	fetchParams["id"] = id
	dataSingleton, _ := new(modelsAdmin.PermitAdmin).FetchPermit(fetchParams)
	var permit modelsAdmin.Permit

	if len(*dataSingleton) > 0 {
		permit = (*dataSingleton)[0]
	}
	return permit
}

//编辑界面
//@author karl.zhao<zhaocj2009@hotmail.com>
//@date 2017/09/12
func (this *PermitController) Edit() {
	var id = this.GetString("pid")
	var act = this.GetString("act")
	var gotoUrl = this.GetString("goto")
	this.Data["Error"] = ""
	var parent_id = ""
	var permitAdd = new(modelsAdmin.Permit)
	switch act {
	case "edit":
		permit := this.getPermitById(id)
		permitAdd = &permit
		if permit.Id != "" {
			parent_id = permit.UppermitId
			break
		}
		if gotoUrl == "" {
			gotoUrl = "/permit/list/"
		}
		this.Data["Error"] = "你要操作的权限不存在或已删除!"
		break
	default:
		break
	}
	this.Data["DataSingleton"] = permitAdd
	this.initCommonData(act, id, gotoUrl, parent_id)
	this.LoadCommon("permit/edit.html")
}

func (this *PermitController) getDomainCache() *[]modelsAdmin.Config {
	var config = (new(modelsAdmin.Config)).GetConfigByLikeKey("domain_")
	return config
}

//获得当前的域名配置列表
func (this *PermitController) getDomainConfig() *[]map[string]string {
	var res = make([]map[string]string, 0)
	var config = this.getDomainCache()
	for _, v := range *config {
		var domain = make(map[string]string)
		domain["DomainMap"] = v.Key
		domain["Name"] = v.Name
		res = append(res, domain)
	}
	return &res
}

//编辑界面
//@author karl.zhao<zhaocj2009@hotmail.com>
//@date 2017/09/12
func (this *PermitController) IframeEdit() {
	//var updateValue = make(map[string]string)
	var permit = new(modelsAdmin.Permit)

	permit.UppermitId = this.dealUpid()
	permit.DomainMap = this.GetString("domainMap")
	permit.Name = this.GetString("name")
	permit.Controller = this.GetString("controller")
	permit.Obyid = this.GetString("obyid")
	permit.Mod = this.GetString("module")
	permit.Action = this.GetString("action ")
	permit.Csscode = this.GetString("csscode")
	permit.Id = this.GetString("pid")

	var id = permit.UpdateDataById(permit)
	var res = new(general.Result)
	var returnValue = make(map[string]string)
	if id != false {
		res.Code = 0
		res.Message = "操作成功"
		returnValue["Goto"] = this.GetString("goto")
	} else {
		res.Code = 1001
		res.Message = "操作失败"
	}
	res.Data = returnValue
	this.Data["Result"] = res
	this.LoadCommon("layout/iframe.html")
}

func (this *PermitController) dealUpid() string {
	var pids = this.GetStrings("uppid[]")
	var pid = ""
	for _, v := range pids {
		if v != "-1" {
			pid = v
		} else {
			break
		}
	}
	return pid
}
func (this *PermitController) getIframeEditParams() *(map[string]interface{}) {
	var params = make(map[string]interface{})

	return &params
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

/**
*
* 放置公共的操作信息
 */
func (this *PermitController) initCommonData(arg ...interface{}) {

	_, action := this.GetControllerAndAction()
	switch action {
	case "List":

		this.Data["PageSmallTitle"], this.Data["TableTitle"] = "权限管理", "权限管理"
		this.Data["Pid"] = arg[0]
		this.Data["Currenturl"] = url.QueryEscape(this.Ctx.Request.RequestURI)
	case "Edit":
		var nowChidList = new([]modelsAdmin.PermitAdmin)
		this.Data["DomainConfig"] = this.getDomainConfig()
		this.Data["PList"], nowChidList = this.getPermitList(arg[3].(string))
		this.Data["NowChidList"] = this.leftJoinUponPermit(nowChidList)
		this.Data["PageTitle"] = "编辑"
		this.Data["DoAct"] = arg[0]
		this.Data["Pid"] = arg[1]
		this.Data["Goto"] = arg[2]
	}

}
