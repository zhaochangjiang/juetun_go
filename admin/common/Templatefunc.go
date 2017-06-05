package common

//"Passport.Register"
//拼接字符串
func CreateUrl(controllerAction string, params map[string]interface{}, domain string) string {
	returnString := ""
	if "" != domain {
		returnString += domain
	}
	for k, v := range params {

		returnString += (k + "=" + v.(string))
		//		switch tp := v.(type) {
		//		case int:
		//			returnString += (k + "=" + strconv.Itoa(v))
		//		case int64:
		//			returnString += k + "=" + (strconv.FormatInt(v, 10))
		//		case string:
		//			returnString += k + "=" + v
		//		default:
		//			panic("操作错误")
		//		}

	}

	return returnString

}
