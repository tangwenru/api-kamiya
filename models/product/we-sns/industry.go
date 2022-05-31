package modelsWeSns

import (
	"api-kamiya/config"
	"api-kamiya/global"
	"api-kamiya/models"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type Industry struct {
	Id        int64
	Name      string
	Enabled   bool
	OrderRank int64
	models.Base
}

func init() {
	orm.RegisterModel(new(Industry))
}

func (this *Industry) GetQueryTable() orm.QuerySeter {
	return this.Orm().QueryTable(this)
}

func (this *Industry) TableName() string {
	return models.TableName("we_sns_industry")
}

func (this *Industry) List(userId int64) global.DataResultModel {
	result := global.GetDataResultModel()

	list := []Industry{}
	sqlQuery := this.GetQueryTable()

	_, err := sqlQuery.
		OrderBy("order_rank").
		All(&list)

	if err != nil {
		fmt.Println("Industry list:", err)
	}

	outData := []config.IndustryList{
		{
			Id:   0,
			Name: "[行业]",
		},
	}
	for _, items := range list {
		outData = append(outData, config.IndustryList{
			Id:   items.Id,
			Name: items.Name,
		})
	}

	result.Data = outData
	result.Success = true
	return result
}