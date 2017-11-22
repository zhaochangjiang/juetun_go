package general

import (
	"strings"

	"github.com/astaxie/beego/logs"
	_ "github.com/astaxie/beego/session/redis"
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

//系统初始化功能
func init() {
	//	logs.Debug("my book is bought in the year of ", 2016)
	//	logs.Info("this %s cat is %v years old", "yellow", 3)
	//	logs.Warn("json is a type of kv like", map[string]int{"key": 2016})
	//	logs.Error(1024, "is a very", "good game")
	logs.SetLogger(logs.AdapterFile, `{"filename":"success.log","level":7,"maxlines":10000,"maxsize":0,"daily":true,"maxdays":10}`)
	logs.Async() //为了提升性能, 可以设置异步输出:

	logs.Async(1e3) //异步输出允许设置缓冲 chan 的大小
	//初始化语言包
	InitLanguage()
	//初始化数据库
	InitDatabase()

	//初始化登录信息
	initSession()

	//初始化模板函数
	InitAddFuncMap()
}

//初始化登录信息
func initSession() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionProvider = "redis"
	beego.BConfig.WebConfig.Session.SessionName = "jtid"
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	//设置session reddis连接地址和密码。
	beego.BConfig.WebConfig.Session.SessionProviderConfig = "127.0.0.1:6379,123456"
}

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

		if err := i18n.SetMessage(lang, "../common/conf/language/"+lang+"/locale_"+lang+".ini"); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			return
		}
	}
}
