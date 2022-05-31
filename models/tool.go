package models

import (
	beego "github.com/beego/beego/v2/adapter"
	_ "github.com/go-sql-driver/mysql"
)

var AppConfig = beego.AppConfig

type Tool struct {
}

func ApiUrl(api string) string {
	HttpApiDomain := AppConfig.String("HttpApiDomain")

	return HttpApiDomain + api
}
