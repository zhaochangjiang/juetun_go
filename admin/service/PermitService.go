package service

import (
	modelsadmin "juetun/common/models/admin"
)

type PermitService struct {
	modelsadmin.Permit
}

//根据上级权限，查询所有下级权限
//@param uponid []interface{}
//@return
//	  *[]Permit,//权限列表
//     int64,	//查询结果数量
//     error
func (this *PermitService) FetchPermitListByUponId(uponid []interface{}) (*[]modelsadmin.Permit, int64, error) {
	var permitList []modelsadmin.Permit
	var querySeter orm.QuerySeter

	querySeter = this.getOrm().QueryTable(this).Filter("uppermit_id__in", uponid...).OrderBy("-id")
	num, err := querySeter.All(&permitList)

	return &permitList, num, err
}

//查询单个权限
func (this *PermitService) FetchPermit(argument map[string]interface{}) (*[]modelsadmin.Permit, error) {

	var permitList []modelsadmin.Permit

	var querySeter orm.QuerySeter

	o := this.getOrm()

	ob := o.QueryTable(this)

	for k, v := range argument {
		querySeter = ob.Filter(k, v)
	}
	_, err := querySeter.All(&permitList)

	return &permitList, err
}

//删除权限
func (this *PermitService) DeletePermit(permitIds []int) (bool, error) {

	//删除相关的数据
	groupPermit := new(GroupPermit)
	_, err1 := groupPermit.DeleteByPermitIds(permitIds)
	if nil != err1 {
		return false, err1
	}

	//删除表头信息
	_, err := this.getOrm().QueryTable(this.TableName()).Filter("id__in", permitIds).Delete()
	if nil != err {
		return false, err
	}
	return true, err
}

//左边展示权限的结构体
type LeftPermit struct {
}

//组织左边权限数据的查询
//@return
func (this *PermitService) OrgPermitLeftData(permitUpon *[]*modelsadmin.Permit) {

	var leftPermit = make(map[int][]*modelsAdmin.Permit)
	for k, v := range *permitUpon {
		if nil == leftPermit[v.UppermitId] {
			leftPermit[v.UppermitId] = make([]modelsAdmin.Permit, 0)

		}
		append(leftPermit[v.UppermitId], v)
	}
}
