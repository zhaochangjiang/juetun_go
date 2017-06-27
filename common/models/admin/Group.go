package admin

import (
	"github.com/astaxie/beego/orm"
)

type Group struct {
	Id         string `orm:"column(id);pk" json:"id"`
	Name       string `orm:"column(name)"`
	SuperAdmin string `orm:varchar(30);orm:"column(super_admin)"`
	UpGroupid  int    `orm:"column(up_groupid)"`
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

//根据group_id删除数据
func (this *Group) DeleteGroup(id []int) (bool, error) {

	//删除相关的数据
	groupPermit := new(GroupPermit)
	_, err1 := groupPermit.DeleteByGroupIds(id)
	if nil != err1 {
		return false, err1
	}
	//删除表头信息
	_, err := this.getOrm().QueryTable(this.TableName()).Filter("id__in", id).Delete()
	if nil != err {
		return false, err
	}

	return true, err
}
