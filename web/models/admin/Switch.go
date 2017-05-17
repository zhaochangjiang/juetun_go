package admin

import (
	"github.com/astaxie/beego/orm"
)

type Switch struct {
	id             string `orm:"column(id);varchar(32);pk;auto" json:"id"`
	Name           string `orm:"varchar(20)"`
	Status         string
	HeigherLevelId string `orm:varchar(32)`
	IsShow         string
	Desc           string `orm:varchar(300)`
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Switch))
}
func (u *Switch) TableName() string {
	return "picture_category"
}
func (this *Switch) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}
