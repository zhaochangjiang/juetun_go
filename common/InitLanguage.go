package common

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

//@author karl.zhao<zhaocj2009@126.com>
//@since 2017/04/17
//初始化语言配置
func InitLanguage() {
	beego.AddFuncMap("i18n", i18n.Tr)

	// Initialized language type list.
	langs := strings.Split(beego.AppConfig.String("lang::types"), "|")
	names := strings.Split(beego.AppConfig.String("lang::names"), "|")
	langTypes := make([]*langType, 0, len(langs))
	for i, v := range langs {
		langTypes = append(langTypes, &langType{
			Lang: v,
			Name: names[i],
		})
	}

	for _, lang := range langs {

		if err := i18n.SetMessage(lang, "conf/language/"+lang+"/locale_"+lang+".ini"); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			return
		}
	}
}
