package middleware

import (
	"github.com/beego/beego/v2/server/web/context"
)

func MiddlewareAllow(ctx *context.Context) {
	//允许访问所有源
	ctx.Output.Header("Access-Control-Allow-Origin", "*")

	//可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
	//其中Options跨域复杂请求预检
	ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	ctx.Output.Header("Server", "ChuOS/1.8.2")



	//指的是允许的Header的种类
	ctx.Output.Header("Access-Control-Allow-Headers", "Origin, Authorization, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

	//公开的HTTP标头列表
	//	ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},

	//如果设置，则允许共享身份验证凭据，例如cookie
	//	AllowCredentials: false,

	//允许访问所有源
	//AllowAllOrigins: true,
	//
	////可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
	////其中Options跨域复杂请求预检
	//	AllowMethods: []string{"GET", "POST", "OPTIONS"},
	//
	//	//指的是允许的Header的种类
	//		AllowHeaders: []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "x-token"},
	//
	//	//公开的HTTP标头列表
	//		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	//
	//	//如果设置，则允许共享身份验证凭据，例如cookie
	//		AllowCredentials: false,
}
