package admin

import (
	"github.com/astaxie/beego/orm"
)

type GroupPermit struct {
	Id       int `orm:"column(id);pk;auto" json:"id"`
	PermitId int
	GroupId  int
}

func (u *GroupPermit) TableName() string {
	return "grouppermit"
}
func init() {
	orm.RegisterModelWithPrefix("admin_", new(GroupPermit))
}

func (this *GroupPermit) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}
