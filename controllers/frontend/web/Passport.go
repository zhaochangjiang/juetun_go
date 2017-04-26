package web

import (
	"fmt"
	"juetun/controllers/frontend/fcommon"
	"juetun/models/user"
)

type PassportController struct {
	fcommon.WebController
}

func (this *PassportController) Login() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie111111111@gmail.com"

	//渲染文件
	this.Display("login.html")

}

//登录提交
func (this *PassportController) IframeLogin() {
	this.Database.Using("default") // 默认使用 default，你可以指定为其他数据库

	userMain := new(user.UserMain)
	userMain.Username = "长江"

	fmt.Println(this.Database.Insert(userMain))

	//	this.Data["Website"] = "beego.me"
	//	this.Data["Email"] = "astaxie111111111@gmail.com"

	//	//渲染文件
	//	this.Display("login.html")

}

//忘记密码
func (this *PassportController) Forget() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie111111111@gmail.com"

	//渲染文件
	this.Display("login.html")

}

//注册
func (this *PassportController) Register() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie111111111@gmail.com"

	//渲染文件
	this.Display("login.html")

}
