package system

import (
	acommon "juetun/admin/common"
	"strings"

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
	var notNeedLogin = []interface{}{"get"} //请填写小写
	if utils.InArrayOrSlice(strings.ToLower(a), notNeedLogin) == true {
		//设置本页面不需要验证权限
		this.ConContext.NotNeedValidatePermit = false
	}
	this.AdminController.Prepare()

	//不需要加载本页的JS,必须放在Prepare方法之后
	this.ConContext.NeedRenderJs = false
}

/**
* @author karl.zhao <zhaocj2009@hotmail.com>
* @date 2017/08/14
* 首页访问默认页面
 */
func (this *MainController) Get() {
	this.LoadCommon("default/index.html")
}
