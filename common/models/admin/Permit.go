package admin

import (
	"github.com/astaxie/beego/orm"
)

type Permit struct {
	Id         string `orm:"column(id);pk" json:"id"`
	Name       string `orm:varchar(50);orm:"column(name)"`
	Module     string `orm:varchar(30);orm:"column(module)"`
	Controller string `orm:varchar(30);orm:"column(controller)"`
	Action     string `orm:varchar(30);orm:"column(action)"`
	UppermitId string `orm:int(10);orm:"column(uppermit_id)"`
	DomainMap  string `orm:"column(domain_map)"`
	Obyid      string `orm:"column(obyid)"`
	Csscode    string `orm:varchar(500);orm:"column(csscode)"`
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Permit))
}
func (this *Permit) TableName() string {
	return "permit"
}
func (this *Permit) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}

func (this *Permit) getQuerySeter() orm.QuerySeter {
	return this.getOrm().QueryTable(this)
}

//根据上级权限，查询所有下级权限
func (this *Permit) FetchPermitListByUponId(uponid []interface{}) (*[]PermitAdmin, int64, error) {
	var permitList []Permit
	var querySeter orm.QuerySeter

	querySeter = this.getQuerySeter().Filter("uppermit_id__in", uponid...).OrderBy("-id")
	num, err := querySeter.All(&permitList)

	result := make([]PermitAdmin, 0)
	for _, v := range permitList {
		params := make(map[string]string)
		params["module"] = v.Module
		result = append(result, *this.orgAdminPermit(v, params))
	}
	return &result, num, err
}

//查询单个权限
func (this *Permit) FetchPermit(argument map[string]interface{}) ([]*Permit, error) {

	var permitList []*Permit
	querySeter := this.getQuerySeter()
	for k, v := range argument {
		querySeter = querySeter.Filter(k, v)
	}
	_, err := querySeter.All(&permitList)

	return permitList, err
}

//删除权限
func (this *Permit) DeletePermit(permitIds []string) (bool, error) {

	//删除相关的数据
	groupPermit := new(GroupPermit)
	_, err1 := groupPermit.DeleteByPermitIds(permitIds)
	if nil != err1 {
		return false, err1
	}

	//删除表头信息
	_, err := this.getQuerySeter().Filter("id__in", permitIds).Delete()
	if nil != err {
		return false, err
	}
	return true, err
}

//左侧权限结构体
type PermitLeft struct {
	Permit
	ChildList []Permit
	Active    bool
}

//后台结构体
type PermitAdmin struct {
	Permit
	Params map[string]string
	Domain string
}

//
func (this *Permit) FetchDefaultPermitByModuleString(moduleString string) {
	this.Controller = "data"
	this.Action = "list"
	this.DomainMap = ""
}

//获得左边的权限列表
func (this *Permit) GetLeftPermit(leftTopId string) *[](map[string]interface{}) {
	result := make([]map[string]interface{}, 0)
	if leftTopId == "" {
		return &result
	}

	var permitList []Permit
	var querySeter orm.QuerySeter
	var childPermitList []Permit

	//查询上级权限为leftTopId的权限列表
	querySeter = this.getQuerySeter().Filter("uppermit_id__exact", leftTopId).OrderBy("obyid")
	querySeter.All(&permitList)

	leftPermitIdList := make([]string, 0)
	for _, v := range permitList {
		leftPermitIdList = append(leftPermitIdList, v.Id)
	}

	this.getQuerySeter().Filter("uppermit_id__in", leftPermitIdList).OrderBy("obyid").All(&childPermitList)

	childPermit := make(map[string][]PermitAdmin)
	for _, v := range childPermitList {
		params := make(map[string]string)
		childPermit[v.UppermitId] = append(childPermit[v.UppermitId], *(this.orgAdminPermit(v, params)))
	}

	for _, v := range permitList {

		everyData := make(map[string]interface{})
		everyData["Permit"] = v
		everyData["Active"] = false //默认设置不为选中的状态
		everyData["ChildList"] = make([]PermitAdmin, 0)

		//判断内容是否存在，相当于PHP中的isset函数
		if _, ok := childPermit[v.Id]; ok {
			everyData["ChildList"] = childPermit[v.Id]
		}
		result = append(result, everyData)
	}

	return &result
}

func (this *Permit) orgAdminPermit(v Permit, params map[string]string) *PermitAdmin {
	domain := "default" //default默认为当前域名,此处为域名的MAP映射
	m := this.getDefaultModuleControllerAction(v)
	permitLeft := PermitAdmin{*m, params, domain}
	return &permitLeft
}

func (this *Permit) getDefaultModuleControllerAction(v Permit) *Permit {

	//判断是否为header属性
	if "" == v.Controller && "" == v.Action && "" != v.Module {
		return this.getModuleDefaultPermit(v)
	}
	return &v
}
func (this *Permit) getModuleDefaultPermit(permit Permit) *Permit {

	//默认访问地址
	permit.Controller = "main"
	permit.Action = "goto"
	return &permit
}
