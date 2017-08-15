package common

import (
	"errors"
	"juetun/common/general"
	modelsAdmin "juetun/common/models/admin"
	"juetun/common/utils"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type AdminController struct {
	PermitService *modelsAdmin.Permit
	general.BaseController
	NotNeedLogin bool //是否需要登录
	isSuperAdmin bool //是否为超级管理员
}

/**
* 返回当前后台的权限列表
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) InitPermitItem() {

	if !this.authSuperAdmin() {

		//如果不是超级管理员
		this.initAllShowNotSuperAdminPermit()

	} else {

		//如果是超级管理员
		this.initAllShowSuperAdminPermit()
	}

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
	permitModel := new(modelsAdmin.Permit)
	fetchParams := make(map[string]string)
	fetchParams["Controller"], fetchParams["Action"] = this.GetControllerAndAction()

	fetchParams["Controller"] = strings.ToLower(utils.Substr(fetchParams["Controller"], 0, strings.Index(fetchParams["Controller"], "Controller")))

	fetchParams["Action"] = strings.ToLower(fetchParams["Action"])

	//获得默认的访问连接路由
	defaultController, actionString := this.DefaultControllerAndAction()

	//如果
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

/**
* 获得当前地址对应的数据库存储的权限及所有上级权限
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) getNowAndAllUponPermit() (*[]interface{}, []interface{}, error) {

	permitModel := new(modelsAdmin.Permit)

	result := make([]interface{}, 0)

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
		fetchParams := make(map[string]string)
		fetchParams["id"] = permitData.UppermitId
		uponPermitId = *utils.Slice_unshift(uponPermitId, permitData.UppermitId)
		permitModelList, err1 = permitModel.FetchPermit(fetchParams)
		if len(permitModelList) <= 0 {
			break
		}
		permitData = (permitModelList[0])
		//往队列的队首添加数据
		result = *utils.Slice_unshift(result, permitData)

	}
	return &result, uponPermitId, err1
}

/**
* 获得当前地址对应的数据库存储的权限及所有上级权限
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) getNowNotSuperAdminAndAllUponPermit(groupIds *[]string) (*[]interface{}, []interface{}, error) {

	permitModel := new(modelsAdmin.Permit)
	log.Println(*groupIds)
	result := make([]interface{}, 0)

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
		fetchParams := make(map[string]string)
		fetchParams["id"] = permitData.UppermitId
		uponPermitId = *utils.Slice_unshift(uponPermitId, permitData.UppermitId)
		permitModelList, err1 = permitModel.FetchPermit(fetchParams)
		if len(permitModelList) <= 0 {
			break
		}
		permitData = (permitModelList[0])
		//往队列的队首添加数据
		result = *utils.Slice_unshift(result, permitData)

	}
	return &result, uponPermitId, err1
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
		headerActive = permit.Module
		activeId = permit.Id
	}
	return headerActive, activeId
}

/**
*获得超级管理员具备的页面展示权限
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
* @param isSuperAdmin 是否为超级管理员
*
 */
func (this *AdminController) initAllShowSuperAdminPermit() {
	var leftTopId string
	permit := make(map[string]interface{})
	// 获得当前页面的所有上级权限
	permitUpon, activeUponId, _ := this.getNowAndAllUponPermit()

	//获得页面头部的信息
	headerPermitList, _, _ := this.PermitService.FetchPermitListByUponId([]interface{}{0})
	//Header信息列表
	//如果是超级管理员，那么权限对于此账号无效

	//Header信息列表
	permit["Header"] = headerPermitList
	permit["HeaderActive"], leftTopId = this.getHeaderDefaultActive(*permitUpon)

	//左侧边栏权限列表
	permit["Left"] = this.PermitService.GetLeftPermit(leftTopId)
	this.Data["Permit"] = permit

	log.Println(activeUponId)
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
* 处理当前非超级管理员的权限
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) initAllShowNotSuperAdminPermit() {

	//用户组权限ID列表
	var groupPermit modelsAdmin.GroupPermit
	permit := make(map[string]interface{})

	//var headerPermitList *[]modelsAdmin.GroupPermit
	//获得当前用户的用户组ID列表
	groupIds := this.getNowUserGroupId()

	//如果用户组不存在，则不用继续操作了
	if len(groupIds) == 0 {
		return
	}

	//根据当前用户的用户组获得用户的权限
	//Header信息列表
	headerPermit, _ := groupPermit.GetGroupPermitList(groupIds, []string{"0", ""})
	permit["Header"] = headerPermit

	// 获得当前页面的所有上级权限
	permitUpon, activeUponId, _ := this.getNowNotSuperAdminAndAllUponPermit(&groupIds)

	permitArray := make([]interface{}, 0)
	for _, v := range *headerPermit {
		permitArray = append(permitArray, &v)
	}

	headerActive, leftTopId := this.getHeaderDefaultActive(permitArray)

	if "" != headerActive {
		permit["HeaderActive"] = headerActive
	}

	//左侧边栏权限列表
	permit["Left"] = this.PermitService.GetLeftPermitByGroupId(leftTopId, groupIds)
	this.Data["Permit"] = permit
	log.Println(activeUponId)
	log.Println(permitUpon)
}

/**
* 头部页面的CSS JS文件引入
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
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

/**
* 上下文调用方法，用作引入初始化用
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) Prepare() {

	//引入父类的处理逻辑
	this.BaseController.Prepare()

	//TODO此为由于SESSION保持有问题的临时解决办法
	//	log.Println("AdminController 设置的临时解决登录的方法!")
	//	this.SetSession("Uid", "1")
	//	this.SetSession("Username", "长江")
	//	this.SetSession("Avater", "/assets/img/avatar5.jpg")
	if true == this.NotNeedLogin {
		return
	}
	//判断是否登录
	//如果需要登录等于
	if this.IsLogin() == false {
		gotoUrl := general.CreateUrl("passport", "login", make(map[string]string), "web")
		this.Redirect(gotoUrl, 301)
		return
	}

	this.Data["SiteName"] = beego.AppConfig.String("sitename")
	time := time.Now()
	y := time.Year()
	this.Data["Copyright"] = "Copyright " + strconv.Itoa(y-1) + "-" + strconv.Itoa(y) + " " + beego.AppConfig.String("appname") + " Corporation. All Rights Reserved."

	//处理用户默认信息，比如头像
	this.UserDataDefault()

	this.Data["PageTitle"] = " 后台管理中心"
	//设置登录信息

	this.Data["NowHourAndMinute"] = strconv.Itoa(time.Hour()) + ":" + strconv.Itoa(time.Minute())
	//加上权限管理
	this.InitPermitItem()

	//引入页面内容
	this.InitPageScript()

}

/**
* 判断是否为超级管理员
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) authSuperAdmin() bool {
	this.isSuperAdmin = false
	if nil != this.GetSession("SuperAdmin") {
		switch this.GetSession("SuperAdmin").(string) {
		case "yes":
			this.isSuperAdmin = true
		default:
			this.isSuperAdmin = false
		}
	}
	//初始化是否为超级管理员的配置
	this.Data["isSuperAdmin"] = this.isSuperAdmin
	return this.isSuperAdmin
}

/**
* 判断是否登录
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) IsLogin() bool {
	uid := this.GetSession("Uid")
	log.Println("---------------IsLogin---------------------")
	log.Println(uid)
	log.Println("---------------IsLogin O---------------------")

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

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "layout/header.html"
	this.LayoutSections["SideBar"] = "layout/left.html"
	this.LayoutSections["ScriptsAfter"] = "layout/scriptsafter.html"

}
