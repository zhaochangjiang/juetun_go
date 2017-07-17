package admin

import (
	"strings"

	"github.com/astaxie/beego/orm"
)

var dbPrefix = "admin_"

type GroupPermit struct {
	Id       string `orm:"column(id);pk" json:"id"`
	PermitId string `orm:"column(permit_id)"`
	GroupId  string `orm:"column(group_id)"`
}

func (this *GroupPermit) TableName() string {
	return "grouppermit"
}

func init() {
	orm.RegisterModelWithPrefix(dbPrefix, new(GroupPermit))
}
func (this *GroupPermit) getQuerySeter() orm.QuerySeter {
	return this.getOrm().QueryTable(this)
}

func (this *GroupPermit) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}

//通过权限ID 删除用户组的权限关系
func (this *GroupPermit) DeleteByGroupIds(groupIds []string) (bool, error) {

	//删除表头信息
	_, err := this.getQuerySeter().Filter("permit_id__in", groupIds).Delete()
	if nil != err {
		return false, err
	}

	return true, err
}

//通过用户群权限ID 删除权限组的权限
func (this *GroupPermit) DeleteByPermitIds(permitIds []string) (bool, error) {

	//删除表头信息
	_, err := this.getQuerySeter().Filter("group_id__in", permitIds).Delete()
	if nil != err {
		return false, err
	}

	return true, err
}

/**
*
 */
func (this *GroupPermit) GetGroupPermitList(groupIds []string, uppermit_id []string) (*[]GroupPermit, error) {

	var permit Permit
	var groupPermitList []GroupPermit
	var where string
	var sliceParams []string
	// 构建查询对象

	nowTableName := dbPrefix + this.TableName()
	leftTableName := dbPrefix + permit.TableName()
	if len(groupIds) == 0 {
		return &groupPermitList, nil
	} else {
		where += nowTableName + ".group_id in (?)"
		sliceParams = append(sliceParams, "\""+strings.Join(groupIds, "\",\"")+"\"")
	}
	if len(uppermit_id) > 0 {
		where += " AND " + leftTableName + ".uppermit_id in (?)"
		sliceParams = append(sliceParams, "\""+strings.Join(uppermit_id, "\",\"")+"\"")
	}

	qb, _ := orm.NewQueryBuilder("mysql")

	sql := qb.Select("*").From(nowTableName).
		LeftJoin(leftTableName).On(nowTableName + ".permit_id=" + leftTableName + ".id").Where(where).OrderBy(leftTableName + ".obyid").Asc().String()
	_, err := this.getOrm().Raw(sql, sliceParams).QueryRows(&groupPermitList)
	if nil != err {
		panic(err)
	}
	return &groupPermitList, err
}
