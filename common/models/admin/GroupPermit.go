package admin

import (
	"strings"

	"github.com/astaxie/beego/orm"
)

type GroupPermit struct {
	CommonModel
	Id       string `orm:"column(id);pk" json:"id"`
	PermitId string `orm:"column(permit_id)"`
	GroupId  string `orm:"column(group_id)"`
}

func (this *GroupPermit) TableName() string {
	return "grouppermit"
}

func init() {
	model := new(GroupPermit)
	orm.RegisterModelWithPrefix(model.GetTablePrefix(), model)
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
* 根据用户组ID获得用户的权限列表
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/14
* @param  groupIds 						[]string
* @param   uppermit_id 					[]string
* @return  *[]PermitAdmin, error
 */
func (this *GroupPermit) GetGroupPermitList(groupIds []string, uppermit_id []string) (*[]PermitAdmin, error) {

	var permit Permit
	var permitList []Permit
	var where string
	var sliceParams []string
	result := make([]PermitAdmin, 0)
	qb, _ := orm.NewQueryBuilder("mysql")

	// 构建查询对象
	tablePrefix := this.GetTablePrefix()
	nowTableName := tablePrefix + this.TableName()
	leftTableName := tablePrefix + permit.TableName()

	if len(groupIds) == 0 {
		return &result, nil
	} else {
		where += nowTableName + ".group_id in (\"" + strings.Join(groupIds, "\",\"") + "\")"
	}
	if len(uppermit_id) > 0 {
		where += " AND " + leftTableName + ".uppermit_id in (\"" + strings.Join(uppermit_id, "\",\"") + "\")"
	}

	sql := qb.Select(leftTableName + ".*").From(nowTableName).LeftJoin(leftTableName).On(nowTableName + ".permit_id=" + leftTableName + ".id").Where(where).OrderBy(leftTableName + ".obyid").Asc().String()

	_, err := this.getOrm().Raw(sql, sliceParams).QueryRows(&permitList)
	if nil != err {
		panic(err)
	}

	for _, v := range permitList {
		params := make(map[string]string)
		params["mod"] = v.Mod
		result = append(result, *permit.OrgAdminPermit(v, params))
	}
	return &result, err
}
