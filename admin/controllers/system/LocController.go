package system

import (
	acommon "juetun/admin/common"
	"juetun/common/general"
	modelsAdmin "juetun/common/models/admin"
	"juetun/common/utils"
	"strings"
	"time"

	wetalkutils "github.com/beego/wetalk/modules/utils"
)

//跳转控制界面处理
type LocController struct {
	acommon.AdminController
}

func (this *LocController) SetPaginator(per int, nums int64) *wetalkutils.Paginator {
	p := wetalkutils.NewPaginator(this.Ctx.Request, per, nums)
	this.Data["paginator"] = p
	return p
}
func (this *LocController) Test() {
	this.SetPaginator(20, 100)

	this.LoadCommon("layout/pageinator.html")
}

/**
* @author karl.zhao <zhaocj2009@hotmail.com>
* @date 2017/08/14
* 301跳转页面
 */
func (this *LocController) Goto() {

	getParams := this.Ctx.Input.Params()
	//如果参数和呼标准,此处判断map类型key是否存在的方式，不适合数组切片的判断
	if _, ok := getParams["0"]; ok {
		permit := new(modelsAdmin.PermitAdmin)
		paramsL := strings.Split(getParams["0"], "_")
		time.Sleep(500 * time.Microsecond)
		//获得默认的moduleString值
		if len(paramsL) == 2 {
			if paramsL[1] == "dashboard" {
				this.Redirect("/", 301)
				return
			}
			//获得本module默认的访问路径
			permit.FetchDefaultPermitByModuleString(paramsL[1], this.ConContext)
			if permit.Controller == "" {
				this.Data["Message"] = "没有找到您选择跳转的链接!"
				this.Data["HttpCode"] = "500"
				this.LoadCommon("error/goto.html")
				return
			}
			gotoUrl := general.CreateUrl(permit.Controller, permit.Action, permit.Params, permit.Domain)

			this.Redirect(gotoUrl, 301)
			return
		}
		this.Abort("404")
	}
	this.Data["HttpCode"] = "500"
	this.Data["Message"] = "您没有选择跳转的链接!"
	this.ConContext.NeedRenderJs = false
	this.LoadCommon("error/goto.html")
}

/**
* @author karl.zhao <zhaocj2009@hotmail.com>
* @date 2017/08/14
* 实现本结构体的基本加载，本文件中所有的界面不需要验证登录
*
 */
func (this *LocController) Prepare() {
	_, a := this.GetControllerAndAction()

	//设置不需要登录的Action
	var notNeedLogin = []interface{}{"Goto"}
	if utils.InArrayOrSlice(a, notNeedLogin) {

		//设置本页面不需要验证权限
		this.ConContext.NotNeedValidatePermit = true
	}

	this.AdminController.Prepare()
}
