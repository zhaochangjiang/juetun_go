package general

type Result struct {
	Code    int         //返回的状态码
	Data    interface{} //返回的数据描述
	Message string      //信息提示
}
