package controllers

import (
	acommon "juetun/admin/common"
	"juetun/common/general"
	modelAdmin "juetun/common/models/admin"
	"juetun/common/models/user"
	"time"
)

type Passport struct {
	acommon.AdminController
}

//实现本结构体的基本加载，本文件中所有的界面不需要验证登录
func (this *Passport) Prepare() {
	//设置本页面不需要登录
	this.NotNeedLogin = true
	this.AdminController.Prepare()
}

//登录页面
func (this *Passport) Login() {

	//渲染文件
	this.Layout = "layout/passport.html"
	this.TplName = "passport/login.html"
}

//登录提交
func (this *Passport) IframeLogin() {
	params, validateAuth := this.validateIframeLogin()
	if !validateAuth {
		return
	}

	modelUserMain := new(user.Main)
	var userMain user.Main
	var err error
	userMain, err = modelUserMain.FetchUserByUserName((*params)["userName"])
	if nil != err {
		this.DisplayIframe(err.Error())
		return
	}

	adminUser := new(modelAdmin.User)
	admin_user, err1 := adminUser.FetchUserById(userMain.User_id)
	if nil != err1 {
		this.DisplayIframe(err1.Error())
		return
	}
	encyption := new(general.PasswordEncyption)
	encyptionString := encyption.Sha1((*params)["pwd"])
	// 判断密码是否正确
	if userMain.Password != encyptionString {
		// 多条的时候报错
		this.DisplayIframe("请输入正确的账号和密码！")
		return
	}

	this.SetSession("Uid", userMain.User_id)
	this.SetSession("Username", userMain.Username)
	this.SetSession("Avater", "/assets/img/user.jpg")
	this.SetSession("Gender", userMain.Gender)
	this.SetSession("SuperAdmin", admin_user.SuperAdmin)

	//延迟500毫秒
	time.Sleep(time.Microsecond * 3000)
	this.Data["LocationHref"] = "/"
	this.DisplayIframe("")

	//	this.Redirect("/", 302)
	//渲染文件
	//this.DisplayIframe("密码正确")
}
func (this *Passport) validateIframeLogin() (*map[string]string, bool) {
	params := make(map[string]string, 0)
	params["userName"] = this.GetString("username")
	params["pwd"] = this.GetString("pwd")

	var pass = true
	if "" == params["userName"] {
		this.DisplayIframe("请输入账号！")
		return &params, false
	}

	if "" == params["pwd"] {
		this.DisplayIframe("请输入密码！")
		return &params, false
	}
	return &params, pass
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
