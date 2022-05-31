package controllersCrawlerKuaishou

import (
	configCrawler "api-kamiya/config/crawler"
	"api-kamiya/controllers"
	"api-kamiya/global"
	crawlerKuaishou "api-kamiya/models/crawler/kuaishou"
	"encoding/json"
)

type VideoController struct {
	controllers.BaseController
}

// 记录
func (this *VideoController) CreateVideoList() {
	query := []configCrawler.CreateKuaiShouVideoListQuery{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &query)
	if err != nil {
		this.Json(global.DataResultError("-1", "请求参数有问题", ""))
		return
	}
	kuaiShou := crawlerKuaishou.KuaiShouVideo{}
	data := kuaiShou.CreateVideoList(&query)
	this.Json(data)
}

func (this *VideoController) CreateVideoListHandle() {
	query := []configCrawler.CreateKuaiShouVideoListQuery{}
	dataPost := []byte(VideoData)
	err := json.Unmarshal(dataPost, &query)
	if err != nil {
		this.Json(global.DataResultError("-1", "请求参数有问题", err.Error()))
		return
	}
	kuaiShou := crawlerKuaishou.KuaiShouVideo{}
	data := kuaiShou.CreateVideoList(&query)
	this.Json(data)
}
