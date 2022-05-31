package models

import (
	"api-kamiya/global"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type Article struct {
	Id         int64 `orm:"primary"`
	ClassifyId int64
	Title      string
	Content    string
	Created    int64
	Updated    int64
	Base
}

type ArticleDetail struct {
	IdKey      string `json:"idKey"`
	ClassifyId int64  `json:"classifyId"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Created    int64  `json:"created"`
	Updated    int64  `json:"updated"`
	Base
}

func init() {
	orm.RegisterModel(new(Article))
}

func (this *Article) Detail(userId int64, idKey string) global.DataResultModel {
	result := global.GetDataResultModel()
	id := global.IdDecrypt(idKey)

	detail := Article{}
	err := this.orm().
		QueryTable(this).
		Filter("id", id).
		Limit(1).
		One(&detail)

	if err != nil {
		fmt.Println("aticle Detail err:", err)
	}

	groupDetail := ArticleDetail{
		IdKey:      global.IdEncrypt(detail.Id),
		ClassifyId: detail.ClassifyId,
		Title:      detail.Title,
		Content:    detail.Content,
		Created:    detail.Created,
		Updated:    detail.Updated,
	}

	result.Success = detail.Id > 0
	if !result.Success {
		result.Message = "内容不存在"
		result.Code = "404"
	}
	result.Data = groupDetail
	return result
}
