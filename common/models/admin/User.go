package admin

import (
	"github.com/astaxie/beego/orm"
)

type User struct {
	Uid        int    `orm:"column(uid);pk;"`
	Name       string `orm:varchar(30)`
	SuperAdmin string `orm:"column(super_admin)"`
	Isdel      string `orm:"column(id_del)"`
}

func (u *User) TableName() string {
	return "user"
}
func init() {
	orm.RegisterModelWithPrefix("admin_", new(User))
}
func (this *User) getOrm() orm.Ormer {
	o := orm.NewOrm()
	// 默认使用 default，你可以指定为其他数据库
	o.Using("db_admin")
	return o
}
