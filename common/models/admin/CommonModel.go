package admin

import (
	//	"juetun/common/utils"
	"github.com/astaxie/beego/cache"
)

var (
	bm cache.Cache
)

type CommonModel struct {
}

func (this *CommonModel) GetTablePrefix() string {
	var tablePrefix = "admin_"
	return tablePrefix
}
