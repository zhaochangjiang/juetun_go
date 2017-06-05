package admin

import (
	"github.com/astaxie/beego/orm"
)

//表头信息存储表
type TableSetting struct {
	Id string `orm:"column(id);varchar(32);pk" json:"id"`

	//字段名称
	FieldName string `orm:"varchar(100)";orm:"column(field_name)`
	//字段类型、
	Type string `orm:"varchar(100)";orm:"column(type)`
	//所属表格类型,每个展示表格以本字段做标记
	LabelType string `orm:"varchar(100)";orm:"column(label_type)`

	//所属模块（当前主要ENUM("admin","web","user")
	ModuleType string `orm:"varchar(100)";orm:"column(module_type)`

	//表头样式class
	CssClass string `orm:"varchar(100)";orm:"column(css_class)`

	//内容样式（如果为空，则与CssClass 值一致）
	TextClass string `orm:"varchar(100)";orm:"column(text_class)`
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(TableSetting))
}
func (u *TableSetting) TableName() string {
	return "table_setting"
}
func (this *TableSetting) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}

//删除表头信息
func (this *TableSetting) DeleteMessageById(id []string) (bool, error) {

	//删除相关的数据
	tableSettingUser := new(TableSettingUser)
	tableSettingUser.DeleteByTableSettingId(id)

	//删除表头信息
	_, err := this.getOrm().QueryTable(this.TableName()).Filter("id__in", id).Delete()
	if nil != err {
		return false, err
	}

	return true, err
}
