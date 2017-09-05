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

/**
* 返回当前后台的权限列表
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) InitPermitItem() {

	if this.ConContext.IsSuperAdmin != true {

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
	var err error
	var permitModelList = make([]modelsAdmin.Permit, 0)

	permitModel := new(modelsAdmin.Permit)
	//var permitModelList []*modelsAdmin.Permit

	//获得默认的访问连接路由
	defaultController, actionString := this.DefaultControllerAndAction()

	fetchParams := make(map[string]string)
	fetchParams["Controller"], fetchParams["Action"] = this.GetControllerAndAction()
	fetchParams["Controller"] = strings.ToLower(utils.Substr(fetchParams["Controller"], 0, strings.Index(fetchParams["Controller"], "Controller")))
	fetchParams["Action"] = strings.ToLower(fetchParams["Action"])

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
func (this *AdminController) getNowAndAllUponPermit() (*[]interface{}, []string, error) {
	var err1 error
	result := make([]interface{}, 0)
	uponPermitId := make([]string, 0)

	permitModel := new(modelsAdmin.Permit)
	permitData, _ := this.getNowPermitData()
	if nil == permitData {
		//默认的上级机构必须查询
		uponPermitId = *utils.SliceUnshiftString(uponPermitId, "")
		return nil, uponPermitId, nil
	}

	uponPermitId = *utils.SliceUnshiftString(uponPermitId, permitData.Id)

	var permitModelList *[]modelsAdmin.Permit
	i := 0
	for {

		i++

		//判断如果循环超过5次还没中断，则强制中断，防止程序异常
		if "" == permitData.UppermitId || "0" == permitData.UppermitId || i > 5 {
			break
		}
		uponPermitId = *utils.SliceUnshiftString(uponPermitId, permitData.UppermitId)

		fetchParams := make(map[string]string)
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
	return &result, uponPermitId, err1
}

/**
* 获得当前地址对应的数据库存储的权限及所有上级权限
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *AdminController) getNowNotSuperAdminAndAllUponPermit(groupIds *[]string) (*[]interface{}, []string, error) {

	var err1 error

	permitModel := new(modelsAdmin.Permit)
	result := make([]interface{}, 0)

	uponPermitId := make([]string, 0)

	permitData, _ := this.getNowPermitData()
	if nil == permitData {
		//默认的上级机构必须查询
		uponPermitId = *utils.SliceUnshiftString(uponPermitId, "")
		return nil, nil, err1
	}
	uponPermitId = *utils.SliceUnshiftString(uponPermitId, permitData.Id)
	var permitModelList *[]modelsAdmin.Permit

	i := 0
	for {
		i++
		if "" == permitData.UppermitId || i > 5 {
			break
		}
		fetchParams := make(map[string]string)
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
		headerActive = permit.Mod
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

	//如果当前权限没查到,则直接跳转404
	if permitUpon == nil {
		this.Abort("404")
		return
	}
	//获得页面头部的信息
	headerPermitList, _, _ := this.PermitService.FetchPermitListByUponId([]interface{}{0})
	//Header信息列表
	//如果是超级管理员，那么权限对于此账号无效

	//Header信息列表
	permit["Header"] = headerPermitList
	permit["HeaderActive"], leftTopId = this.getHeaderDefaultActive(*permitUpon)

	var args = make(map[string]interface{})
	args["LeftTopId"] = leftTopId
	args["SuperAdmin"] = true
	leftPermit := this.PermitService.GetLeftPermit(args)

	//设置左侧权限active
	var err2 error
	leftPermit, err2 = this.setLeftActive(leftPermit, activeUponId)
	if nil != err2 {
		panic(err2)
	}

	//左侧边栏权限列表
	permit["Left"] = leftPermit
	this.Data["Permit"] = permit

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

	//如果用户组不存在，则不用继续操作了
	if len(this.ConContext.GroupIds) == 0 {
		return
	}

	// 获得当前页面的所有上级权限
	res, activeUponId, _ := this.getNowNotSuperAdminAndAllUponPermit(&this.ConContext.GroupIds)

	//如果没有查询到当前的权限
	if res == nil {
		//如果没有查询到信息，则直接跳转到未找到页面
		this.Abort("404")
		return
	}
	//如果没有查询到当前信息
	if nil == activeUponId {
		this.Data["Permit"] = permit
		return
	}
	//根据当前用户的用户组获得用户的权限
	//Header信息列表
	headerPermit, _ := groupPermit.GetGroupPermitList(this.ConContext.GroupIds, []string{"0", ""})
	permit["Header"] = headerPermit

	//处理当前头部选中的选项
	permitArray := make([]interface{}, 0)
	for _, v := range *headerPermit {
		permitArray = append(permitArray, &v)
	}

	headerActive, leftTopId := this.getHeaderDefaultActive(permitArray)

	if "" != headerActive {
		permit["HeaderActive"] = headerActive
	}

	//左侧边栏权限列表
	var args = make(map[string]interface{})
	args["SuperAdmin"] = false
	args["LeftTopId"] = leftTopId
	args["GroupIds"] = this.ConContext.GroupIds
	leftPermit := this.PermitService.GetLeftPermit(args)

	var err2 error
	leftPermit, err2 = this.setLeftActive(leftPermit, activeUponId)
	if nil != err2 {
		panic(err2)
	}
	permit["Left"] = leftPermit
	this.Data["Permit"] = permit
}

/**
* 组织左侧权限的高亮显示设置
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/16
* @return *[](map[string]interface{}), error
 */
func (this *AdminController) setLeftActive(leftPermit *[](map[string]interface{}), activeUponId []string) (*[](map[string]interface{}), error) {
	var res = make([](map[string]interface{}), 0)
	var errR error
	//	this.Debug(leftPermit)
	//如果没有数据，说明没有标明选中项
	if len(activeUponId) < 3 {
		return leftPermit, nil
	}

	//去掉第一，第二条数据为空字符串
	activeUponId = activeUponId[2:]

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
				v["ChildList"], errR = this.setLeftActive(&cList, activeUponId[1:])
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

	this.Data["SiteName"] = beego.AppConfig.String(beego.BConfig.RunMode + "::sitename")
	time := time.Now()
	y := time.Year()
	this.Data["Copyright"] = "Copyright " + strconv.Itoa(y-1) + "-" + strconv.Itoa(y) + " " + beego.AppConfig.String("appname") + " Corporation. All Rights Reserved."

	this.Data["PageTitle"] = " 后台管理中心"
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

	this.LayoutSections["HtmlHead"] = "layout/header.html"
	this.LayoutSections["HtmlSideBar"] = "layout/left.html"
	this.LayoutSections["HtmlScriptsAfter"] = "layout/scriptsafter.html"

}
