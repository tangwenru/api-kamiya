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
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type WrongWord struct {
	Id           int64
	UserId       int64
	Text         string
	Status       string // ENUM('init', 'fail', 'success')
	ErrorMessage string
	Created      int64
	Updated      int64
	models.Base
}

func init() {
	//orm.RegisterModel(new(WrongWord))
}

func (this *WrongWord) TableName() string {
	return "watermark_sensitive_word_detect"
}

func (this *WrongWord) GetQueryTable() orm.QuerySeter {
	return this.Base.GetOrm().QueryTable(this)
}

func (this *WrongWord) Create(
	userId int64,
	query *configWatermark.WrongWordCreateQuery,
) global.DataResultModel {
	result := global.GetDataResultModel()

	//sqlQuery := this.GetQueryTable()
	// 校验一下
	text := global.SafeSqlText(query.Text)
	if text == "" {
		result.Message = "请输入太短了"
		return result
	}

	// 必须登录
	//if userId <= 0 {
	//	result.Message = "请先登录"
	//	result.Code = "not-login"
	//	return result
	//}
	//
	//// 检查会员
	//userVip := models.UserVip{}
	//userVipDetail := userVip.Detail(userId, "remove-watermark")
	//if !userVipDetail.IsVip {
	//	result.Code = "need-vip"
	//	result.Message = "此功能为会员专享，请购买会员"
	//	return result
	//}

	//初期简单粗暴一些
	//systemConfig := global.SystemConfig{}
	//monthMaxCount := systemConfig.Get("remove-watermark.sensitive-word-detect.month-max-count").(int64)
	//count := this.GetCreateCount(userId, 30*86400)
	//if count > monthMaxCount {
	//	result.Message = "你本月使用已超出配额，请联系管理员微信"
	//	return result
	//}

	//i, _ := sqlQuery.PrepareInsert()
	//
	//businessLog := models.BusinessLog{}
	//
	//insertData := WrongWord{
	//	Text:    text,
	//	UserId:  userId,
	//	Status:  "init",
	//	Created: time.Now().Unix(),
	//}
	//
	//id, err := i.Insert(&insertData)
	//if err != nil {
	//	businessLog.Create(
	//		userId,
	//		"error",
	//		"WrongWord Create error",
	//		map[string]interface{}{
	//			"id":    id,
	//			"query": query,
	//		},
	//		err,
	//	)
	//}
	//i.Close()

	// 抓取
	resultData := this.GetUrlData(userId, query.Text)
	//resultData := this.GetUrlDataMock( shareUrl, paging )
	result.Message = resultData.Message

	//fmt.Println("resultData：", resultData )

	//status := ""
	//errorMessage := ""
	if resultData.Success {
		//status = "success"
		result.Success = true
		//apiDataListResult := resultData.Data.( []configWatermark.WrongWordApiResult )
		result.Data = resultData.Data
	} else {
		//status = "fail"
		//errorMessage = global.Slice(resultData.Message, 0, 250)
	}

	// 更新数据；
	//this.GetQueryTable().
	//	Filter("id", id).
	//	Limit(1).
	//	Update(orm.Params{
	//		"Status":       status, // ENUM('init', 'fail', 'success')
	//		"ErrorMessage": errorMessage,
	//		"Updated":      time.Now().Unix(),
	//	})
	return result
}

// timeSecond: 往前推多长时间；
func (this *WrongWord) GetCreateCount(userId int64, timeSecond int64) int64 {
	now := time.Now().Unix()
	count, _ := this.GetQueryTable().
		Filter("userId", userId).
		Filter("status", "success").
		Filter("created__gt", now-timeSecond).
		Count()
	return count
}

func (this *WrongWord) GetUrlData(userId int64, text string) global.DataResultModel {
	result := global.GetDataResultModel()

	host := "https://aip.baidubce.com/rpc/2.0/nlp/v1/ecnet"

	accessTokenResult := this.GetBaiduAccessToken(userId)
	if !accessTokenResult.Success {
		return accessTokenResult
	}

	accessToken := accessTokenResult.Data.(configWatermark.WrongWordAccessToken)

	postUrl := host +
		"?access_token=" + accessToken.AccessToken +
		"&charset=UTF-8"

	req := httplib.Post(postUrl)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	params := map[string]string{
		"text": text,
	}
	sendBody, err := json.Marshal(params)
	req.Body(string(sendBody))
	req.Header("Content-Type", "application/json")

	bytesResult, err := req.Bytes()
	if err != nil {
		fmt.Println("百度错别字:", err)
		result.Message = "获取远程文件失败"
		return result
	}

	fmt.Println("bytesResult 222", postUrl, string(bytesResult))

	//反序列化
	query := configWatermark.WrongWordApiResult{}
	errJson := json.Unmarshal(bytesResult, &query)
	if errJson != nil {
		result.Message = errJson.Error() // errJson.Error()
		return result
	}
	result.Data = query
	result.Success = true

	return result
}

func (this *WrongWord) GetBaiduAccessToken(userId int64) global.DataResultModel {
	result := global.GetDataResultModel()
	apiKey := beego.AppConfig.String("WrongWord.ApiKey")
	secretKey := beego.AppConfig.String("WrongWord.SecretKey")

	log := models.BusinessLog{}

	postUrl := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + apiKey + "&client_secret=" + secretKey

	req := httplib.Post(postUrl)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	bytesResult, err := req.Bytes()
	if err != nil {
		fmt.Println("get GetBaiduAccessToken RemoteUrl :", err)
		result.Message = "获取远程文件失败"
		log.Create(userId, "error", "GetBaiduAccessToken", "", err)

		return result
	}

	//fmt.Println("bytesResult 11111", postUrl, string( bytesResult ) )
	query := configWatermark.WrongWordAccessToken{}
	errJson := json.Unmarshal(bytesResult, &query)
	if errJson != nil {
		result.Message = errJson.Error()
		fmt.Println("GetBaiduAccessToken errJson", errJson)
		return result
	}
	//fmt.Println("GetBaiduAccessToken out", string( bytesResult ) )

	result.Data = query
	result.Success = true
	return result
}
