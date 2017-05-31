package admin

import (
	"github.com/astaxie/beego/orm"
)

type Permitgroup struct {
	Id       int `orm:"column(id);pk;auto" json:"id"`
	PermitId int
	GroupId  int
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Permitgroup))
}
func (u *Permitgroup) TableName() string {
	return "permitgroup"
}
func (this *Permitgroup) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}
