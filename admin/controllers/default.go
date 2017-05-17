package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Layout = "layout/main.html"
	this.TplName = "default/index.html"
}
