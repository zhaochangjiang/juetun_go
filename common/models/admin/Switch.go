package admin

import (
	"github.com/astaxie/beego/orm"
)

type Switch struct {
	CommonModel
	Id             string `orm:"column(id);varchar(32);pk" json:"id"`
	Name           string `orm:"varchar(20)";orm:"column(name)"`
	Status         string `orm:"column(status)"`
	HeigherLevelId string `orm:varchar(32);orm:"column(heigher_level_id)"`
	IsShow         string `orm:varchar(15);orm:"column(is_show)"`
	Desc           string `orm:varchar(300);orm:"column(desc)"`
}

func init() {
	switchModel := new(Switch)
	orm.RegisterModelWithPrefix(switchModel.GetTablePrefix(), switchModel)
}
func (u *Switch) TableName() string {
	return "switch"
}
func (this *Switch) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}
