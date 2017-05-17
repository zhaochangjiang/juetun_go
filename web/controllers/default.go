package controllers

import (
	fcommon "juetun/web/common"
)

// IndexController home controller
type Default struct {
	fcommon.WebController
}

// Get home page
func (this *Default) Get() {
	// 基础布局页面
	this.Layout = "layout/main.html"
	this.TplName = "default/index.html"
	this.Data["UserId"] = this.GetSession("uid")
}
