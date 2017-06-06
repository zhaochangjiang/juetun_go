package admin

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type Permit struct {
	Id         int    `orm:"column(id);pk;auto" json:"id"`
	Name       string `orm:varchar(50)`
	Module     string `orm:varchar(30)`
	Controller string `orm:varchar(30)`
	Action     string `orm:varchar(30)`
	UppermitId int    `orm:int(10)`
	Obyid      int
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
func (this *Permit) FetchPermitListByUponId(uponid []interface{}) (*[]Permit, int64, string) {
	var permitList []Permit
	var message string
	var querySeter orm.QuerySeter

	querySeter = this.getOrm().QueryTable(this).Filter("uppermit_id__in", uponid...).OrderBy("-id")

	num, err := querySeter.All(&permitList)
	//	err := o.QueryTable(this).Filter("username", userName).Filter("flag_del", "no").All(&permitList)
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		message = "data exception,please cotact the administrator!"
		return nil, num, message
	}
	return &permitList, num, message
}

//查询单个权限
func (this *Permit) FetchPermit(argument map[string]interface{}) ([]*Permit, string) {

	var permitList []*Permit
	var message string
	var querySeter orm.QuerySeter

	o := this.getOrm()

	ob := o.QueryTable(this)
	fmt.Println("fdafdasd:")

	for k, v := range argument {
		fmt.Println(v)
		querySeter = ob.Filter(k, v)
	}
	//num, err := o.QueryTable("user").Filter("name", "slene").All(&users)

	_, err := querySeter.All(&permitList)
	//	err := o.QueryTable(this).Filter("username", userName).Filter("flag_del", "no").All(&permitList)
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		message = "data exception,please cotact the administrator!"
		return nil, message
	}
	return permitList, message
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
