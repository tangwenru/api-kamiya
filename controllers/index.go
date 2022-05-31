package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/adapter"
)

// Operations about Users
type IndexController struct {
	beego.Controller
}

func (p *IndexController) Prepare() {
	fmt.Println("ori")
}

func (u *IndexController) Get() {

	u.Data["json"] = beego.AppConfig.String("AppTitle")
	//
	//u.RenderString();

	u.ServeJSON()
}
