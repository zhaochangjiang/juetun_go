package admin

import (
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
func (this *Permit) GetOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}
func (this *Permit) GetQuerySelect() orm.QuerySeter {
	return this.GetOrm().QueryTable(this.TableName())
}

//根据上级权限，查询所有下级权限
//@param uponid []interface{}
//@return
//	  *[]Permit,//权限列表
//     int64,	//查询结果数量
//     error
func (this *Permit) FetchPermitListByUponId(uponid []interface{}) (*[]Permit, int64, error) {
	var permitList []Permit
	num, err := this.GetQuerySelect().Filter("uppermit_id__in", uponid...).OrderBy("-id").All(&permitList)
	return &permitList, num, err
}

//查询单个权限
func (this *Permit) FetchPermit(argument map[string]interface{}) (*[]Permit, error) {

	var permitList []Permit

	var querySeter orm.QuerySeter

	ob := this.GetOrm().QueryTable(this)

	for k, v := range argument {
		querySeter = ob.Filter(k, v)
	}
	_, err := querySeter.All(&permitList)

	return &permitList, err
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
	_, err := this.GetOrm().QueryTable(this.TableName()).Filter("id__in", permitIds).Delete()
	if nil != err {
		return false, err
	}
	return true, err
}

//左边展示权限的结构体
type LeftPermit struct {
}

//组织左边权限数据的查询
//@return
func (this *Permit) OrgPermitLeftData(permitUpon *[]*Permit) *map[int][]*Permit {

	var leftPermit = make(map[int][]*Permit)
	for _, v := range *permitUpon {
		if nil == leftPermit[v.UppermitId] {
			leftPermit[v.UppermitId] = make([]*Permit, 0)

		}
		leftPermit[v.UppermitId] = append(leftPermit[v.UppermitId], v)
	}
	return &leftPermit
}
