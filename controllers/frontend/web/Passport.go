package web

import (
	"juetun/controllers/frontend/common"
)

type PassportController struct {
	CFCommon.WebController
}

func (this *PassportController) Login() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie111111111@gmail.com"

	//渲染文件
	this.Display("login.html")

}

//登录提交
func (this *PassportController) IframeLogin() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie111111111@gmail.com"

	//渲染文件
	this.Display("login.html")

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
