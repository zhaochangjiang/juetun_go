package system

import (
	acommon "juetun/admin/common"
)

type SwitchController struct {
	acommon.AdminController
}

func (this *SwitchController) Website() {
	this.LoadCommon("system/switch/website.html")
}
