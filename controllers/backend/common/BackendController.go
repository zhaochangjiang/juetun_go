package CBCommon

import (
	"juetun/common"
)

type BackendController struct {
	common.BaseController
	userid         int64
	username       string
	nickname       string
	controllerName string
	actionName     string
	backend        string
}
