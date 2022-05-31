package modelsWatermark

import (
	configWatermark "api-kamiya/config/watermark"
	"api-kamiya/global"
	"api-kamiya/models"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type AudioText struct {
	Id           int64
	UserId       int64
	UseOurApi    bool
	FileId       string
	FileName     string
	Status       string // ENUM('init', 'fail', 'success')
	VideoSize    int64
	Duration     int64
	ErrorMessage string
	Created      int64
	models.Base
}

func init() {
	orm.RegisterModel(new(AudioText))
}

func (this *AudioText) TableName() string {
	return "watermark_audio_text"
}

func (this *AudioText) GetQueryTable() orm.QuerySeter {
	return this.Base.GetOrm().QueryTable(this)
}

func (this *AudioText) Create(userId int64, query *configWatermark.AudioTextCreateQuery) global.DataResultModel {
	result := global.GetDataResultModel()

	sqlQuery := this.GetQueryTable()

	checkCanUseInfo := this.CheckCanUse(userId)
	if !checkCanUseInfo.Success {
		return checkCanUseInfo
	}

	i, _ := sqlQuery.PrepareInsert()

	businessLog := models.BusinessLog{}

	for _, items := range query.FileList {
		// 校验；
		items.FileId = global.SafeSqlText(global.Slice(items.FileId, 0, 20))
		items.Status = global.SafeSqlText(global.Slice(items.Status, 0, 20))
		items.FileName = global.SafeSqlText(global.Slice(items.FileName, 0, 200))
		items.ErrorMessage = global.SafeSqlText(global.Slice(items.ErrorMessage, 0, 240))

		detail := this.Detail(items.FileId)
		if detail.Success {
			// 更新数据；
			upData := orm.Params{
				"FileName":  items.FileName,
				"Status":    items.Status, // ENUM('init', 'fail', 'success')
				"VideoSize": items.VideoSize,
				"Duration":  items.Duration,
				"UseOurApi": items.UseOurApi,
			}
			if items.ErrorMessage != "" {
				upData["error_message"] = items.ErrorMessage
			}
			_, upError := this.GetQueryTable().
				Filter("FileId", items.FileId).
				Limit(1).
				Update(upData)

			if upError != nil {

			}
		} else {
			insertData := AudioText{
				UserId:       userId,
				FileId:       items.FileId,
				UseOurApi:    items.UseOurApi,
				FileName:     items.FileName,
				Status:       items.Status, // ENUM('init', 'fail', 'success')
				VideoSize:    items.VideoSize,
				Duration:     items.Duration,
				ErrorMessage: items.ErrorMessage,
				Created:      time.Now().Unix(),
			}

			id, err := i.Insert(&insertData)
			if err != nil {
				businessLog.Create(
					userId,
					"error",
					"AudioText Create error",
					map[string]interface{}{
						"id":    id,
						"query": *query,
					},
					err,
				)
			}
		}
	}
	i.Close()

	result.Success = true
	return result
}

// timeSecond: 往前推多长时间； , 使用平台的；
func (this *AudioText) GetCreateCount(userId int64, timeSecond int64) int64 {
	now := time.Now().Unix()
	count, _ := this.GetQueryTable().
		Filter("userId", userId).
		Filter("status", "success").
		Filter("use_our_api", true).
		Filter("created__gt", now-timeSecond).
		Count()
	return count
}

func (this *AudioText) Detail(fileId string) global.DataResultModel {
	result := global.GetDataResultModel()
	detail := AudioText{}
	err := this.GetQueryTable().
		Filter("fileId", fileId).
		Limit(1).
		One(&detail)

	if err != nil {

	}
	result.Data = detail
	result.Success = detail.Id > 0
	return result
}

func (this *AudioText) CheckCanUse(userId int64) global.DataResultModel {
	// 看看是不是会员，看看有没有超过 数量；
	result := global.GetDataResultModel()

	// 必须登录
	if userId <= 0 {
		result.Message = "请先登录"
		result.Code = "not-login"
		return result
	}

	// 检查会员
	userVip := models.UserVip{}
	userVipDetail := userVip.Detail(userId, "remove-watermark")
	if !userVipDetail.IsVip {
		result.Code = "need-vip"
		result.Message = "此功能为会员专享，请购买会员"
		return result
	}

	// 简直 会员级别对应的视频数量；
	//初期简单粗暴一些
	systemConfig := global.SystemConfig{}
	monthMaxCount := systemConfig.Get("remove-watermark.audio_text.month-max-count").(int64)
	count := this.GetCreateCount(userId, 30*86400)
	if count > monthMaxCount {
		result.Message = "你本月使用已超出配额，请配置自己的阿里云参数"
		return result
	}

	result.Success = true
	return result
}