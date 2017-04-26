package user

type UserMain struct {
	User_id           int
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
