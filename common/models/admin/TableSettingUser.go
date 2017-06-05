package admin

import (
	"github.com/astaxie/beego/orm"
)

//表头信息存储表,用户偏好设置表头信息表
type TableSettingUser struct {
	Id string `orm:"column(id);varchar(40);pk" json:"id"`
	//用户ID
	User_Id string `orm:"varchar(40)";orm:"column(user_id)"`
	//table_setting 表ID
	TableSettingId string `orm:"varchar(40)";orm:"column(table_setting_id)"`
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(TableSettingUser))
}
func (this *TableSettingUser) TableName() string {
	return "table_setting_user"
}
func (this *TableSettingUser) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}

//删除信息
func (this *TableSettingUser) DeleteByTableSettingId(tableSettingId []string) (bool, error) {
	//删除表头信息
	_, err := this.getOrm().QueryTable(this.TableName()).Filter("table_setting_id__in", tableSettingId).Delete()
	if nil != err {
		return false, err
	}
	return true, err
}
