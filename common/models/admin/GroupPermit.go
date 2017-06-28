package admin

import (
	"github.com/astaxie/beego/orm"
)

type GroupPermit struct {
	Id       string `orm:"column(id);pk" json:"id"`
	PermitId string `orm:"column(permit_id)"`
	GroupId  string `orm:"column(group_id)"`
}

func (u *GroupPermit) TableName() string {
	return "grouppermit"
}
func init() {
	orm.RegisterModelWithPrefix("admin_", new(GroupPermit))
}
func (this *GroupPermit) getQuerySeter() orm.QuerySeter {
	return this.getOrm().QueryTable(this)
}

func (this *GroupPermit) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}

//通过权限ID 删除用户组的权限关系
func (this *GroupPermit) DeleteByGroupIds(groupIds []string) (bool, error) {

	//删除表头信息
	_, err := this.getQuerySeter().Filter("permit_id__in", groupIds).Delete()
	if nil != err {
		return false, err
	}

	return true, err
}

//通过用户群权限ID 删除权限组的权限
func (this *GroupPermit) DeleteByPermitIds(permitIds []string) (bool, error) {

	//删除表头信息
	_, err := this.getQuerySeter().Filter("group_id__in", permitIds).Delete()
	if nil != err {
		return false, err
	}

	return true, err
}
