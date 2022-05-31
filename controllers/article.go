package controllers

import (
	"api-kamiya/models"
)

type ArticleController struct {
	BaseController
}

func (this *ArticleController) Detail() {
	article := models.Article{}
	idKey := this.GetString("idKey")
	data := article.Detail(this.GetUserId(), idKey)
	this.Json(data)
}
