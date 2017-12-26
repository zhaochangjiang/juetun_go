package admin

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"juetun/common/general"
	"juetun/common/utils"
	"strings"
	"time"
)

type Permit struct {
	CommonModel
	Id         string `orm:"column(id);pk" json:"id";form:"-"`
	Name       string `orm:varchar(50);orm:"column(name)";form:"name"`
	Mod        string `orm:varchar(30);orm:"column(mod)";form:"mod"`
	Controller string `orm:varchar(30);orm:"column(controller)";form:"controller"`
	Action     string `orm:varchar(30);orm:"column(action)";form:"action"`
	UppermitId string `orm:int(10);orm:"column(uppermit_id)";form:"uppermit_id"`
	DomainMap  string `orm:"column(domain_map)";form:"domain_map"`
	Obyid      string `orm:"column(obyid)";form:"obyid"`
	Csscode    string `orm:varchar(500);orm:"column(csscode)";form:"csscode"`
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
	Params    map[string]string //其他参数
	Domain    string            //域名mapkey
	UrlString string            //Url字符串
	Active    bool              //当前是否为选中项
}

func (p *Permit) GetCacheKeys(key string) string {
	if key == "" {
		panic("cache key must be not null!")
	}
	return "per_" + key
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

//从redis缓存中读取相应的信息
func (this *Permit) GetContentFromCache(key string) interface{} {

	var cache = *this.GetCache()
	return cache.Get(key)
}

//将数据放入redis缓存中。
func (this *Permit) SetContentFromCache(key string, content interface{}, lifeTime time.Duration) {
	(*this.GetCache()).Put(key, content, lifeTime) //time.Second*3600
}

func (this *Permit) GetRedisConfig() (string, string, string) {

	var cacheType = "redis"
	var redisName = "redisCache"
	var dbNumber = "0"
	return cacheType, redisName, dbNumber
}

//获得redis缓存中的内容
func (this *Permit) GetCache() *cache.Cache {
	if nil == bm {

		cacheType, redisName, dbNumber := this.GetRedisConfig()
		var err error
		var redisconfig = new(utils.RedisConfig)
		var con = redisconfig.GetJsonCode(redisName, dbNumber)
		bm, err = cache.NewCache(cacheType, con)
		if nil != err {
			panic(err.Error())
		}
	}
	return &bm
}

//根据上级权限，查询所有下级权限
func (this *Permit) FetchPermitListByUponId(uponid *[]string) (*[]PermitAdmin, int64, error) {

	var err error
	var permitList = make([]Permit, 0)
	var uponidParams = make([]interface{}, 0)
	var cachePrefix = this.GetCacheKeys("p_")
	var result = make([]PermitAdmin, 0)
	var num int64 = 0

	for _, v := range *uponid {
		var cachekey = cachePrefix + v
		var ct = this.GetContentFromCache(cachekey)

		if ct == nil {
			uponidParams = append(uponidParams, v)
			continue
		}

		var c = ct.([]byte)
		j2 := new(Permit)
		err = json.Unmarshal(c, j2)
		if err != nil { //如果错误信息存在，则说明数据异常，直接从数据库中取数据
			uponidParams = append(uponidParams, v)
			continue
		}
		//如果缓存中不存在，则从数据库中查询
		permitList = append(permitList, *j2)

	}
	if len(uponidParams) == 0 {
		return &result, num, err
	}
	var opermitList []Permit
	num, err = this.getQuerySeter().Filter("uppermit_id__in", uponidParams...).All(&opermitList) //OrderBy("-id")

	if err != nil {
		return &result, num, err
	}
	for _, v := range opermitList {
		permitList = append(permitList, v)

		//将数据转换为json字符串。
		js1, err := json.Marshal(v)
		if err == nil { //只有没有错误的时候才保存到redis中
			this.SetContentFromCache(cachePrefix+v.Id, string(js1[:]), time.Second*3600)
		}
	}

	for _, v := range permitList {
		params := make(map[string]string)
		params["mod"] = v.Mod
		result = append(result, *this.OrgAdminPermit(v, params))
	}
	return &result, num, err
}

//查询单个权限
func (this *Permit) FetchPermit(argument map[string]interface{}) (*[]Permit, error) {

	var permitList = make([]Permit, 0)
	querySeter := this.getQuerySeter()
	var flag = false
	for k, v := range argument {
		switch v.(type) {
		case uint8, uint16, uint32, int64, string, int8, int16, int, int32, float32, float64:
			querySeter = querySeter.Filter(k, v)
			flag = true
			break
		case []string:
			tmp := v.([]string)
			if len(tmp) > 0 {
				con := make([]interface{}, 0)
				for _, v1 := range tmp {
					con = append(con, v1)
				}
				flag = true
				querySeter = querySeter.Filter(k, con...)

			}
		}
	}
	if flag == false {
		return &permitList, nil
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

func (this *PermitAdmin) GetOperate() (string, error) {
	fmt.Println("-----------")
	return "213123", nil
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

	var args = make(map[string]interface{})
	args["SuperAdmin"] = controllerContext.IsSuperAdmin
	if controllerContext.IsSuperAdmin == true {
		var fetchParams = make(map[string]interface{})
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
		p := (v["Permit"]).(PermitAdmin)
		if p.Controller != "" && p.Action != "" {
			this.Controller = p.Controller
			this.Action = p.Action
			this.DomainMap = p.DomainMap
			break
		}
		var flag = 0
		childList := v["ChildList"].(*[]map[string]interface{})
		for _, v := range *childList {
			permitAdmin := v["Permit"].(PermitAdmin)
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

	result := this.iteration(leftTopId, childPermit)

	return result
}

/**
*
* 递归获得子选项
* @author karl.zhao<zhaochangjiang@huoyunren.com>
* @date 20170905
*
*
 */
func (this *Permit) iteration(topId string, childPermit map[string][]interface{}) *[]map[string]interface{} {
	var childList = make([](map[string]interface{}), 0)
	if _, ok := childPermit[topId]; ok {

		if len(childPermit[topId]) > 0 {
			for _, v := range childPermit[topId] {
				vk := v.(map[string]interface{})
				v1 := vk["Permit"].(PermitAdmin)
				vk["ChildList"] = this.iteration(v1.Id, childPermit)

				childList = append(childList, vk)
			}
		}
	}
	return &childList
}

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
		LeftJoin(leftTableName).On(leftTableName + ".permit_id=" + nowTableName + ".id").Where(where).OrderBy(nowTableName + ".obyid").Desc().String()
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
		LeftJoin(leftTableName).On(leftTableName + ".`permit_id`=" + nowTableName + ".id").Where(where).OrderBy(nowTableName + ".obyid").Desc().String()
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
func (this *Permit) FetchPermitByGroupId(groupIds []string, condition map[string]interface{}) (*[]Permit, error) {

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

		where += " AND " + nowTableName + "." + k + " = \"" + v.(string) + "\""
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

	permitLeft := PermitAdmin{*m, params, domain, "", false}

	//获得链接字符串
	permitLeft.UrlString = general.CreateUrl(m.Controller, m.Action, params, m.DomainMap)

	return &permitLeft
}

/**
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/07/20
 */
func (this *Permit) getDefaultModuleControllerAction(v Permit) *Permit {

	//判断是否为header属性
	if "" == strings.Trim(v.Controller, "") && "" == strings.Trim(v.Action, "") && "" != strings.Trim(v.Mod, "") {
		return this.GetModuleDefaultPermit(v)
	}
	return &v
}

/**
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/07/20
 */
func (this *Permit) GetModuleDefaultPermit(permit Permit) *Permit {

	//默认访问地址
	permit.Controller = "loc"
	permit.Action = "goto"
	return &permit
}

func (this *Permit) GetList(args map[string]interface{}) (*[][]PermitAdmin, *[]PermitAdmin) {

	pid, isSuperAdmin, groupIdsPointer := this.getDefaultArgs(&args)

	pidsPointer, _ := this.getAllUponByPid(pid, isSuperAdmin, groupIdsPointer)

	pList := this.GetAllChildByPids(pidsPointer, isSuperAdmin, groupIdsPointer)

	var res = make([][]PermitAdmin, 0)
	var resChild = make([]PermitAdmin, 0)
	for _, v := range *pidsPointer {
		tmp := make([]PermitAdmin, 0)
		if v == "" {
			v = "-1"
		}

		if _, ok := (*pList)[v]; ok {
			res = append(res, (*pList)[v])
		} else {
			res = append(res, tmp)
		}

	}

	for k, v := range res {
		for k1, v1 := range v {
			for _, v2 := range *pidsPointer {
				if v2 == v1.Id {
					res[k][k1].Active = true
				}
			}
		}
	}
	lenRes := len(res)
	if lenRes > 0 {
		f := (lenRes - 1)
		resChild = res[f]
	}

	return &res, &resChild
}

func (this *Permit) GetAllChildListByPids(pidsPointer *[]string, isSuperAdmin bool, groupIdsPointer *[]string) *[]PermitAdmin {
	var childList []PermitAdmin

	//如果是超级管理员
	if isSuperAdmin == true {
		var params = make(map[string]interface{})
		params["uppermit_id__in"] = *pidsPointer
		tmp, _ := this.FetchPermit(params)
		for _, v := range *tmp {
			params := make(map[string]string)
			permitDataTmp := this.OrgAdminPermit(v, params)
			childList = append(childList, *permitDataTmp)
		}
	} else { //如果不是超级管理员
		groupPermit := new(GroupPermit)

		tmp, _ := groupPermit.GetGroupPermitList(*groupIdsPointer, *pidsPointer)
		childList = *tmp
	}
	return &childList
}

/**
*获得子权限
*
 */
func (this *Permit) GetAllChildByPids(pidsPointer *[]string, isSuperAdmin bool, groupIdsPointer *[]string) *map[string][]PermitAdmin {

	var res = make(map[string][]PermitAdmin)
	childList := this.GetAllChildListByPids(pidsPointer, isSuperAdmin, groupIdsPointer)
	for _, v := range *childList {
		upid := v.UppermitId
		if upid == "" {
			upid = "-1"
		}
		if _, ok := res[upid]; ok {
			res[upid] = append(res[upid], v)
		} else {
			res[upid] = make([]PermitAdmin, 0)
			res[upid] = append(res[upid], v)
		}
	}
	return &res
}

/**
* 获得所有上级权限
*
*
 */
func (this *Permit) getAllUponByPid(pid string, isSuperAdmin bool, groupIdsPointer *[]string) (*[]string, *[]PermitAdmin) {
	var pList = make([]PermitAdmin, 0)

	permitData := this.GetPermitByPid(pid, isSuperAdmin, groupIdsPointer)

	if nil != permitData {
		//将当前信息放入进去
		params := make(map[string]string)
		permitDataTmp := this.OrgAdminPermit(*permitData, params)
		pList = append(pList, *permitDataTmp)
		i := 0
		for {
			i++
			//判断如果循环超过5次还没中断，则强制中断，防止程序异常
			if i > 6 {
				break
			}

			permitData = this.GetPermitByPid(permitData.UppermitId, isSuperAdmin, groupIdsPointer)
			//如果没有查询到数据，那么跳出循环
			if permitData == nil {
				break
			}

			params := make(map[string]string)
			permitDataTmp := this.OrgAdminPermit(*permitData, params)

			//往队列的队首添加数据
			slice := []PermitAdmin{*permitDataTmp}
			pList = append(slice, pList...)

		}
	}

	var pids = make([]string, 0)
	var permitDataTmpCommon = new(PermitAdmin)
	permitDataTmpCommon.Id = ""
	slice := []PermitAdmin{*permitDataTmpCommon}
	pList = append(slice, pList...)

	for _, v := range pList {
		pids = append(pids, v.Id)
	}

	return &pids, &pList
}

func (this *Permit) GetPermitByPid(pid string, isSuperAdmin bool, groupIdsPointer *[]string) *Permit {
	var permitModelList *[]Permit

	var permitData Permit
	var params = make(map[string]interface{})
	params["id"] = pid
	if isSuperAdmin == true {

		permitModelList, _ = this.FetchPermit(params)
	} else {

		permitModelList, _ = this.FetchPermitByGroupId(*groupIdsPointer, params)
	}
	if len(*permitModelList) <= 0 {
		return nil
	}
	permitData = (*permitModelList)[0]
	return &permitData
}
func (this *Permit) getDefaultArgs(args *map[string]interface{}) (string, bool, *[]string) {
	var pid = ""
	var isSuperAdmin = true
	var groupIds = make([]string, 0)
	if utils.Isset("IsSuperAdmin", *args) == false {
		panic("args is not exists key names 'IsSuperAdmin'")
	} else {
		isSuperAdmin = (*args)["IsSuperAdmin"].(bool)
	}

	if utils.Isset("Pid", *args) == false {
		panic("args is not exists key names 'Pid'")
	} else {
		pid = (*args)["Pid"].(string)
	}
	if utils.Isset("GroupIds", *args) == false {
		panic("args is not exists key names 'GroupIds'")
	} else {
		groupIds = (*args)["GroupIds"].([]string)
	}
	return pid, isSuperAdmin, &groupIds
}

//更新数据
func (this *Permit) UpdateDataById(update *Permit) bool {
	var o = this.getOrm()

	var e = o.Read(&Permit{Id: update.Id})
	if nil != e {
		panic(e.Error())
	}
	update.Id = this.Id
	if _, err := o.Update(update); err == nil {
		return true
	} else {
		return false
		//panic(err.Error())
	}

	return true
}
