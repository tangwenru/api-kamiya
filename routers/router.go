package routers

import (
	"api-kamiya/controllers"
	controllersCrawlerKuaishou "api-kamiya/controllers/crawler/kuaishou"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.IndexController{})
	web.AutoRouter(&controllers.ProductController{})

	web.AutoRouter(&controllers.ArticleController{})

	nsCrawler :=
		web.NewNamespace("/crawler",
			web.NSAutoRouter(&controllersCrawlerKuaishou.VideoController{}),
		)

	//注册 namespace
	web.AddNamespace(nsCrawler)
}
