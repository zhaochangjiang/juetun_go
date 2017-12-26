package system

import (
	acommon "juetun/admin/common"
)

type GroupController struct {
	acommon.AdminController
}

/**
* 基础数据管理
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/16
* @return void
 */
func (this *GroupController) List() {
	//	this.Data["UserId"] = this.GetSession("uid")
	//	this.Data["PageTitle"] = " 后台管理中心"
	//	this.Data["Avater"] = "/assets/img/user.jpg"
	//	this.Data["Username"] = "长江"
	this.LoadCommon("group/list.html")
}
