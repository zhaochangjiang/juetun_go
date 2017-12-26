package dashboard

type Permit struct {
	CommonModel
	Id         string `orm:"column(id);pk" json:"id";form:"-"`
	Name       string `orm:varchar(50);orm:"column(name)";form:"name"`
	Mod        string `orm:varchar(30);orm:"column(mod)";form:"mod"`
	Controller string `orm:varchar(30);orm:"column(controller)";form:"controller"`
	Action     string `orm:varchar(30);orm:"column(action)";form:"action"`
	UppermitId string `orm:int(10);orm:"column(uppermit_id)";form:"uppermit_id"`
	DomainMap  string `orm:"column(domain_map)";form:"domain_map"`
	Obyid      string `orm:"column(obyid)";form:"obyid"`
	Csscode    string `orm:varchar(500);orm:"column(csscode)";form:"csscode"`
}

//获得子权限
func (this *Permit) GetChildPermit(args *map[string]string) {
	if args["isSuperAdmin"] == "yes" {
		this.GetSuperChildPermit(args)
	} else {
		this.GetNotSuperChildPermit(args)
	}
}

//如果不是超级管理员的子权限
func (this *Permit) GetNotSuperChildPermit(args *map[string]string) {
}

//如果是超级管理员的子权限
func (this *Permit) GetSuperChildPermit(args *map[string]string) *[]Permit {
	var permitList = make([]Permit, 0)
	_, err := this.getQuerySeter().qs.Filter("name__istartswith", args["parentString"]).All(&permitList)
	if err != nil {
		panic(err)
	}
	return &permitList
}
func (this *Permit) GetSuperChildPermitHasDealData(args *map[string]string) {
	var permitList = this.GetSuperChildPermit(args)
	var upPemitId = make(map[string][]Permit)
	for _, v := range *permitList {
		upPemitId[v.UppermitId] = append(upPemitId[v.UppermitId], v)
	}
}

/****************以下是基础公共配置********************/
func init() {
	permit := new(Permit)
	orm.RegisterModelWithPrefix(permit.GetTablePrefix(), permit)
}

func (this *Permit) TableName() string {
	return "permit"
}
func (this *Permit) getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("db_admin") // 默认使用 default，你可以指定为其他数据库
	return o
}

func (this *Permit) getQuerySeter() orm.QuerySeter {
	return this.getOrm().QueryTable(this)
}
