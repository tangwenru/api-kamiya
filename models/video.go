package models

import (
	"api-kamiya/global"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type Video struct {
	Id               int64
	Title            string
	VideoKey         string
	AnimatedCoverUrl string
	CoverUrl         string
	Duration         float64
	LikeCount        int64
	ViewCount        int64
	Platform         string
	Created          int64
	Updated          int64
	Base
}

type VideoDetail struct {
}

func init() {
	orm.RegisterModel(new(Video))
}

func (this *Video) GetQueryTable() orm.QuerySeter {
	return this.orm().QueryTable(this)
}

func (this *Video) Detail(userId int64, productType string) global.DataResultModel {
	result := global.GetDataResultModel()

	product := Video{}
	err := this.orm().
		QueryTable(this).
		Filter("productType", productType).
		Limit(1).
		One(&product)

	if err != nil {
		fmt.Println("product Detail err:", userId, err)
	}

	result.Success = product.Id > 0
	result.Data = product

	return result
}

func (this *Video) List() []Video {
	list := []Video{}
	sqlQuery := this.GetQueryTable()

	_, err := sqlQuery.
		OrderBy("id").
		All(&list)

	if err != nil {

	}
	return list
}
