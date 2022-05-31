package global



func NotLogin () DataResultModel {
	data := GetDataResultModel()
	data.Code = "not-login"
	data.Message = "登录失效，请重新登录"
	return data
}

