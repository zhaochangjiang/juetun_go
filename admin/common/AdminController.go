package common

import (
	"fmt"
	"juetun/common/general"
	modelsAdmin "juetun/common/models/admin"
)

type AdminController struct {
	general.BaseController
}

//返回当前后台的权限列表
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
func (this *AdminController) getNowPermitData() (*modelsAdmin.Permit, string) {
	errorMessage := ""
	permitModel := new(modelsAdmin.Permit)

	fetchParams := make(map[string]interface{})
	fetchParams["Controller"], fetchParams["Action"] = this.GetControllerAndAction()
	defaultController, actionString := this.DefaultControllerAndAction()
	if defaultController == fetchParams["Controller"] && actionString == fetchParams["Action"] {
		return permitModel, ""
	}
	var permitModelList []*modelsAdmin.Permit
	permitModelList, message := permitModel.FetchPermit(fetchParams)

	if "" != message {
		//	this.DisplayIframe(message)
		return permitModel, message
	}
	fmt.Println(permitModelList)
	if len(permitModelList) > 0 {
		permitModel = permitModelList[0]
	}
	return permitModel, errorMessage
}

//获得当前地址对应的数据库存储的权限及所有上级权限
func (this *AdminController) getNowAndAllUponPermit() (*[]interface{}, []interface{}, string) {

	permitData, message := this.getNowPermitData()
	permitModel := new(modelsAdmin.Permit)

	result := make([]interface{}, 0)
	utils := new(general.Utils)
	uponPermitId := make([]interface{}, 0)

	//默认的上级机构必须查询
	uponPermitId = *utils.Slice_unshift(uponPermitId, 0)
	for {
		if 0 == permitData.UppermitId {
			break
		}
		fetchParams := make(map[string]interface{})
		fetchParams["UppermitId"] = permitData.UppermitId
		uponPermitId = *utils.Slice_unshift(uponPermitId, permitData.UppermitId)
		permitModelList, msg := permitModel.FetchPermit(fetchParams)
		if "" != msg {
			message = msg
		}
		if nil != permitModelList {
			permitData = permitModelList[0]
			//往队列的队首添加数据
			result = *utils.Slice_unshift(result, permitData)
		}

	}
	return &result, uponPermitId, message
}

//获得超级管理员具备的页面展示权限
func (this *AdminController) initAllShowPermit() {
	//	item := make([]interface{}, 0)

	// 获得当前页面的所有上级权限
	permitUpon, arrayUponId, errorMessage := this.getNowAndAllUponPermit()

	if "" != errorMessage {
		this.DisplayIframe(errorMessage)
	}
	permitModel := new(modelsAdmin.Permit)
	uponIdList, _, msg := permitModel.FetchPermitListByUponId(arrayUponId)

	if "" != msg {
		this.DisplayIframe(msg)
	}

	permit := make(map[string]interface{})
	permit["Header"] = *uponIdList
	permit["HeaderActive"] = "dashboard"
	permit["Left"] = permitUpon

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
