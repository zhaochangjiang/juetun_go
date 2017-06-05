package common

//"Passport.Register"
//拼接字符串
func CreateUrl(controllerAction string, params map[string]interface{}, domain string) {
	returnString := ""
	if "" != domain {
		returnString += domain
	}
	for k, v := range params {
		switch(v.type())
		{
			case :""
			returnString += k + "=" + v;
		}
		
	}

	return returnString

}
