package common

import (
	"errors"
	"fmt"
	modelsAdmin "juetun/common/models/admin"

	serviceAdmin "juetun/admin/service"
	"juetun/common/general"
)

type AdminController struct {
	general.BaseController
	permitService permitService
}

//返回当前后台的权限列表
//@return void
func (this *AdminController) InitPermitItem() {

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
//@return *modelsAdmin.Permit, error
func (this *AdminController) getNowPermitData() (*modelsAdmin.Permit, error) {

	var permitModel modelsAdmin.Permit
	fetchParams := make(map[string]interface{})
	fetchParams["Controller"], fetchParams["Action"] = this.GetControllerAndAction()

	defaultController, actionString := this.DefaultControllerAndAction()

	if defaultController == fetchParams["Controller"] && actionString == fetchParams["Action"] {
		return &permitModel, errors.New("")
	}
	permitModelList, err := this.permitService.FetchPermit(fetchParams)
	permitList := *permitModelList
	if len(permitList) > 0 {
		permitModel = permitList[0]
	}
	return &permitModel, err
}

//获得当前地址对应的数据库存储的权限及所有上级权限
//@return *[]*modelsAdmin.Permit, []interface{}, error
func (this *AdminController) getNowAndAllUponPermit() (*[]*modelsAdmin.Permit, []interface{}, error) {

	result := make([]*modelsAdmin.Permit, 0)

	permitData, _ := this.getNowPermitData()

	utils := new(general.Utils)

	//默认的上级机构必须查询
	uponPermitId := make([]interface{}, 0)
	uponPermitId = *utils.Slice_unshift(uponPermitId, 0)

	var permitModelList *[]modelsAdmin.Permit
	i := 0

	for {
		i++
		//防止死循环，跳不出功能
		if i > 20 || 0 == permitData.UppermitId {
			break
		}

		//将数据添加到数组的队首
		uponPermitId = *utils.Slice_unshift(uponPermitId, permitData.UppermitId)

		fetchParams := make(map[string]interface{})
		fetchParams["Id"] = permitData.UppermitId
		permitModelList, _ = this.permitService.FetchPermit(fetchParams)

		//如果数据库查询的结果集不为空
		if len(*permitModelList) > 0 {
			permitData = &((*permitModelList)[0])
			//往队列的队首添加数据
			slice := []*modelsAdmin.Permit{(permitData)}
			result = append(slice, result...)
		}

	}
	return &result, uponPermitId, errors.New("")
}

//在终端输出数据 用于调试数据
func (this *AdminController) output(p interface{}) {

	fmt.Println("**********echo params start ***********")
	fmt.Println(p)
	fmt.Println("------------echo params  over ---------")
}

//获得header默认的Type
//@return string
func (this *AdminController) getHeaderDefaultActive(permitUpon []*modelsAdmin.Permit) string {
	headerActive := "dashboard" //默认的选中地址
	length := len(permitUpon)
	if length > 0 {
		headerActive = (*(permitUpon[0])).Module
	}
	return headerActive
}

//获得超级管理员具备的页面展示权限
//@return void
func (this *AdminController) initAllShowPermit() {

	//初始化权限Service
	this.permitService = new(serviceAdmin.PermitService)

	//	item := make([]interface{}, 0)

	// 获得当前页面的所有上级权限
	permitUpon, arrayUponId, _ := this.getNowAndAllUponPermit()

	//查询所有的上级ID的下级权限列表
	uponIdList, _, _ := this.permitService.FetchPermitListByUponId(arrayUponId)

	permit := make(map[string]interface{})

	//后台界面header信息
	permit["HeaderActive"] = this.getHeaderDefaultActive(*permitUpon)
	permit["Header"] = *uponIdList

	permit["Left"] = this.permitService.OrgPermitLeftData(permitUpon)

	this.Data["Permit"] = permit

	//        $permitIdArray = $this->getNowPermitLink($permit);

	//        if (!empty($permit['id'])) {
	//            $permitData ['childPermit'] = $this->findAll(array(
	//                'uppermit_id' => array(
	//                    'doType' => 'in',
	//                    'value' => $permit['id'])));
	//        }

	//        $permitData ['header'] = $this->findAll(array(
	//            'uppermit_id' => array(
	//                'doType' => 'in',
	//                'value' => 0)), '', '`obyid` asc');
	//        $headerActive = array_shift($permitIdArray);
	//        $uppermitIdArray = array();
	//        foreach ($permitData ['header'] as $key => $value) {
	//            if (($value['id'] == $headerActive['id'])) {
	//                $permitData['header'][$key]['active'] = true;
	//                $uppermitIdArray[] = $value['id'];
	//            } else {
	//                $permitData['header'][$key]['active'] = false;
	//            }
	//        }

	//        $uppermitIdData = array();

	//        //   stop($permitIdArray);

	//        $i = 0;
	//        while (true) {
	//            $temp = $this->findAll(array(
	//                'uppermit_id' => array(
	//                    'doType' => 'in',
	//                    'value' => $uppermitIdArray)), '', '`obyid` asc');
	//            $uppermitIdArray = array();
	//            $permitList = array();

	//            foreach ($temp as $value) {
	//                $permitList[$value['uppermit_id']][] = $value;
	//                $uppermitIdArray[] = $value['id'];
	//            }
	//            $uppermitIdData[] = $permitList;
	//            if ($i > 1) {
	//                break;
	//            }
	//            $i++;
	//        }

	//        $permitData['left'] = $this->organizationPermit($uppermitIdData, $permitIdArray);
	//        //    stop($permitData['left']);
	//        return $permitData;

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
