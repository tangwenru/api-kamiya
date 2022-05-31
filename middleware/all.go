package middleware

import (
	"github.com/beego/beego/v2/server/web/context"
)

func MiddlewareAll(ctx *context.Context) {
	if ctx.Request.Method == "OPTIONS" {
		ctx.WriteString("")
	}
}
