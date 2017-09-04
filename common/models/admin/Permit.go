package admin

import (
	"juetun/common/general"
	"log"

	"strings"

	"github.com/astaxie/beego/orm"
)

type Permit struct {
	CommonModel
	Id         string `orm:"column(id);pk" json:"id"`
	Name       string `orm:varchar(50);orm:"column(name)"`
	Mod        string `orm:varchar(30);orm:"column(mod)"`
	Controller string `orm:varchar(30);orm:"column(controller)"`
	Action     string `orm:varchar(30);orm:"column(action)"`
	UppermitId string `orm:int(10);orm:"column(uppermit_id)"`
	DomainMap  string `orm:"column(domain_map)"`
	Obyid      string `orm:"column(obyid)"`
	Csscode    string `orm:varchar(500);orm:"column(csscode)"`
}

func init() {
	permit := new(Permit)
	orm.RegisterModelWithPrefix(permit.GetTablePrefix(), permit)
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
		params["mod"] = v.Mod
		result = append(result, *this.OrgAdminPermit(v, params))
	}
	return &result, num, err
}

//查询单个权限
func (this *Permit) FetchPermit(argument map[string]string) (*[]Permit, error) {

	var permitList []Permit
	querySeter := this.getQuerySeter()
	for k, v := range argument {
		querySeter = querySeter.Filter(k, v)
	}
	_, err := querySeter.All(&permitList)

	return &permitList, err
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
	Params    map[string]string
	Domain    string
	UrlString string
}

/**
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *Permit) FetchDefaultPermitByModuleString(moduleString string, controllerContext general.ControllerContext) {

	if "" == moduleString {
		panic("moduleString is null string")
		return
	}

	var leftPermit *[](map[string]interface{})

	var leftTopId string

	var args map[string]interface{}
	args["SuperAdmin"] = controllerContext.IsSuperAdmin
	if controllerContext.IsSuperAdmin == true {
		var fetchParams = make(map[string]string)
		fetchParams["mod"] = moduleString
		fetchParams["uppermit_id"] = ""
		permit, err := this.FetchPermit(fetchParams)
		if err != nil {
			panic(err)
		}
		if len(*permit) != 0 {
			leftTopId = (*permit)[0].Id
		}
		args["LeftTopId"] = leftTopId

	} else {
		var fetchParams = make(map[string]interface{})
		var permitParams = make(map[string]interface{})
		permitParams["mod"] = moduleString
		permitParams["uppermit_id"] = ""
		fetchParams["PermitTable"] = permitParams

		var groupPermitParams = make(map[string]interface{})
		groupPermitParams["GroupIds"] = controllerContext.GroupIds
		fetchParams["GrouppermitTable"] = groupPermitParams

		permit := this.FetchPermitByCondition(fetchParams)

		if len(*permit) != 0 {
			leftTopId = (*permit)[0].Id
		}
		args["LeftTopId"] = leftTopId
		args["GroupIds"] = controllerContext.GroupIds
	}
	leftPermit = this.GetLeftPermit(args)

	for _, v := range *leftPermit {
		p := (v["Permit"]).(*PermitAdmin)
		if p.Controller != "" && p.Action != "" {
			this.Controller = p.Controller
			this.Action = p.Action
			this.DomainMap = p.DomainMap
			break
		}
		var flag = 0
		childList := v["ChildList"].([]PermitAdmin)
		for _, permitAdmin := range childList {
			//permitAdmin := v1.(PermitAdmin)
			if permitAdmin.Controller != "" && permitAdmin.Action != "" {
				this.Controller = permitAdmin.Controller
				this.Action = permitAdmin.Action
				this.DomainMap = permitAdmin.DomainMap
				flag = 1
				break
			}
		}
		if flag == 1 {
			break
		}
	}

}

/**
* 获得左边的权限列表
* @author karl.zhao<zhaocj2009@126.com>
* @Date 2017/08/01
*
 */
func (this *Permit) GetLeftPermit(args map[string]interface{}) *[](map[string]interface{}) {
	var permitList, childPermitList []Permit
	var result = make([]map[string]interface{}, 0)

	var leftTopId string = ""
	var groupIds []string = make([]string, 0)
	var superAdmin = false
	if _, ok := args["LeftTopId"]; ok {
		leftTopId = args["LeftTopId"].(string)
	}

	if _, ok := args["SuperAdmin"]; ok {
		superAdmin = args["SuperAdmin"].(bool)
	}

	if leftTopId == "" {
		return &result
	}
	if superAdmin == true {
		querySeter := this.getQuerySeter()

		//查询上级权限为leftTopId的权限列表
		querySeter = querySeter.Filter("uppermit_id__exact", leftTopId).OrderBy("obyid")
		querySeter.All(&permitList)
		leftPermitIdList := make([]string, 0)
		leftPermitIdList = append(leftPermitIdList, leftTopId)
		for _, v := range permitList {
			leftPermitIdList = append(leftPermitIdList, v.Id)
		}

		this.getQuerySeter().Filter("uppermit_id__in", leftPermitIdList).OrderBy("obyid").All(&childPermitList)

	} else {
		if _, ok := args["GroupIds"]; ok {
			groupIds = args["GroupIds"].([]string)
		}
		var result = make([](map[string]interface{}), 0)
		var leftPermitIdList []string

		if leftTopId == "" {
			return &result
		}
		permitList := this.FetchPermitByGroupIdAndUppermit(groupIds, []string{leftTopId})
		for _, v := range *permitList {
			leftPermitIdList = append(leftPermitIdList, v.Id)
		}
		tmp := this.FetchPermitByGroupIdAndUppermit(groupIds, leftPermitIdList)
		childPermitList = *tmp
	}

	return this.orgPermitSepData(&childPermitList, leftTopId)
}

/**
*
*
*
 */
func (this *Permit) orgPermitSepData(childPermitList *[]Permit, leftTopId string) *([](map[string]interface{})) {
	var childPermit = make(map[string][]interface{})
	for _, v := range *childPermitList {
		params := make(map[string]string)
		tmp := this.OrgAdminPermit(v, params)
		obj := make(map[string]interface{})
		obj["Permit"] = *tmp
		obj["Active"] = false
		obj["ChildList"] = make([]interface{}, 0)
		childPermit[v.UppermitId] = append(childPermit[v.UppermitId], obj)
	}
	var result = make([](map[string]interface{}), 0)

	if _, ok := childPermit[leftTopId]; ok {
		result = childPermit[leftTopId]

		if len(result) > 0 {

			for k, v := range result {
				permit := v["Permit"].(Permit)
				if _, ok := childPermit[permit.Id]; ok {
					result[k]["ChildList"] = childPermit[permit.Id]

					if len(result[k]["ChildList"]) > 0 {

					}
				}
			}
		}

	}
	result := this.iteration(childPermit)
	return result
}
func (this *Permit) iteration(everyData Permit, childPermit map[string][]interface{}) *[]map[string]interface{} {
	var result = make([](map[string]interface{}), 0)

	if _, ok := childPermit[permit.Id]; ok {
		if len(childPermit[permit.Id]) > 0 {
			result[k]["ChildList"] = childPermit[permit.Id]

			for k, v := range result[k]["ChildList"] {

			}
		}

		//this.iteration(,childPermit)

	}

	return &result
}

//func (this *Permit) iteration(childPermit map[string][]interface{}) *[]map[string]interface{} {
//	var result = make([](map[string]interface{}), 0)

//	for _, v := range *permitList {

//		everyData := make(map[string]interface{})

//		everyData["Permit"] = &v

//		everyData["Active"] = false //默认设置不为选中的状态
//		everyData["ChildList"] = make([]interface{}, 0)

//		//判断内容是否存在，相当于PHP中的isset函数
//		if _, ok := childPermit[v.Id]; ok {
//			tmp := childPermit[v.Id]

//			if len(tmp) > 0 {
//				tmp1 := make([]PermitAdmin, 0)
//				for _, v := range tmp {
//					tmp1 = append(tmp1, v["Permit"].(PermitAdmin))
//				}
//				//	tmp := this.iteration(&tmp1, childPermit)
//			}
//			everyData["ChildList"] = tmp

//		}

//		result = append(result, everyData)
//	}
//	return &result
//}

/**
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/07/20
 */
func (this *Permit) FetchPermitByGroupIdAndUppermit(groupIds []string, uppermitIds []string) *[]Permit {

	var sliceParams []string
	var where string
	var nowTableName string
	var leftTableName string
	var permitList []Permit
	var groupPermit GroupPermit

	if len(groupIds) == 0 || len(uppermitIds) == 0 {
		return &permitList
	}

	commonModel := new(CommonModel)
	//获得表前缀
	tablePrefix := commonModel.GetTablePrefix()

	// 构建查询对象
	nowTableName = tablePrefix + this.TableName()
	leftTableName = tablePrefix + groupPermit.TableName()

	//查询上级权限为leftTopId的权限列表
	where += leftTableName + ".group_id in (\"" + strings.Join(groupIds, "\",\"") + "\")"
	where += " AND " + nowTableName + ".uppermit_id in (\"" + strings.Join(uppermitIds, "\",\"") + "\")"

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select(nowTableName + ".*").From(nowTableName).
		LeftJoin(leftTableName).On(leftTableName + ".permit_id=" + nowTableName + ".id").Where(where).OrderBy(nowTableName + ".obyid").Asc().String()
	_, err := this.getOrm().Raw(sql, sliceParams).QueryRows(&permitList)
	if nil != err {
		panic(err)
	}
	return &permitList
}

/**
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/07/20
* @param condition map[string]interface{}
* 							PermitTable map[string]interface{}
*									mod	 string
*									uppermit_id string
*									uppermit_ids []string
*							GrouppermitTable map[string]interface{}
*									GroupIds []string
 */
func (this *Permit) FetchPermitByCondition(condition map[string]interface{}) *[]Permit {
	//   groupIds []string, uppermitIds []string
	var sliceParams []string
	var where string = "1 "
	var nowTableName string
	var leftTableName string
	var permitList []Permit
	var groupPermit GroupPermit
	if condition == nil {
		return &permitList
	}
	commonModel := new(CommonModel)
	//获得表前缀
	tablePrefix := commonModel.GetTablePrefix()

	if _, ok := condition["PermitTable"]; ok {
		// 构建查询对象
		nowTableName = tablePrefix + this.TableName()

		tmp := condition["PermitTable"].(map[string]interface{})
		if _, ok := tmp["mod"]; ok {
			mod := tmp["mod"].(string)
			where += " AND " + nowTableName + ".`mod`=\"" + mod + "\""

		}
		if _, ok := tmp["uppermit_id"]; ok {
			uppermitId := tmp["uppermit_id"].(string)
			where += " AND " + nowTableName + ".`uppermit_id` =\"" + uppermitId + "\""
		}
		if _, ok := tmp["uppermit_ids"]; ok {
			uppermitIds := tmp["uppermit_ids"].([]string)
			where += " AND " + nowTableName + ".`uppermit_id` in (\"" + strings.Join(uppermitIds, "\",\"") + "\")"
		}
	}
	if _, ok := condition["GrouppermitTable"]; ok {
		leftTableName = tablePrefix + groupPermit.TableName()
		tmp := condition["GrouppermitTable"].(map[string]interface{})
		if _, ok := tmp["GroupIds"]; ok {
			groupIds := tmp["GroupIds"].([]string)
			//查询上级权限为leftTopId的权限列表
			where += " AND " + leftTableName + ".`group_id` in (\"" + strings.Join(groupIds, "\",\"") + "\")"
		}

	}

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select(nowTableName + ".*").From(nowTableName).
		LeftJoin(leftTableName).On(leftTableName + ".`permit_id`=" + nowTableName + ".id").Where(where).OrderBy(nowTableName + ".obyid").Asc().String()
	_, err := this.getOrm().Raw(sql, sliceParams).QueryRows(&permitList)
	if nil != err {
		panic(err)
	}
	return &permitList
}

/**
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/07/20
 */
func (this *Permit) FetchPermitByGroupId(groupIds []string, condition map[string]string) (*[]Permit, error) {

	var sliceParams []string
	var where string
	var nowTableName string
	var leftTableName string
	var permitList []Permit
	var groupPermit GroupPermit

	if len(groupIds) == 0 {
		return &permitList, nil
	}

	commonModel := new(CommonModel)
	//获得表前缀
	tablePrefix := commonModel.GetTablePrefix()

	// 构建查询对象
	nowTableName = tablePrefix + this.TableName()
	leftTableName = tablePrefix + groupPermit.TableName()

	//查询上级权限为leftTopId的权限列表
	where += leftTableName + ".group_id in (\"" + strings.Join(groupIds, "\",\"") + "\")"
	for k, v := range condition {
		where += " AND " + nowTableName + "." + k + " = \"" + v + "\""
	}

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select(nowTableName + ".*").From(nowTableName).
		LeftJoin(leftTableName).On(leftTableName + ".permit_id=" + nowTableName + ".id").Where(where).OrderBy(nowTableName + ".obyid").Asc().String()
	_, err := this.getOrm().Raw(sql, sliceParams).QueryRows(&permitList)
	if nil != err {
		panic(err)
	}
	return &permitList, nil
}

/**
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/07/20
 */
func (this *Permit) OrgAdminPermit(v Permit, params map[string]string) *PermitAdmin {

	domain := "default" //default默认为当前域名,此处为域名的MAP映射

	m := this.getDefaultModuleControllerAction(v)

	//获得链接字符串
	urlString := general.CreateUrl(v.Controller, v.Action, params, v.DomainMap)

	permitLeft := PermitAdmin{*m, params, domain, urlString}
	return &permitLeft
}

/**
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/07/20
 */
func (this *Permit) getDefaultModuleControllerAction(v Permit) *Permit {

	//判断是否为header属性
	if "" == v.Controller && "" == v.Action && "" != v.Mod {
		return this.getModuleDefaultPermit(v)
	}
	return &v
}

/**
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/07/20
 */
func (this *Permit) getModuleDefaultPermit(permit Permit) *Permit {

	//默认访问地址
	permit.Controller = "main"
	permit.Action = "goto"
	return &permit
}
