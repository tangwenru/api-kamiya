package controllers

import (
	"api-kamiya/models"
)

// Operations about Users
type ProductController struct {
	BaseController
	trade models.Product
}

func (this *ProductController) Prepare() {
	this.trade = models.Product{}
}

func (this *ProductController) Detail() {
	detail := this.trade.Detail(this.GetUserId(), this.GetString("productType"))
	this.Json(detail)
}
