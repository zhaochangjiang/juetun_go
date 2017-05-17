package admin

import (
	"github.com/astaxie/beego/orm"
)

type Group struct {
	Id         int `orm:"column(id);pk;auto" json:"id"`
	Name       string
	SuperAdmin `orm:varchar(30)`
	UpGroupid  int
}

func (u *Group) TableName() string {
	return "group"
}
func init() {
	orm.RegisterModelWithPrefix("admin_", new(Group))
}
func (this *Group) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}
