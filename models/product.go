package models

import (
	"api-kamiya/global"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type Product struct {
	Id                int64
	Name              string
	ThumbUrl          string
	ProductType       string
	HighIntroduction  string
	Introduction      string
	Content           string
	ScreenshotPreview string
	Tag               string
	Created           int64
	Updated           int64
	Base
}

type ProductDetail struct {
	Name              string   `json:"name"`
	ThumbUrl          string   `json:"thumbUrl"`
	ProductType       string   `json:"productType"`
	HighIntroduction  string   `json:"highIntroduction"`
	Introduction      string   `json:"introduction"`
	Content           string   `json:"content"`
	ScreenshotPreview []string `json:"screenshotPreview"`
	Tag               []string `json:"tag"`
}

func init() {
	orm.RegisterModel(new(Product))
}

func (this *Product) GetQueryTable() orm.QuerySeter {
	return this.orm().QueryTable(this)
}

func (this *Product) Detail(userId int64, productType string) global.DataResultModel {
	result := global.GetDataResultModel()

	product := Product{}
	err := this.orm().
		QueryTable(this).
		Filter("productType", productType).
		Limit(1).
		One(&product)

	if err != nil {
		fmt.Println("product Detail err:", userId, err)
	}

	detail := ProductDetail{
		Name:              product.Name,
		ThumbUrl:          product.ThumbUrl,
		ProductType:       product.ProductType,
		HighIntroduction:  product.HighIntroduction,
		Introduction:      product.Introduction,
		Content:           product.Content,
		ScreenshotPreview: global.String2Array(product.ScreenshotPreview, ","),
		Tag:               global.String2Array(product.Tag, ","),
	}

	result.Success = product.Id > 0
	result.Data = detail

	return result
}

func (this *Product) List() []Product {
	list := []Product{}
	sqlQuery := this.GetQueryTable()

	_, err := sqlQuery.
		OrderBy("id").
		All(&list)

	if err != nil {

	}
	return list
}
