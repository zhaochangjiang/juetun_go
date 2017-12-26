package data

import (
	acommon "juetun/admin/common"
)

type DataController struct {
	acommon.AdminController
}

/**
* 基础数据管理
* @author karl.zhao<zhaocj2009@hotmail.com>
* @date 2017/08/16
* @return void
 */
func (this *DataController) List() {

	this.LoadCommon("data/list.html")
}
