package admin

import (
	"github.com/astaxie/beego/orm"
)

type Groupuser struct {
	AdminUserid int `orm:"column(admin_userid);pk;auto" json:"admin_userid"`
	GroupId     int `orm:"column(group_id)"`
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Groupuser))
}

func (u *Groupuser) TableName() string {
	return "groupuser"
}

func (this *Groupuser) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}
