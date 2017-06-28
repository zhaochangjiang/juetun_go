package admin

import (
	"github.com/astaxie/beego/orm"
)

type Serverlist struct {
	Id          string `orm:"column(id);pk" json:"id"`
	IpAddr      string `orm:varchar(30);orm:"column(ip_addr);`
	MachineRoom string `orm:varchar(30);orm:"column(machine_room);`
	UniqueKey   string `orm:varchar(255);orm:"column(unique_key);`
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Serverlist))
}
func (u *Serverlist) TableName() string {
	return "serverlist"
}
func (this *Serverlist) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}
