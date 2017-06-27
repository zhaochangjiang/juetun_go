package admin

import (
	"github.com/astaxie/beego/orm"
)

type Permit struct {
	Id         string `orm:"column(id);pk" json:"id"`
	Name       string `orm:varchar(50)`
	Module     string `orm:varchar(30)`
	Controller string `orm:varchar(30)`
	Action     string `orm:varchar(30)`
	UppermitId int    `orm:int(10)`
	Obyid      string
	Csscode    string `orm:varchar(500)`
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Permit))
}
func (this *Permit) TableName() string {
	return "permit"
}
func (this *Permit) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}

//根据上级权限，查询所有下级权限
func (this *Permit) FetchPermitListByUponId(uponid []interface{}) (*[]Permit, int64, error) {
	var permitList []Permit
	var querySeter orm.QuerySeter

	querySeter = this.getOrm().QueryTable(this).Filter("uppermit_id__in", uponid...).OrderBy("-id")
	num, err := querySeter.All(&permitList)

	return &permitList, num, err
}

//查询单个权限
func (this *Permit) FetchPermit(argument map[string]interface{}) ([]*Permit, error) {

	var permitList []*Permit
	o := this.getOrm()
	querySeter := o.QueryTable(this)
	for k, v := range argument {
		querySeter = querySeter.Filter(k, v)
	}
	_, err := querySeter.All(&permitList)

	return permitList, err
}

//删除权限
func (this *Permit) DeletePermit(permitIds []int) (bool, error) {

	//删除相关的数据
	groupPermit := new(GroupPermit)
	_, err1 := groupPermit.DeleteByPermitIds(permitIds)
	if nil != err1 {
		return false, err1
	}

	//删除表头信息
	_, err := this.getOrm().QueryTable(this.TableName()).Filter("id__in", permitIds).Delete()
	if nil != err {
		return false, err
	}
	return true, err
}

func (this *Permit) getLeftPermit(leftTopId string) {
	var permitList []Permit
	var querySeter orm.QuerySeter

	querySeter = this.getOrm().QueryTable(this).qs.Filter("uppermit_id__gt", "").OrderBy("obyid")
	num, err := querySeter.All(&permitList)

	return &permitList, num, err
}
