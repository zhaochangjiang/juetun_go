package general

type ControllerContext struct {
	GroupIds              []string //当前用户所属用户组
	NotNeedLogin          bool     //是否需要登录
	IsSuperAdmin          bool     //是否为超级管理员,
	NotNeedValidatePermit bool     //是否需要验证权限(true:不需要验证权限)
}
