package controllers

import (
	acommon "juetun/admin/common"

	"juetun/common/utils"
)

type MainController struct {
	acommon.AdminController
}

/**
* @author karl.zhao <zhaocj2009@hotmail.com>
* @date 2017/08/14
* 实现本结构体的基本加载，本文件中所有的界面不需要验证登录
*
 */
func (this *MainController) Prepare() {
	_, a := this.GetControllerAndAction()

	//设置不需要登录的Action
	var notNeedLogin = []interface{}{"Get"}
	if utils.InArrayOrSlice(a, notNeedLogin) {

		//设置本页面不需要验证权限
		this.ConContext.NotNeedValidatePermit = true

	}
	//不需要加载本页的JS
	this.ConContext.NeedRenderJs = false
	this.AdminController.Prepare()
}

/**
* @author karl.zhao <zhaocj2009@hotmail.com>
* @date 2017/08/14
* 首页访问默认页面
 */
func (this *MainController) Get() {
	this.LoadCommon("default/index.html")
}
