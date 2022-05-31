package global

import (
	beego "github.com/beego/beego/v2/adapter"
)

func Url( api string ) string {
	return beego.AppConfig.String("HttpApiDomain") + api
}
