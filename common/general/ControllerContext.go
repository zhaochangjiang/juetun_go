package general

type IncludePageProperty struct {
	HaveTable    bool
	HaveCheckbox bool
}
type Breadcrumb struct {
	Href   string
	FaCss  string
	Name   string
	Active bool
}

type ControllerContext struct {
	IncludePageProperty
	GroupIds              []string //当前用户所属用户组
	NotNeedLogin          bool     //是否需要登录
	IsSuperAdmin          bool     //是否为超级管理员,
	NotNeedValidatePermit bool     //是否需要验证权限(true:不需要验证权限)
	Controller            string
	Action                string
	JsFileAfter           []string     //页面尾部JS文件
	JsFileBefore          []string     //页面头部JS文件
	CssFile               []string     //样式文件
	NeedRenderJs          bool         //是否需要引入必要的渲染JS文件，controller_action.js文件,
	Breadcrumbs           []Breadcrumb //面包屑对象,
	Permit                *map[string]interface{}
	OutputResult          Result
}
