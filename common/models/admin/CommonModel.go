package admin

type CommonModel struct {
}

func (this *CommonModel) GetTablePrefix() string {
	var tablePrefix = "admin_"
	return tablePrefix
}
