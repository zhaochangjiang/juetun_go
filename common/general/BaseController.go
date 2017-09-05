package general

import (
	modelsUser "juetun/common/models/user"

	"log"
	"strings"

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

func (this *BaseController) Debug(data interface{}) {
	log.Println("")
	log.Println("")
	log.Println("--------------------Debug----------------------")

	log.Print(data)
	log.Println("------------------- Debug-----------------------")
	log.Println("")
	log.Println("")
}

/**
* 获得域名MAP key列表
* @author karl.zhao<zhaocj2009@126.com>
* @date 2017/08/24
 */
func (this *BaseController) GetDomainMapKey() *[]string {
	domainMapArray := beego.AppConfig.Strings("domainmap")
	return &domainMapArray
}

/**
* 获得域名MAP对应的域名地址
* @author karl.zhao<zhaocj2009@126.com>
* @date 2017/08/24
 */
func (this *BaseController) GetDomainMapList() *map[string]string {
	domainMapArray := this.GetDomainMapKey()
	var returnContent = make(map[string]string)
	for _, dataName := range *domainMapArray {
		returnContent[dataName] = beego.AppConfig.String(beego.BConfig.RunMode + "::domain_" + dataName)
	}
	return &returnContent
}

//公共的
//func (this *BaseController) Common() {
//	w := this.Ctx.ResponseWriter
//	//Origin := this.Header.Get("Origin")
//	Origin := "*.test.com"
//	if Origin != "" {
//		w.Header().Add("Access-Control-Allow-Origin", Origin)
//		w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
//		w.Header().Add("Access-Control-Allow-Headers", "x-requested-with,content-type")
//		w.Header().Add("Access-Control-Allow-Credentials", "true")

//		fmt.Println("safasdfsdf")
//	}
//}
func (this *BaseController) DisplayIframe(message string) {
	this.Data["Message"] = message
	this.TplName = "iframe.html"
}

func (this *BaseController) UserDataDefault() *modelsUser.Main {

	userMain := new(modelsUser.Main)
	avater := this.GetSession("Avater")
	gender := this.GetSession("Gender")
	user_id := this.GetSession("User_id")
	username := this.GetSession("Username")

	if nil != avater {
		userMain.Avater = avater.(string)
	}
	if nil != gender {
		userMain.Gender = gender.(string)
	}
	if nil != user_id {
		userMain.User_id = user_id.(string)
	}
	if nil != username {
		userMain.Username = username.(string)
	}
	//设置默认头像
	if "" == userMain.Avater {
		//判断性别
		if userMain.Gender == "female" {
			userMain.Avater = "/assets/img/avatar3.png"
		} else {
			userMain.Avater = "/assets/img/avatar5.png"
		}
	} else {
		userMain.Avater = this.GetAvaterByPictureId(userMain.Avater)
	}
	this.Data["Username"] = userMain.Username
	this.Data["Avater"] = userMain.Avater
	this.Data["Uid"] = userMain.User_id
	return userMain
}

/**
* 获得用户头像
 */
func (this *BaseController) GetAvaterByPictureId(pictureId string) string {

	return "/assets/img/avatar3.png"
}

// Prepare implemented Prepare method for baseRouter.
func (this *BaseController) Prepare() {

	this.Database = orm.NewOrm()
	// Setting properties.
	//	this.Data["AppVer"] = AppVer
	//	this.Data["IsPro"] = IsPro
	//	this.Data["PageStartTime"] = time.Now()
	//	this.Data["PageTitle"] = ""
	//	this.Data["PageKeyword"] = ""
	//	this.Data["PageDescription"] = ""
	//	this.Data["PageAuthor"] = ""

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
