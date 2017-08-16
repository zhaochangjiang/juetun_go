package controllers

import (
	acommon "juetun/admin/common"
	"juetun/common/general"
	modelsAdmin "juetun/common/models/admin"
	"juetun/common/utils"
	"strings"

	"log"
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
	var notNeedLogin = []interface{}{"Goto"}
	if utils.InArrayOrSlice(a, notNeedLogin) {
		//设置本页面不需要登录
		this.ConContext.NotNeedLogin = true
	}

	this.AdminController.Prepare()
}

/**
* @author karl.zhao <zhaocj2009@hotmail.com>
* @date 2017/08/14
* 301跳转页面
 */
func (this *MainController) Goto() {
	getParams := this.Ctx.Input.Params()
	log.Println("MainController.go,line:37")
	//如果参数和呼标准,此处判断map类型key是否存在的方式，不适合数组切片的判断
	if _, ok := getParams["0"]; ok {
		permit := new(modelsAdmin.PermitAdmin)
		paramsL := strings.Split(getParams["0"], "_")

		//获得默认的moduleString值
		var moduleString string
		if len(paramsL) == 2 {
			moduleString = paramsL[1]
			//获得本module默认的访问路径
			permit.FetchDefaultPermitByModuleString(moduleString)
			this.Redirect(general.CreateUrl(permit.Controller, permit.Action, permit.Params, permit.Domain), 301)
			return
		}
	}
	this.Redirect("地址不正确!", 404)

}

/**
* @author karl.zhao <zhaocj2009@hotmail.com>
* @date 2017/08/14
* 首页访问默认页面
 */
func (this *MainController) Get() {
	this.LoadCommon("default/index.html")
}
