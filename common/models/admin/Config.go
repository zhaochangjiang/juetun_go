package admin

import (
	"encoding/json"
	"time"

	"juetun/common/utils"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
)

var (
	bm cache.Cache
)

type Config struct {
	CommonModel
	Id   string `orm:"column(id);pk" json:"id"`
	Name string `orm:varchar(30);orm:"column(name)"`
	Key  string `orm:varchar(50);orm:"column(key)"`
	Val  string `orm:varchar(1000);orm:"column(val)"`
}

func (u *Config) TableName() string {
	return "config"
}
func init() {
	var config = new(Config)
	orm.RegisterModelWithPrefix(config.GetTablePrefix(), config)
}
func (this *Config) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}
func (this *Config) getQuerySeter() orm.QuerySeter {
	return this.getOrm().QueryTable(this)
}
func (this *Config) GetCachePrefix() string {
	return "adminconfig"
}

func (this *Config) getCache() *cache.Cache {
	if nil == bm {
		var err error
		var redisconfig = new(utils.RedisConfig)
		var con = redisconfig.GetJsonCode("redisCache", "0")
		bm, err = cache.NewCache("redis", con)
		if nil != err {
			panic(err.Error())
		}
	}
	return &bm
}

//从redis缓存中读取相应的信息
func (this *Config) getContentFromCache(key string) interface{} {
	return (*this.getCache()).Get(key)
}

//将数据放入redis缓存中。
func (this *Config) setContentFromCache(key string, content interface{}) {
	(*this.getCache()).Put(key, content, time.Second*3600)
}

//根据key的前缀查询配置列表
func (this *Config) GetConfigByLikeKey(key string) *[]Config {
	var configList = make([]Config, 0)
	var keys = this.GetCachePrefix() + "_like_" + key
	var err error

	//读取缓存数据，如果缓存数据中为空，才从数据库读取。
	var ct = this.getContentFromCache(keys)
	if ct != nil {
		var c = ct.([]byte)
		j2 := make([]Config, 0)
		err = json.Unmarshal(c, &j2)
		if err != nil {
			panic(err)
		}
		return &j2
	}

	var querySelect = this.getQuerySeter()
	_, err = querySelect.Filter("key__istartswith", key).All(&configList)
	if nil != err {
		panic(err.Error())
	}

	//将数据转换为json字符串。
	js1, err := json.Marshal(configList)
	if err != nil {
		panic(err)
	}
	this.setContentFromCache(keys, string(js1[:]))

	return &configList
}
