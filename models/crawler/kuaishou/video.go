package crawlerKuaishou

import (
	configCrawler "api-kamiya/config/crawler"
	"api-kamiya/global"
	"api-kamiya/models"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type KuaiShouVideo struct {
	models.Video
	models.Base
}

func init() {

}

func (this *KuaiShouVideo) detailByVideoKey(videoKey string) models.Video {
	detail := models.Video{}
	err := this.GetQueryTable().
		Filter("video_key", videoKey).
		Limit(1).
		One(&detail)

	if err != nil {

	}
	return detail
}

func (this *KuaiShouVideo) CreateVideoList(list *[]configCrawler.CreateKuaiShouVideoListQuery) global.DataResultModel {
	result := global.GetDataResultModel()

	log := models.BusinessLog{}

	sqlQuery := this.GetQueryTable()
	i, _ := sqlQuery.PrepareInsert()

	now := time.Now().Unix()

	for _, items := range *list {
		for _, item := range items.Data.VisionProfilePhotoList.Feeds {
			detail := this.detailByVideoKey(item.Photo.Id)
			// 还不存在
			if detail.Id == 0 {
				insertData := models.Video{
					Title:            item.Photo.Title,
					VideoKey:         item.Photo.Id,
					AnimatedCoverUrl: item.Photo.AnimatedCoverUrl,
					CoverUrl:         item.Photo.CoverUrl,
					Duration:         float64(item.Photo.Duration),
					LikeCount:        item.Photo.LikeCount,
					ViewCount:        global.String2Int64(item.Photo.ViewCount),
					Platform:         "kuaishou",
					Created:          now,
				}

				_, insertError := i.Insert(&insertData)

				if insertError != nil {
					log.Create(0, "error", "kuaishou video-insert", insertData, insertError)
				}
			} else {
				// 更新
				_, updateError := this.GetQueryTable().
					Filter("id", detail.Id).
					Limit(1).
					Update(orm.Params{
						"Title":    item.Photo.Title,
						"VideoKey": item.Photo.Id,
						//"AnimatedCoverUrl": item.Photo.AnimatedCoverUrl,
						//"CoverUrl":         item.Photo.CoverUrl,
						"Duration":  float64(item.Photo.Duration),
						"LikeCount": item.Photo.LikeCount,
						"ViewCount": global.String2Int64(item.Photo.ViewCount),
						"Updated":   now,
					})
				if updateError != nil {
					log.Create(
						0,
						"error",
						"kuaishou video-insert",
						item,
						updateError,
					)
				}
			}
		}
	}

	closeError := i.Close()
	if closeError != nil {
		log.Create(0, "error", "kuaishou video close", *list, closeError)
	}

	result.Success = true
	return result
}
