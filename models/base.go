package models

import (
	beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type ResultData struct {
	State   string `json:"state"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	//Time string `json:"time"`
	Data interface{} `json:"data"`
}

// 远程i接受的数据
type ResponseData struct {
	State   string `json:"state"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	//Time string `json:"time"`
	Data interface{} `json:"data"`
}

type Base struct {
}

func Init() {
	//// 设置为 UTC 时间
	//orm.DefaultTimeLoc = time.UTC;
}

func (b *Base) Orm() (o orm.Ormer) {
	return orm.NewOrm()
}

func (b *Base) orm() (o orm.Ormer) {
	return orm.NewOrm()
}

func (b *Base) GetOrm() (o orm.Ormer) {
	return orm.NewOrm()
}

func (b *Base) GetQB() orm.QueryBuilder {
	// 获取 QueryBuilder 对象. 需要指定数据库驱动参数。
	// 第二个返回值是错误对象，在这里略过
	qb, err := orm.NewQueryBuilder("mysql")
	log := BusinessLog{}
	if err != nil{
		log.Create(0, "error", "base get query builder ", "", err )
	}
	return qb
}

func (b *Base) GetTableName(str string) string {
	return TableName( str )
}


//func ( this *Base) GetQueryTable() orm.QuerySeter  {
//	return this.orm().QueryTable( this );
//}

//返回带前缀的表名
func TableName(str string) string {
	return beego.AppConfig.String("MySqlPrefix") + str
}