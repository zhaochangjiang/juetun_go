package general

import (
	"log"

	"github.com/astaxie/beego"
)

//添加模板函数
func InitAddFuncMap() {
	beego.AddFuncMap("createurl", CreateUrl)
}

//创建URL
func CreateUrl(p ...interface{}) string {
	log.Println(p)
	return "adfasdf"
}
