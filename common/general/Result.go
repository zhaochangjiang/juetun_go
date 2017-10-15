package general

type Result struct {
	Code    int         `json:"code"`    //返回的状态码
	Data    interface{} `json:"data"`    //返回的数据描述
	Message string      `json:"message"` //信息提示
}
