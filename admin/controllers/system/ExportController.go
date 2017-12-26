package system

import (
	acommon "juetun/admin/common"
	"juetun/common/general"
	modelsAdmin "juetun/common/models/admin"
)

type ExportController struct {
	acommon.AdminController
	Model  modelsAdmin.Export
	Result general.Result
}

//开始导出
func (this *ExportController) Start() {

	this.Result = *this.Model.AddExport()
}

//查看进度
func (this *ExportController) Process() {

}

//取消下载
func (this *ExportController) Cancel() {

}

//当前导出的列表
func (this *ExportController) List() {

}

//删除导出信息
func (this *ExportController) Del() {

}
