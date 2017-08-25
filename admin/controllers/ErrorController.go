package controllers

import (
	acommon "juetun/admin/common"
)

type ErrorController struct {
	acommon.AdminController
}

//实现本结构体的基本加载，本文件中所有的界面不需要验证登录
func (this *ErrorController) Prepare() {
	//	_,a:=this.GetControllerAndAction()
	//	utils:=new(general.Utils)

	//	//设置不需要登录的Action
	//	var notNeedLogin=[]interface{"Goto"}
	//	if(utils.InArrayOrSlice(a,notNeedLogin)){
	//设置本页面不需要登录
	this.ConContext.NotNeedLogin = true

	this.Data["Breadcrumbs"] = ""
	this.AdminController.Prepare()
}

/**
* 404界面
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/16
*
 */
func (this *ErrorController) Error404() {
	this.Data["Content"] = "page not found"
	this.LoadCommon("error/404.html")
	//this.TplName = "404.tpl"
}

/**
* 501界面
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/16
*
 */
func (this *ErrorController) Error501() {
	this.Data["Content"] = "server error"
	this.LoadCommon("error/501.html")
}

/**
* 数据库异常提醒
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/16
*
 */
func (this *ErrorController) ErrorDb() {
	this.Data["Content"] = "database is now down"
	this.LoadCommon("error/db.html")
}
