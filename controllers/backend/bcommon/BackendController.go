package bcommon

import (
	"juetun/general"
)

type BackendController struct {
	general.BaseController
	userid         int64
	username       string
	nickname       string
	controllerName string
	actionName     string
	backend        string
}
