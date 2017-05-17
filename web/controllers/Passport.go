package controllers

import (
	"juetun/general"
	fcommon "juetun/web/common"
	"juetun/web/models/user"
)

type Passport struct {
	fcommon.WebController
}

//登录页面
func (this *Passport) Login() {

	//渲染文件

	this.Layout = "layout/passport.html"
	this.TplName = "passport/login.html"
}

//登录提交
func (this *Passport) IframeLogin() {
	userName := this.GetString("username")
	pwd := this.GetString("pwd")

	userMain := new(user.Main)
	umain, message := userMain.FetchUserByUserName(userName)

	if "" != message {
		this.DisplayIframe(message)
		return
	}
	encyption := new(general.PasswordEncyption)
	encyptionString := encyption.Sha1(pwd)
	// 判断密码是否正确
	if umain.Password != encyptionString {
		// 多条的时候报错
		this.DisplayIframe("请输入正确的账号和密码！")
		return
	}
	this.SetSession("uid", umain.User_id)
	this.Data["LocationHref"] = "/"
	this.DisplayIframe("")
	//	this.Redirect("/", 302)
	//渲染文件
	//this.DisplayIframe("密码正确")
}

//忘记密码
func (this *Passport) Forget() {

	this.Layout = "layout/passport.html"
	this.TplName = "passport/forget.html"
}

//验证码请求
func (this *Passport) Authcode() {
	authcode := new(general.AuthCode)
	authcode.GetImage(&this.BaseController.Controller)
}

//注册
func (this *Passport) Register() {
	this.Layout = "layout/passport.html"
	this.TplName = "passport/register.html"
}

//注册数据处理
func (this *Passport) IframeRegister() {
	authcode := new(general.AuthCode)
	validate := authcode.Validate(&this.BaseController.Controller)
	//如果验证失败
	if !validate {
		//渲染文件
		this.DisplayIframe("密码正确")
		return
	}

}
