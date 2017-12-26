package admin

import (
	"encoding/json"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"juetun/common/utils"
	"time"
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

//根据key的前缀查询配置列表
func (this *Config) GetConfigByLikeKey(key string) *[]Config {
	var configList = make([]Config, 0)
	var keys = this.GetCachePrefix() + "_like_" + key
	var err error

	//读取缓存数据，如果缓存数据中为空，才从数据库读取。
	var ct = this.GetContentFromCache(keys)
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
	this.SetContentFromCache(keys, string(js1[:]), time.Second*3600)

	return &configList
}

//从redis缓存中读取相应的信息
func (this *Config) GetContentFromCache(key string) interface{} {

	var cache = *this.GetCache()

	return cache.Get(key)
}

//将数据放入redis缓存中。
func (this *Config) SetContentFromCache(key string, content interface{}, lifeTime time.Duration) {
	(*this.GetCache()).Put(key, content, lifeTime) //time.Second*3600
}

func (this *Config) GetRedisConfig() (string, string, string) {

	var cacheType = "redis"
	var redisName = "redisCache"
	var dbNumber = "0"
	return cacheType, redisName, dbNumber
}

//获得redis缓存中的内容
func (this *Config) GetCache() *cache.Cache {
	if nil == bm {

		cacheType, redisName, dbNumber := this.GetRedisConfig()
		var err error
		var redisconfig = new(utils.RedisConfig)
		var con = redisconfig.GetJsonCode(redisName, dbNumber)
		bm, err = cache.NewCache(cacheType, con)
		if nil != err {
			panic(err.Error())
		}
	}
	return &bm
}
