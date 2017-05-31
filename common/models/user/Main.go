package user

import (
	"github.com/astaxie/beego/orm"
)

type Main struct {
	User_id           int     `orm:"column(user_id);pk;auto" json:"user_id"`
	Username          string  //用户名
	Name              string  //姓名
	Email             string  //邮箱
	Mobile            string  //手机号
	Password          string  //密码
	Avater            int     //头像ID
	Gender            string  //性别默认 male
	Step              int     //会员等级
	Score             int     //积分
	Gaode_mapx        float32 //高德地图坐标经度
	Gaode_mapy        float32 //高德地图坐标纬度
	Have_admin_permit int     //是否有客服后台权限0：没有，1:有
	Flag_del          string  //是否删除,''yes'':是，''no'':否
}

func init() {
	orm.RegisterModelWithPrefix("user_", new(Main))
}
func (this *Main) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_user") // 默认使用 default，你可以指定为其他数据库
	return o
}

/**
* 根据用户名获得用户信息
 */
func (this *Main) FetchUserByUserName(userName string) (Main, string) {
	var umain Main
	var message string
	o := this.getOrm()
	err := o.QueryTable(this).Filter("username", userName).Filter("flag_del", "no").One(&umain)
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		message = "data exception,please cotact the administrator!"
		return nil, message
	}
	return umain, message
}
