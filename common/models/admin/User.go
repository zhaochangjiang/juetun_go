package admin

import (
	"errors"

	"github.com/astaxie/beego/orm"
)

type User struct {
	CommonModel
	Uid        string `orm:"column(uid);pk;"`
	Name       string `orm:varchar(30)`
	SuperAdmin string `orm:"column(super_admin)"`
	Isdel      string `orm:"column(isdel)"`
}

func (u *User) TableName() string {
	return "user"
}
func init() {
	user := new(User)
	orm.RegisterModelWithPrefix(user.GetTablePrefix(), user)
}
func (this *User) getOrm() orm.Ormer {
	o := orm.NewOrm()
	// 默认使用 default，你可以指定为其他数据库
	o.Using("db_admin")
	return o
}
func (this *User) getQuerySeter() orm.QuerySeter {
	return this.getOrm().QueryTable(this)
}

func (this *User) FetchUserById(userid string) (*User, error) {
	if "" == userid {
		return this, errors.New("userid is null")
	}
	err := this.getQuerySeter().Filter("uid", userid).One(this)
	return this, err
}
