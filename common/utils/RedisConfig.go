package utils

import (
	"encoding/json"

	"github.com/astaxie/beego"
)

type RedisConfig struct {
	Conn     string ` json:"conn"`
	Password string `json:"key"`
	DbNum    string `json:"dbNum"`
}

/**
* 将数据转换为json字符串
* @params name  配置文件redis连接的名称
* @params dbnum redis数据库编号一般默认传“0”
* @return 返回字符串
 */
func (this *RedisConfig) GetJsonCode(name string, dbnum string) string {
	this.Conn = beego.AppConfig.String(name+"::host") + ":" + beego.AppConfig.String(name+"::port")
	this.Password = beego.AppConfig.String(name + "::pwd")
	s, err := json.Marshal(this)
	if nil != err {
		panic(err.Error())
	}
	return string(s[:])
}
