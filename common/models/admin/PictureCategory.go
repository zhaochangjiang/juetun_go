package admin

import (
	"github.com/astaxie/beego/orm"
)

type PictureCategory struct {
	PictureCategoryId   string `orm:"column(picture_category_id);pk;auto" json:"picture_category_id"`
	PictureCategoryname string `orm:"column(picture_categoryname)"`
	PictureCategorykey  string `orm:"column(picture_category_key)"`
	PictureSavepath     string `orm:"column(picture_savepath)"`
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(PictureCategory))
}
func (u *PictureCategory) TableName() string {
	return "picture_category"
}
func (this *PictureCategory) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}
