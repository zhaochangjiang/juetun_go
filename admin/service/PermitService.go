package service

import (
	"juetun/common/models/admin"
)

type PermitService struct {
	admin.Permit
}

func (this *PermitService) PermitService(*[]*admin.Permit) interface{} {
	var result interface{}
}
