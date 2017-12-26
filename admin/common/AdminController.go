package common

import (
	"errors"
	"juetun/common/general"
	modelsAdmin "juetun/common/models/admin"
	"juetun/common/utils"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type AdminController struct {
	PermitService *modelsAdmin.Permit
	general.BaseController
	ConContext general.ControllerContext //本次请求上下文存储的数据
}

func (this *AdminController) Edit() {
	this.LoadCommon("layout/edit.html")
}

func (this *AdminController) Delete() {

}

func (this *AdminController) List() {
	this.Debug("AdminController_List")
	this.LoadCommon("layout/list.html")
}

func (this *AdminController) SetListPageMessage() {

	if this.ConContext.IncludePageProperty.HaveTable {
		this.ConContext.JsFileAfter = append(this.ConContext.JsFileAfter, []string{"plugins/datatables/jquery.dataTables.js", "plugins/datatables/dataTables.bootstrap.js"}...)
		this.ConContext.CssFile = append(this.ConContext.CssFile, []string{"datatables/dataTables.bootstrap.css"}...)
	}
	//如果页面属性中有checkbox
	if this.ConContext.IncludePageProperty.HaveCheckbox {

		this.ConContext.JsFileAfter = append(this.ConContext.JsFileAfter, []string{"plugins/iCheck/icheck.min.js"}...)
		this.ConContext.CssFile = append(this.ConContext.CssFile, []string{"iCheck/minimal/blue.css"}...)

	}

}

/**
* 返回当前后台的权限列表
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) InitPermitItem() {

	var permit = make(map[string]interface{})
	var headerPermit = new([]modelsAdmin.PermitAdmin)
	var leftPermit = new([](map[string]interface{}))
	var permitUpon = new([]interface{})
	var activeUponId = make([]string, 0)
	var leftTopId string
	var headerActive string
	var getLeftPermitArgs = make(map[string]interface{})

	var nowChildList = new([]modelsAdmin.PermitAdmin)

	if this.ConContext.IsSuperAdmin != true { //如果不是超级管理员

		//用户组权限ID列表
		var groupPermit modelsAdmin.GroupPermit

		// 获得当前页面的所有上级权限
		permitUpon, activeUponId, nowChildList, _ = this.getNowNotSuperAdminAndAllUponPermit(&this.ConContext.GroupIds)

		//如果没有查询到当前信息
		if len(activeUponId) <= 0 {
			this.Data["Permit"] = permit
			return
		}
		//根据当前用户的用户组获得用户的权限
		//Header信息列表
		headerPermit, _ = groupPermit.GetGroupPermitList(this.ConContext.GroupIds, []string{""})

		//处理当前头部选中的选项
		permitArray := make([]interface{}, 0)
		for _, v := range *headerPermit {
			permitArray = append(permitArray, &v)
		}

		headerActive, leftTopId = this.getHeaderDefaultActive(permitArray)

		//左侧边栏权限查询参数
		getLeftPermitArgs["SuperAdmin"] = false
		getLeftPermitArgs["LeftTopId"] = leftTopId
		getLeftPermitArgs["GroupIds"] = this.ConContext.GroupIds

	} else { //如果是超级管理员
		this.Debug("kkk")
		// 获得当前页面的所有上级权限
		permitUpon, activeUponId, nowChildList, _ = this.getNowAndAllUponPermit()
		this.Debug("bbb")
		//如果当前权限没查到,则直接跳转404
		if len(activeUponId) <= 0 {
			this.Data["Permit"] = permit
			return
		}

		//获得页面头部的信息
		headerPermit, _, _ = this.PermitService.FetchPermitListByUponId(&[]string{""})

		//Header信息列表
		//如果是超级管理员，那么权限对于此账号无效

		headerActive, leftTopId = this.getHeaderDefaultActive(*permitUpon)

		getLeftPermitArgs["LeftTopId"] = leftTopId
		getLeftPermitArgs["SuperAdmin"] = true

	}

	//获得当前权限的子权限列表(此处值设置为int为了后期扩展用)
	var permitShow = make(map[string]int)
	for _, v := range *nowChildList {
		keys := v.Mod + "_" + v.Controller + "_" + v.Action
		permitShow[keys] = 1
	}
	this.Data["PermitShow"] = permitShow
	leftPermit = this.PermitService.GetLeftPermit(getLeftPermitArgs)
	leftPermit, _ = this.setLeftActive(leftPermit, activeUponId, false)

	//Header信息列表
	permit["Header"] = this.setHeaderActive(headerPermit, headerActive)

	//左侧边栏权限列表
	permit["Left"] = leftPermit
	this.ConContext.Permit = &permit
	this.Data["Permit"] = permit
}

func (this *AdminController) setHeaderActive(headerPermit *[]modelsAdmin.PermitAdmin, headerActiveId string) *[]modelsAdmin.PermitAdmin {

	for k, v := range *headerPermit {

		if headerActiveId == v.Mod {
			(*headerPermit)[k].Active = true
		}
	}
	return headerPermit
}

/**
* 默认访问的页面
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) DefaultControllerAndAction() (string, string) {
	return "main", "get"
}

/**
* 获得当前的权限
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) getNowPermitData() (*modelsAdmin.Permit, error) {
	var err error
	var permitModelList = make([]modelsAdmin.Permit, 0)

	permitModel := new(modelsAdmin.Permit)
	//var permitModelList []*modelsAdmin.Permit

	//获得默认的访问连接路由
	defaultController, actionString := this.DefaultControllerAndAction()

	fetchParams := make(map[string]interface{})

	fetchParams["Controller"] = this.ConContext.Controller
	fetchParams["Action"] = this.ConContext.Action

	//如果
	if defaultController == fetchParams["Controller"] && actionString == fetchParams["Action"] {
		return permitModel, errors.New("the controller and action is equal default!")
	}

	//如果是超级管理员
	if this.ConContext.IsSuperAdmin == true {
		var listTmp *[]modelsAdmin.Permit
		listTmp, err = this.PermitService.FetchPermit(fetchParams)
		permitModelList = *listTmp
	} else {
		//如果不是超级管理员
		var pList *[]modelsAdmin.Permit
		pList, err = this.PermitService.FetchPermitByGroupId(this.ConContext.GroupIds, fetchParams)
		if 0 == len(*pList) {
			return nil, nil
		}
		for _, v := range *pList {
			permitModelList = append(permitModelList, v)
		}
	}

	if len(permitModelList) > 0 {
		permitModel = &permitModelList[0]
	}
	return permitModel, err
}

/**
* 获得当前地址对应的数据库存储的权限及所有上级权限
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
* @return *[]interface{}, []interface{}, error
 */
func (this *AdminController) getNowAndAllUponPermit() (*[]interface{}, []string, *[]modelsAdmin.PermitAdmin, error) {
	var err1 error
	var result = make([]interface{}, 0)
	var uponPermitId = make([]string, 0)
	var nowChildList = new([]modelsAdmin.PermitAdmin)

	permitModel := new(modelsAdmin.Permit)
	//获得当前的权限
	permitData, _ := this.getNowPermitData()
	if nil == permitData {
		//默认的上级机构必须查询
		uponPermitId = *utils.SliceUnshiftString(uponPermitId, "")
		return &result, uponPermitId, nowChildList, nil
	}

	//获得当前页面的子权限
	var ids = []string{permitData.Id}
	nowChildList = permitModel.GetAllChildListByPids(&ids, this.ConContext.IsSuperAdmin, &this.ConContext.GroupIds)

	uponPermitId = *utils.SliceUnshiftString(uponPermitId, permitData.Id)

	var permitModelList *[]modelsAdmin.Permit
	i := 0
	for {

		i++

		//判断如果循环超过5次还没中断，则强制中断，防止程序异常
		if "" == permitData.UppermitId || "0" == permitData.UppermitId || i > 6 {
			break
		}
		uponPermitId = *utils.SliceUnshiftString(uponPermitId, permitData.UppermitId)

		fetchParams := make(map[string]interface{})
		fetchParams["id"] = permitData.UppermitId
		permitModelList, err1 = permitModel.FetchPermit(fetchParams)

		//如果没有查询到数据，那么跳出循环
		if len(*permitModelList) <= 0 {
			break
		}
		permitData = &((*permitModelList)[0])

		params := make(map[string]string)
		permitDataTmp := permitModel.OrgAdminPermit(*permitData, params)

		//往队列的队首添加数据
		result = *utils.SliceUnshift(result, permitDataTmp)

	}

	//默认的上级机构必须查询
	uponPermitId = *utils.SliceUnshiftString(uponPermitId, "")
	return &result, uponPermitId, nowChildList, err1
}

/**
* 获得当前地址对应的数据库存储的权限及所有上级权限
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) getNowNotSuperAdminAndAllUponPermit(groupIds *[]string) (*[]interface{}, []string, *[]modelsAdmin.PermitAdmin, error) {

	var err1 error = nil
	var uponPermitId = make([]string, 0)
	var result = make([]interface{}, 0)
	var nowChildList = new([]modelsAdmin.PermitAdmin)

	if len(*groupIds) <= 0 {
		return &result, uponPermitId, nowChildList, err1
	}

	permitModel := new(modelsAdmin.Permit)

	permitData, _ := this.getNowPermitData()
	if nil == permitData {
		//默认的上级机构必须查询
		uponPermitId = *utils.SliceUnshiftString(uponPermitId, "")
		return nil, nil, nowChildList, err1
	}

	//获得当前权限的所有子权限
	var ids = []string{permitData.Id}
	nowChildList = permitModel.GetAllChildListByPids(&ids, this.ConContext.IsSuperAdmin, &this.ConContext.GroupIds)

	uponPermitId = *utils.SliceUnshiftString(uponPermitId, permitData.Id)
	var permitModelList *[]modelsAdmin.Permit

	i := 0
	for {
		i++
		if "" == permitData.UppermitId || i > 5 {
			break
		}
		fetchParams := make(map[string]interface{})
		fetchParams["id"] = permitData.UppermitId
		uponPermitId = *utils.SliceUnshiftString(uponPermitId, permitData.UppermitId)

		permitModelList, err1 = permitModel.FetchPermit(fetchParams)
		if len(*permitModelList) <= 0 {
			break
		}
		permitData = &((*permitModelList)[0])
		//往队列的队首添加数据
		result = *utils.SliceUnshift(result, permitData)

	}

	//默认的上级机构必须查询
	uponPermitId = *utils.SliceUnshiftString(uponPermitId, "")
	return &result, uponPermitId, nowChildList, err1
}

/**
* 获得header默认的Type
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) getHeaderDefaultActive(permitUpon []interface{}) (string, string) {
	var activeId string
	var headerActive = "dashboard"

	length := len(permitUpon)
	if length > 0 {
		permit := permitUpon[0].(*modelsAdmin.PermitAdmin)
		headerActive = permit.Mod
		activeId = permit.Id
	}
	return headerActive, activeId
}

/**
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) orgPermit(uponIdList *[]modelsAdmin.Permit) *map[string][]modelsAdmin.Permit {
	var result = make(map[string][]modelsAdmin.Permit)

	for _, v := range *uponIdList {
		result[v.UppermitId] = append(result[v.UppermitId], v)
	}
	return &result
}

/**
*获得当前用户的用户组
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) getNowUserGroupId() []string {
	var groupIds []string
	var uid = this.GetSession("Uid")
	if nil == uid {
		return groupIds
	}

	groupuser := new(modelsAdmin.Groupuser)

	//获得当前用户的用户组列表
	getGoupList, err := groupuser.GetGoupList(uid.(string))
	if nil != err {
		panic(err)
	}

	for _, v := range *getGoupList {
		groupIds = append(groupIds, v.GroupId)
	}
	return groupIds
}

/**
* 组织左侧权限的高亮显示设置
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/16
* @return *[](map[string]interface{}), error
 */
func (this *AdminController) setLeftActive(leftPermit *[](map[string]interface{}), activeUponId []string, needDevid bool) (*[](map[string]interface{}), error) {
	var res = make([](map[string]interface{}), 0)
	var errR error

	if len(activeUponId) == 0 {
		return leftPermit, nil
	}
	if needDevid != true {
		//	this.Debug(leftPermit)
		//如果没有数据，说明没有标明选中项
		if len(activeUponId) < 3 {
			return leftPermit, nil
		}

		//去掉第一，第二条数据为空字符串
		activeUponId = activeUponId[2:]
	}

	for _, v := range *leftPermit {

		//将数据转换为Permit格式
		p := v["Permit"].(modelsAdmin.PermitAdmin)

		//如果ID相等
		if p.Id == activeUponId[0] {

			//设置当前为高亮选中
			v["Active"] = true
			childList := v["ChildList"].(*[]map[string]interface{})
			cList := make([]map[string]interface{}, 0)

			for _, v := range *childList {
				cList = append(cList, v)
			}

			//如果有子选项
			if len(cList) > 0 {

				//此处为一个递归处理
				v["ChildList"], errR = this.setLeftActive(&cList, activeUponId[1:], true)
				if nil != errR {
					panic(errR)
				}
			}

		}

		res = append(res, v)
	}

	return &res, nil
}

/**
* 头部页面的CSS JS文件引入
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) InitPageScript() {

	this.Data["PageVersion"] = "1.0"

	this.Data["CssFile"] = []string{
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

	this.Data["JsFileBefore"] = []string{
		"jquery.min.js",
		"jquery-ui-1.10.3.min.js",
		"bootstrap.min.js",
		"fileinput/fileinput.js",
		"fileinput/fileinput_locale_zh.js",
		"base.js"}

	this.Data["JsFileAfter"] = []string{
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
		"AdminLTE/dashboard.js",
	}

}

func (this *AdminController) initConContextControllerAndAction() {
	con, act := this.GetControllerAndAction()
	con = strings.ToLower(con)
	con = utils.TrimRight(con, "controller")
	act = strings.ToLower(act)
	this.ConContext.Controller = con
	this.ConContext.Action = act
}

/**
* 上下文调用方法，用作引入初始化用
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) Prepare() {

	//初始化一些必要的参数
	this.initConContextControllerAndAction()
	this.ConContext.NeedRenderJs = false

	this.ConContext.OutputResult.Code = 0
	this.ConContext.OutputResult.Message = "操作成功"
	this.ConContext.OutputResult.Data = nil

	siteName := beego.AppConfig.String(beego.BConfig.RunMode + "::sitename")
	this.Data["SiteName"] = siteName
	time := time.Now()
	this.Data["Copyright"] = "Copyright " + strconv.Itoa(time.Year()-1) + "-" + strconv.Itoa(time.Year()) + " " + beego.AppConfig.String("appname") + " Corporation. All Rights Reserved."
	this.Data["PageTitle"] = siteName
	this.Data["DashboardTitle"] = "-后台管理中心[" + siteName + "]"
	//设置登录信息
	this.Data["NowHourAndMinute"] = strconv.Itoa(time.Hour()) + ":" + strconv.Itoa(time.Minute())

	//引入页面内容
	this.InitPageScript()
	//初始化模板引入
	this.LayoutSections = make(map[string]string)

	if true == this.ConContext.NotNeedLogin {

		//引入父类的处理逻辑
		this.BaseController.Prepare()
		return
	}

	//判断是否登录
	//如果需要登录等于
	if this.IsLogin() == false {
		gotoUrl := general.CreateUrl("passport", "login", make(map[string]string), "web")
		this.Redirect(gotoUrl, 301)
		return
	}

	//初始化当前用户是否为超级管理员
	if !this.authSuperAdmin() {
		//获得当前用户的用户组ID列表
		this.ConContext.GroupIds = this.getNowUserGroupId()
	}

	//处理用户默认信息，比如头像
	this.UserDataDefault()
	if this.ConContext.NotNeedValidatePermit == false {
		//加上权限管理
		this.InitPermitItem()
	}
	//引入父类的处理逻辑
	this.BaseController.Prepare()
	return
}

/**
* 判断是否为超级管理员
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) authSuperAdmin() bool {
	this.ConContext.IsSuperAdmin = false

	if nil != this.GetSession("SuperAdmin") {
		switch this.GetSession("SuperAdmin").(string) {
		case "yes":
			this.ConContext.IsSuperAdmin = true
		default:
			this.ConContext.IsSuperAdmin = false
		}
	} else {
		panic(errors.New("the session message is not having SuperAdmin!"))
	}
	//初始化是否为超级管理员的配置
	this.Data["isSuperAdmin"] = this.ConContext.IsSuperAdmin
	return this.ConContext.IsSuperAdmin
}

/**
* 判断是否登录
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) IsLogin() bool {
	uid := this.GetSession("Uid")
	if nil == uid {
		return false
	}

	switch uid.(type) {
	case string:
		if "" == uid.(string) {
			return false
		} else {
			return true
		}
	case int:
		if 0 == uid.(int) {
			return false
		} else {
			return true
		}
	default:
		panic("Uid is not Int or Strig")
	}

}

/**
* 设置用于页面排版Layout
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) LoadCommon(tplName string) {

	this.Layout = "layout/main.html"
	this.TplName = tplName
	this.SetListPageMessage()
	this.LayoutSections["HtmlHead"] = "layout/header.html"
	this.LayoutSections["HtmlSideBar"] = "layout/left.html"
	this.LayoutSections["HtmlScriptsAfter"] = "layout/scriptsafter.html"

	this.Data["PageVersion"] = "1.0"

	CssFile := []string{
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

	jsFileBefore := []string{
		"jquery.min.js",
		"jquery-ui-1.10.3.min.js",
		"bootstrap.min.js",
		"fileinput/fileinput.js",
		"fileinput/fileinput_locale_zh.js", "base.js"}

	jsFileAfter := []string{
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
		"AdminLTE/dashboard.js",
	}
	otherJs := "themes/" + this.ConContext.Controller + "/" + this.ConContext.Action + ".js"

	for _, v := range this.ConContext.JsFileBefore {
		jsFileBefore = append(jsFileBefore, v)
	}
	for _, v := range this.ConContext.CssFile {
		CssFile = append(CssFile, v)
	}
	for _, v := range this.ConContext.JsFileAfter {
		jsFileAfter = append(jsFileAfter, v)
	}

	if this.ConContext.NeedRenderJs == true {
		jsFileAfter = append(jsFileAfter, otherJs)

	}

	this.ConContext.JsFileBefore = jsFileBefore
	this.ConContext.JsFileAfter = jsFileAfter
	this.ConContext.CssFile = CssFile

	this.Data["JsFileAfter"] = this.ConContext.JsFileAfter
	this.Data["JsFileBefore"] = this.ConContext.JsFileBefore
	this.Data["CssFile"] = this.ConContext.CssFile

	this.initBreadcrumbParams()
	this.InitBreadcrumb()
}
func (this *AdminController) initBreadcrumbParams() {

	//设置面包屑
	this.ConContext.Breadcrumbs = []general.Breadcrumb{general.Breadcrumb{"/", "fa-dashboard", "主页", false}}
	if this.ConContext.Permit != nil {

		if _, ok := (*this.ConContext.Permit)["Header"]; ok {
			for _, v := range *(*this.ConContext.Permit)["Header"].(*[]modelsAdmin.PermitAdmin) {
				if v.Active == true {
					bc := new(general.Breadcrumb)
					bc.Name = v.Name
					bc.Href = v.UrlString
					bc.FaCss = v.Csscode
					bc.Active = false
					this.ConContext.Breadcrumbs = append(this.ConContext.Breadcrumbs, *bc)
				}

			}
		}

		if _, ok := (*this.ConContext.Permit)["Left"]; ok {
			for _, tmp1 := range *(*this.ConContext.Permit)["Left"].(*[]map[string]interface{}) {
				if tmp1["Active"].(bool) == true {
					tmp2Permit := tmp1["Permit"].(modelsAdmin.PermitAdmin)
					bc := new(general.Breadcrumb)
					bc.Name = tmp2Permit.Name
					bc.Href = tmp2Permit.UrlString
					bc.FaCss = tmp2Permit.Csscode
					bc.Active = false
					this.ConContext.Breadcrumbs = append(this.ConContext.Breadcrumbs, *bc)
				}
				childList := tmp1["ChildList"].(*[]map[string]interface{})
				if len(*childList) > 0 {
					for _, tmp2 := range *childList {
						if tmp2["Active"].(bool) == true {
							tmp2Permit := tmp2["Permit"].(modelsAdmin.PermitAdmin)
							bc := new(general.Breadcrumb)
							bc.Name = tmp2Permit.Name
							bc.Href = tmp2Permit.UrlString
							bc.FaCss = tmp2Permit.Csscode
							bc.Active = false
							this.ConContext.Breadcrumbs = append(this.ConContext.Breadcrumbs, *bc)
						}
					}
				}
			}
		}
	}

}
func (this *AdminController) InitBreadcrumb() {

	this.Data["Breadcrumbs"] = ""
	if len(this.ConContext.Breadcrumbs) <= 0 {
		return
	}
	var s = "<ol class=\"breadcrumb\">"
	for _, v := range this.ConContext.Breadcrumbs {
		if v.Href == "" {
			v.Href = "javascript:void(0);"
		}
		var active, fa string = "", ""
		if v.Active == true {
			active = "class =\"active\""
		}
		if v.FaCss != "" {
			fa = "<i class=\"fa " + v.FaCss + "\"></i>"
		}
		s += "<li " + active + "><a href=\"" + v.Href + "\">" + fa + v.Name + "</a></li>"
	}

	this.Data["Breadcrumbs"] = s + "</ol>"
	return
}
