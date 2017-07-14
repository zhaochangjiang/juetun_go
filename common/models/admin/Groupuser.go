package admin

import (
	"github.com/astaxie/beego/orm"
)

type Groupuser struct {
	AdminUserid string `orm:"column(admin_userid);pk" json:"admin_userid"`
	GroupId     string `orm:"column(group_id)"`
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
func (this *Groupuser) getQuerySeter() orm.QuerySeter {
	return this.getOrm().QueryTable(this)
}

/**
*根据用户ID 获得用户的权限组列表
 */
func (this *Groupuser) GetGoupList(uid string) (*[]Groupuser, error) {

	var groupuser []Groupuser
	var querySeter orm.QuerySeter

	//查询上级权限为leftTopId的权限列表
	querySeter = this.getQuerySeter().Filter("admin_userid__exact", uid)
	_, err := querySeter.All(&groupuser)
	return &groupuser, err
}
