package general

import (
	"errors"

	"strings"

	"github.com/astaxie/beego"
)

//添加模板函数
func InitAddFuncMap() {
	beego.AddFuncMap("createurl", CreateUrl)
}

//创建URL
func CreateUrl(p ...interface{}) string {
	length := len(p)
	var url string
	switch length {
	case 1:
		p0, err := convertInterfaceToString(p[0])
		if nil == err {
			panic(err.Error())
		}
		url += p0
		break
	case 2:
		p0, err := convertInterfaceToString(p[0])
		if nil != err {
			panic(err.Error())
		}
		p1, err := convertInterfaceToString(p[1])
		if nil != err {
			panic(err.Error())
		}

		url += p0 + "/" + p1
		break
	case 3:

		url += getThreeParams(p)
	case 4:
		url += getThreeParams(p)
		domain, err := convertInterfaceToString(p[3])
		if nil != err {
			panic(err.Error())
		}
		domainStirng := beego.AppConfig.String(beego.BConfig.RunMode + "::" + domain)
		if url != "///" {
			url = domainStirng + url
		} else {
			url = "javascript:void(0);"
		}
	default:
		panic("now CreateUrl params length must be 1-4")
	}
	return url
}
func getThreeParams(p []interface{}) string {
	url := "/"
	for i := 0; i < 2; i++ {
		p0, err := convertInterfaceToString(p[i])
		if nil != err {
			panic(err.Error())
		}
		url += p0 + "/"
	}
	switch p[2].(type) { //多选语句switch
	case string: //是字符时做的事情
		url += p[2].(string)
	case map[string]string:
		params := make([]string, 0)
		p2 := p[2].(map[string]string)
		for k, v := range p2 {
			params = append(params, k+"_"+v)
		}
		url += strings.Join(params, "&")
	default:
		panic("The arguments is error!")
	}
	return url
}
func convertInterfaceToString(p interface{}) (string, error) {
	var c string
	var e error
	switch p.(type) {
	case string:
		c = p.(string)
		return c, nil
	default:
		e = errors.New("you send params type must be string")
		return c, e
	}

}
