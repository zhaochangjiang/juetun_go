package controllers

import (
	acommon "juetun/admin/common"
	"juetun/common/general"
	modelsAdmin "juetun/common/models/admin"
	"strings"

	"log"
)

type MainController struct {
	acommon.AdminController
}

/**
* 301跳转页面
 */
func (this *MainController) Goto() {
	getParams := this.Ctx.Input.Params()

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
func (this *MainController) Get() {
	this.LoadCommon("default/index.html")
}
