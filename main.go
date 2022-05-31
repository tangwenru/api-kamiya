package main

import (
	"api-kamiya/middleware"
	_ "api-kamiya/routers"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"runtime"
	"time"
	//beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"html/template"
)

var AppConfig = beego.AppConfig

//type Base struct {
//	orm orm.Ormer
//}

func main() {
	//InsertFilter是提供一个过滤函数
	beego.InsertFilter("*", beego.BeforeRouter, middleware.MiddlewareAllow)

	//beego.BConfig.WebConfig.Session.SessionOn = true

	//获取 token
	beego.InsertFilter("*", beego.BeforeRouter, middleware.MiddlewareAll)

	//beego.InsertFilter("/user/*", beego.BeforeRouter, middleware.MiddlewareUser)

	beego.ErrorHandler("404", page404)

	//fmt.Println("runtime.GOOS", runtime.GOOS )

	runMode := beego.BConfig.RunMode
	//如果是 linux
	if runtime.GOOS == "linux" {
		//local 防止本地代码，发到线上，
		if runMode == "local" {
			runMode = "prod"
			beego.BConfig.RunMode = "prod"
			AppConfig.Set("RunMode", runMode)
		}
	}

	if runMode == "dev" || runMode == "local" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	orm.RegisterDriver("mysql", orm.DRMySQL)

	//"username:password@tcp(127.0.0.1:3306)/db_name?charset=utf8", 30
	var MysqlUrls, _ = AppConfig.String("MysqlUrls")
	var MySqlDbName, _ = AppConfig.String("MySqlDbName")
	var MySqlUser, _ = AppConfig.String("MySqlUser")
	var MySqlPass, _ = AppConfig.String("MySqlPass")
	var MySqlCharset, _ = AppConfig.String("MySqlCharset")

	orm.RegisterDataBase(
		"default",
		"mysql",
		MySqlUser+`:`+MySqlPass+`@tcp(`+MysqlUrls+`)/`+MySqlDbName+`?charset=`+MySqlCharset,
	)

	// 设置为 UTC 时间
	orm.DefaultTimeLoc = time.UTC

	if runMode == "dev" || runMode == "local" {
		orm.Debug = true
	}

	//orm.RegisterModel( new( models.Cartoon ) )

	//beego.SessionConfig.SessionProvider = ""

	beego.Run()
}

func page404(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("404.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/404.html")
	data := make(map[string]interface{})
	data["content"] = "page not found"
	t.Execute(rw, data)
}
