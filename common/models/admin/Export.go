package admin

import (
	"juetun/common/general"
)

type Export struct {
	CommonModel
	Id string `orm:"column(id);pk" json:"id";form:"-"`
}

func (this *Export) AddExport() *general.Result {
	res := new(general.Result)

	return res
}
