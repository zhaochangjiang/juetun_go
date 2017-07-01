package common

import (
	"errors"

	"strings"
	//	"fmt"
	"juetun/common/general"
	modelsAdmin "juetun/common/models/admin"
)

type AdminController struct {
	PermitService *modelsAdmin.Permit
	general.BaseController
}

//返回当前后台的权限列表
func (this *AdminController) InitPermitItem() {

	//初始化权限操作Service
	this.PermitService = new(modelsAdmin.Permit)
	this.initAllShowPermit()

	//如果不是超级管理员
	if !this.authSuperAdmin() {
		//获得当前不是超级管理员的权限列表。
		this.Data["Permit"] = this.getListNotSuperAdmin()
	}

}
func (this *AdminController) DefaultControllerAndAction() (string, string) {
	return "MainController", "GET"
}

//获得当前的权限
func (this *AdminController) getNowPermitData() (*modelsAdmin.Permit, error) {
	permitModel := new(modelsAdmin.Permit)

	fetchParams := make(map[string]interface{})
	fetchParams["Controller"], fetchParams["Action"] = this.GetControllerAndAction()

	fetchParams["Controller"] = strings.ToLower(strings.TrimRight(fetchParams["Controller"].(string), "Controller"))
	fetchParams["Action"] = strings.ToLower(fetchParams["Action"].(string))

	defaultController, actionString := this.DefaultControllerAndAction()
	if defaultController == fetchParams["Controller"] && actionString == fetchParams["Action"] {
		return permitModel, errors.New("")
	}
	var permitModelList []*modelsAdmin.Permit
	permitModelList, err := this.PermitService.FetchPermit(fetchParams)
	if len(permitModelList) > 0 {
		permitModel = permitModelList[0]
	}
	return permitModel, err
}

//获得当前地址对应的数据库存储的权限及所有上级权限
func (this *AdminController) getNowAndAllUponPermit() (*[]interface{}, []interface{}, error) {

	permitModel := new(modelsAdmin.Permit)

	result := make([]interface{}, 0)
	utils := new(general.Utils)
	uponPermitId := make([]interface{}, 0)
	permitData, _ := this.getNowPermitData()

	//默认的上级机构必须查询
	uponPermitId = *utils.Slice_unshift(uponPermitId, 0)
	var err1 error
	var permitModelList []*modelsAdmin.Permit
	i := 0
	for {
		i++
		if "" == permitData.UppermitId || i > 5 {
			break
		}
		fetchParams := make(map[string]interface{})
		fetchParams["id"] = permitData.UppermitId
		uponPermitId = *utils.Slice_unshift(uponPermitId, permitData.UppermitId)
		permitModelList, err1 = permitModel.FetchPermit(fetchParams)

		if len(permitModelList) > 0 {
			permitData = (permitModelList[0])
			//往队列的队首添加数据
			result = *utils.Slice_unshift(result, permitData)
		} else {
			break
		}

	}
	return &result, uponPermitId, err1
}

//获得header默认的Type
func (this *AdminController) getHeaderDefaultActive(permitUpon []interface{}) (string, string) {
	headerActive := "dashboard"
	var activeId string
	length := len(permitUpon)
	if length > 0 {
		permit := permitUpon[0].(*modelsAdmin.Permit)
		headerActive = permit.Module
		activeId = permit.Id
	}
	return headerActive, activeId
}

//获得超级管理员具备的页面展示权限
func (this *AdminController) initAllShowPermit() {
	var leftTopId string

	// 获得当前页面的所有上级权限
	permitUpon, arrayUponId, _ := this.getNowAndAllUponPermit()

	uponIdList, _, _ := this.PermitService.FetchPermitListByUponId(arrayUponId)
	//data := this.orgPermit(uponIdList)

	permit := make(map[string]interface{})

	//Header信息列表
	permit["Header"] = uponIdList
	permit["HeaderActive"], leftTopId = this.getHeaderDefaultActive(*permitUpon)

	//左侧边栏权限列表
	permit["Left"] = this.PermitService.GetLeftPermit(leftTopId)
	this.Data["Permit"] = permit
}

func (this *AdminController) orgPermit(uponIdList *[]modelsAdmin.Permit) *map[string][]modelsAdmin.Permit {
	var result = make(map[string][]modelsAdmin.Permit)

	for _, v := range *uponIdList {
		result[v.UppermitId] = append(result[v.UppermitId], v)
	}
	return &result
}

//获得普通账号具备的账号展示权限
func (this *AdminController) getListNotSuperAdmin() []interface{} {
	item := make([]interface{}, 0)
	return item
}

//判断是否为超级管理员
func (this *AdminController) authSuperAdmin() bool {
	return true
}

func (this *AdminController) InitPageScript() {

	this.Data["PageVersion"] = "1.0"

	this.Data["CssFile"] = [...]string{
		"bootstrap.min.css",
		"font-awesome.min.css",
		"ionicons.min.css",
		"morris/morris.css",
		"jvectormap/jquery-jvectormap-1.2.2.css",
		"fullcalendar/fullcalendar.css",
		"daterangepicker/daterangepicker-bs3.css",
		"bootstrap-wysihtml5/bootstrap3-wysihtml5.min.css",
		"AdminLTE.css",
		"fileinput/fileinput.css"}

	this.Data["JsFileBefore"] = [...]string{
		"jquery.min.js",
		"jquery-ui-1.10.3.min.js",
		"bootstrap.min.js",
		"fileinput/fileinput.js",
		"fileinput/fileinput_locale_zh.js"}

	this.Data["JsFileAfter"] = [...]string{
		"raphael-min.js",
		//     'plugins/morris/morris.min.js',
		"plugins/sparkline/jquery.sparkline.min.js",
		"plugins/jvectormap/jquery-jvectormap-1.2.2.min.js",
		"plugins/jvectormap/jquery-jvectormap-world-mill-en.js",
		"plugins/fullcalendar/fullcalendar.min.js",
		"plugins/jqueryKnob/jquery.knob.js",
		"plugins/daterangepicker/daterangepicker.js",
		"plugins/bootstrap-wysihtml5/bootstrap3-wysihtml5.all.min.js",
		"plugins/iCheck/icheck.min.js",
		"AdminLTE/app.js",
		"AdminLTE/dashboard.js"}
}
func (this *AdminController) Prepare() {

	//引入父类的处理逻辑
	this.BaseController.Prepare()

	//加上权限管理
	this.InitPermitItem()

	//引入页面内容
	this.InitPageScript()

}

//设置Layout
func (this *AdminController) LoadCommon(tplName string) {

	this.Layout = "layout/main.html"
	this.TplName = tplName

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "layout/header.html"
	this.LayoutSections["SideBar"] = "layout/left.html"
	this.LayoutSections["ScriptsAfter"] = "layout/scriptsafter.html"

}
