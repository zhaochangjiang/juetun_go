package data

import (
	acommon "juetun/admin/common"
)

type UserController struct {
	acommon.AdminController
}

func (this *UserController) List() {
	this.LoadCommon("layout/list.html")
}
