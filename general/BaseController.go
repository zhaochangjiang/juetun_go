package general

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
)

var (
	AppVer    string
	IsPro     bool
	langTypes []*langType // Languages are supported.
)

//国际化语言包
type langType struct {
	Lang string
	Name string
}

type BaseController struct {
	beego.Controller
	i18n.Locale
	Database orm.Ormer //数据库操作对象
}

// Prepare implemented Prepare method for baseRouter.
func (this *BaseController) Prepare() {

	this.Database = orm.NewOrm()
	// Setting properties.
	this.Data["AppVer"] = AppVer
	this.Data["IsPro"] = IsPro

	y := time.Now().Year()
	this.Data["Copyright"] = "Copyright 2016-" + strconv.Itoa(y) + " " + beego.AppConfig.String("appname") + " Corporation. All Rights Reserved."

	this.Data["PageStartTime"] = time.Now()

	// Redirect to make URL clean.
	if this.setLangVer() {
		i := strings.Index(this.Ctx.Request.RequestURI, "?")
		this.Redirect(this.Ctx.Request.RequestURI[:i], 302)
		return
	}
}

// setLangVer sets site language version.
func (this *BaseController) setLangVer() bool {
	isNeedRedir := false
	hasCookie := false

	// 1. Check URL arguments.
	lang := this.Input().Get("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = this.Ctx.GetCookie("lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir = false
		hasCookie = false
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := this.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. Default language is English.
	if len(lang) == 0 {
		lang = "en-US"
		isNeedRedir = false
	}

	curLang := langType{
		Lang: lang,
	}

	// Save language information in cookies.
	if !hasCookie {
		this.Ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")
	}

	sliceLength := len(langTypes) - 1

	//判断如果长度小于0，则默认等于0
	if sliceLength < 0 {
		sliceLength = 0
	}

	//创建一个切片
	restLangs := make([]*langType, 0, sliceLength)

	for _, v := range langTypes {
		if lang != v.Lang {

			//给切片（restLangs）追加元素 V
			restLangs = append(restLangs, v)

		} else {
			curLang.Name = v.Name
		}
	}
	// Set language properties.
	this.Lang = lang
	this.Data["Lang"] = curLang.Lang
	this.Data["CurLang"] = curLang.Name
	this.Data["RestLangs"] = restLangs

	return isNeedRedir
}
