package modelsWatermark

import (
	configWatermark "api-kamiya/config/watermark"
	"api-kamiya/global"
	"api-kamiya/models"
	"crypto/tls"
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/client/httplib"
	"net/url"

	"github.com/beego/beego/v2/client/orm"
	"time"
)

type Homepage struct {
	Id           int64
	UserId       int64
	Title        string
	ShareUrl     string
	Status       string // ENUM('init', 'fail', 'success')
	ErrorMessage string
	Created      int64
	Updated      int64
	models.Base
}

type HomepageOutData struct {
	HasMore   bool                   `json:"hasMore"`
	Paging    interface{}            `json:"paging"`
	VideoList []HomepageOutDataVideo `json:"videoList"`
}

type HomepageOutDataVideo struct {
	CoverUrl   string   `json:"coverUrl"`
	ImageUrls  []string `json:"imageUrls"`
	Title      string   `json:"title"`
	VideoUrl   string   `json:"videoUrl"`
	VideoType  string   `json:"videoType"` // 'image' | 'video' | '';
	AudioUrl   string   `json:"audioUrl"`
	CreateTime int64    `json:"created"`
}

func init() {
	orm.RegisterModel(new(Homepage))
}

func (this *Homepage) TableName() string {
	return "watermark_homepage"
}

func (this *Homepage) GetQueryTable() orm.QuerySeter {
	return this.Base.GetOrm().QueryTable(this)
}

func (this *Homepage) Create(userId int64, shareUrl string, paging string) global.DataResultModel {
	result := global.GetDataResultModel()

	sqlQuery := this.GetQueryTable()
	// 校验一下
	shareUrl = global.SafeSqlText(shareUrl)
	if shareUrl == "" {
		result.Message = "请输入分享网址"
		return result
	}

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
	monthMaxCount := systemConfig.Get("remove-watermark.homepage.month-max-count").(int64)
	count := this.GetCreateCount(userId, 30*86400)
	if count > monthMaxCount {
		result.Message = "你本月使用已超出配额，请联系管理员微信"
		return result
	}

	i, _ := sqlQuery.PrepareInsert()

	businessLog := models.BusinessLog{}

	insertData := Homepage{
		ShareUrl: shareUrl,
		UserId:   userId,
		Status:   "init",
		Created:  time.Now().Unix(),
	}

	id, err := i.Insert(&insertData)
	if err != nil {
		businessLog.Create(
			userId,
			"error",
			"Homepage Create error",
			map[string]interface{}{
				"id":    id,
				"query": shareUrl,
			},
			err,
		)
	}
	i.Close()

	// 抓取
	resultData := this.GetUrlData(shareUrl, paging)
	//resultData := this.GetUrlDataMock( shareUrl, paging )

	status := ""
	errorMessage := ""
	if resultData.Success {
		status = "success"
		result.Success = true
		apiDataListResult := resultData.Data.(configWatermark.HomepageApiResult)
		if !apiDataListResult.Succ {
			result.Success = false
			if apiDataListResult.Code == 400 {
				result.Message = "主页的网址有问题，请检查"
				result.Code = "url-error"
			}
			status = "fail"
		} else {
			videoList := []HomepageOutDataVideo{}

			for _, dataItem := range apiDataListResult.Video {
				videoUrl := dataItem.PlayUrl
				if videoUrl == "" {
					videoUrl = dataItem.Url
				}
				videoType := "video"
				if videoUrl == "" {
					videoType = "image"
				}
				videoList = append(videoList, HomepageOutDataVideo{
					CoverUrl:   dataItem.Thumb,
					ImageUrls:  []string{},
					Title:      dataItem.Describe,
					VideoUrl:   videoUrl,
					AudioUrl:   "",
					VideoType:  videoType,
					CreateTime: dataItem.CreateTime,
				})
			}
			result.Data = HomepageOutData{
				HasMore:   apiDataListResult.HasMore,
				VideoList: videoList,
				Paging:    apiDataListResult.Pcursor,
			}
			if len(videoList) == 0 {
				status = "list-empty"
			} else {
				status = "success"
			}

		}
	} else {
		status = "fail"
		errorMessage = global.Slice(resultData.Message, 0, 250)
	}

	// 更新数据；
	_, upError := this.GetQueryTable().
		Filter("id", id).
		Limit(1).
		Update(orm.Params{
			"status":       status, // ENUM('init', 'fail', 'success')
			"ErrorMessage": errorMessage,
			"updated":      time.Now().Unix(),
		})

	if upError != nil {
		businessLog.Create(
			userId,
			"error",
			"Homepage upError error",
			map[string]interface{}{
				"id":       id,
				"shareUrl": shareUrl,
			},
			upError,
		)
		fmt.Println("home page Create Update:", status, errorMessage, upError)
	}

	return result
}

// timeSecond: 往前推多长时间；
func (this *Homepage) GetCreateCount(userId int64, timeSecond int64) int64 {
	now := time.Now().Unix()
	count, _ := this.GetQueryTable().
		Filter("userId", userId).
		Filter("status", "success").
		Filter("created__gt", now-timeSecond).
		Count()
	return count
}

func (this *Homepage) GetUrlDataMock(shareUrl string, paging string) global.DataResultModel {
	data := global.Mock("watermark", "homepage", "homepage")
	result := global.GetDataResultModel()

	//反序列化
	query := configWatermark.HomepageApiResult{}
	errJson := json.Unmarshal([]byte(data), &query)
	if errJson != nil {
		fmt.Println("err 12", errJson)
		result.Message = errJson.Error()
		return result
	}

	//if paging != "" {
	//	query.HasMore = false
	//}

	result.Data = query

	result.Success = true
	return result
}

func (this *Homepage) GetUrlData(shareUrl string, paging string) global.DataResultModel {
	result := global.GetDataResultModel()

	host := beego.AppConfig.String("RemoveWatermark.HostHomepage")
	userId := beego.AppConfig.String("RemoveWatermark.UserId")
	secretKey := beego.AppConfig.String("RemoveWatermark.SecretKey")

	postUrl := host +
		"?userId=" + userId +
		"&secretKey=" + secretKey +
		"&url=" + url.QueryEscape(shareUrl) +
		"&paging=" + paging +
		"&count=24"

	req := httplib.Get(postUrl)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	bytesResult, err := req.Bytes()
	if err != nil {
		fmt.Println("get Upload RemoteUrl :", err)
		result.Message = "获取远程文件失败"
		return result
	}

	//fmt.Println("bytesResult 11111", postUrl, string( bytesResult ) )

	//反序列化
	query := configWatermark.HomepageApiResult{}
	errJson := json.Unmarshal(bytesResult, &query)
	if errJson != nil {
		result.Message = errJson.Error()
		return result
	}
	result.Data = query
	result.Success = true

	return result
}
