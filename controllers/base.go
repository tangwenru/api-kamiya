package controllers

import (
	"api-kamiya/config"
	"api-kamiya/global"
	"api-kamiya/libs"
	"encoding/json"
	beego "github.com/beego/beego/v2/adapter"
	"net"
	"regexp"
	"strconv"
)

type BaseController struct {
	beego.Controller
}

func (p *BaseController) Prepare() {
	//controllerName, actionName := p.GetControllerAndAction()
	//p.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	//p.actionName = strings.ToLower(actionName)
	//p.o = orm.NewOrm();
	//if strings.ToLower( p.controllerName) == "admin" && strings.ToLower(p.actionName)  !=  "login"{
	//	if p.GetSession("user") == nil{
	//		p.History("未登录","/admin/login")
	//		//p.Ctx.WriteString(p.controllerName +"==="+ p.actionName)
	//	}
	//}
	//
	////初始化前台页面相关元素
	//if strings.ToLower( p.controllerName) == "blog"{
	//
	//	p.Data["actionName"] = strings.ToLower(actionName)
	//	var result []*models.Config
	//	p.o.QueryTable(new(models.Config).TableName()).All(&result)
	//	configs := make(map[string]string)
	//	for _, v := range result {
	//		configs[v.Name] = v.Value
	//	}
	//	p.Data["config"] = configs
	//}

}

func (this *BaseController) GetId(id string) int64 {
	idNumber, err := strconv.ParseInt(this.GetString(id), 10, 64)
	if err != nil {
		idNumber = 0
	}
	return idNumber
}

func (this *BaseController) GetUserId() int64 {
	userId := libs.GetRoleId("u", this.Ctx.Input.Query("_token_"))
	return userId
}

func (this *BaseController) GetRequestBody() []byte {
	return this.Ctx.Input.RequestBody
}

func (this *BaseController) GetStaffId() int64 {
	staffId := libs.GetRoleId("s", this.Ctx.Input.Query("_token_"))
	return staffId
}

func (this *BaseController) Json(data interface{}) {
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *BaseController) JsonSuccess(data interface{}) {
	this.Data["json"] = global.DataResultSuccess("", data)
	this.ServeJSON()
}
func (this *BaseController) JsonError(message string, data interface{}) {
	this.Data["json"] = global.DataResultError("-1", message, data)
	this.ServeJSON()
}

func (this *BaseController) GetPaginationQuery() config.Pagination {
	var pagination config.Pagination
	page := global.String2Int64(this.GetString("page"))
	if page == 0 {
		page = global.String2Int64(this.GetString("current"))
	}
	pageSize := global.String2Int64(this.GetString("pageSize"))

	// 没取到，就到 body 里面找；
	if page == 0 {
		body := this.Ctx.Input.RequestBody

		if string(body) == "" {
			pagination.Current = 1
			pagination.PageSize = 10
			return pagination
		}

		query := config.PaginationQuery{}
		err2 := json.Unmarshal(body, &query)
		if err2 == nil {
			page = query.Current
			pageSize = query.PageSize
		} else {
			page = 1
			pageSize = 10
		}
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	pagination.Current = page
	pagination.PageSize = pageSize

	return pagination
}

//
//func (this *BaseController) GetRequestBodyJsonData( query *interface{} ) error {
//	err := json.Unmarshal( this.Ctx.Input.RequestBody, &query )
//
//	return err
//}

func (this *BaseController) GetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	ip := ""
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}

	// 查询 https://lbs.qq.com/ 服务
	//局域网地址：192.168.*.* 需要鉴权
	match5, _ := regexp.MatchString(`192\.168\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])`, ip)

	if match5 || ip == "127.0.0.1" {
		ip = global.GetIp()
	}

	return ip
}
