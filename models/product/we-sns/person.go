package modelsWeSns

import (
	"api-kamiya/config"
	"api-kamiya/global"
	"api-kamiya/models"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type Person struct {
	Id           int64
	Nickname     string
	AvatarUrl    string
	WechatQrUrl  string
	Gender       int
	ProvinceId   int64
	CityId       int64
	Introduction string
	UserId       int64
	FromWebName  string
	FromPersonId int64
	Created      int64
	Deleted      int64
	models.Base
}

type PersonList struct {
	Id           int64
	Nickname     string
	AvatarUrl    string
	Gender       int
	ProvinceName string
	CityName     string
	Introduction string
	JoinInId     int64
}

type PersonListData struct {
	IdKey        string `json:"idKey"`
	Nickname     string `json:"nickname"`
	AvatarUrl    string `json:"avatarUrl"`
	Gender       int    `json:"gender"`
	ProvinceName string `json:"provinceName"`
	CityName     string `json:"cityName"`
	Introduction string `json:"introduction"`
	IsJoinIn     bool   `json:"isJoinIn"`
}

type PersonDetail struct {
	WechatQrUrl string `json:"wechatQrUrl"`
	Nickname    string `json:"nickname"`
	AvatarUrl   string `json:"avatarUrl"`
}

func init() {
	orm.RegisterModel(new(Person))
}

func (this *Person) TableName() string {
	return models.TableName("we_sns_person")
}

func (this *Person) GetQueryTable() orm.QuerySeter {
	return this.Orm().QueryTable(this)
}

func (this *Person) Detail(userId int64, idKey string, appQuery config.AppClientQuery) global.DataResultModel {
	result := global.GetDataResultModel()
	id := global.IdDecrypt(idKey)

	//检查是否开启收费模式
	if userId <= 0 {
		result.Message = "请先登录"
		result.Code = "not-login"
		return result
	}

	joinIn := JoinIn{}

	// 检查会员
	userVip := models.UserVip{}
	userVipDetail := userVip.Detail(userId, "we-sns")
	if !userVipDetail.IsVip {
		// 免费会员，每天可以查看 2 个群；
		systemConfig := global.SystemConfig{}
		dayMaxCount := systemConfig.Get("we-sns.person-day-max-count").(int64)
		count := joinIn.GetCreateCount(userId, "person", 86400)
		if count > dayMaxCount {
			result.Message = "免费用户每天只能查看 " + global.Int64ToString(dayMaxCount) + " 个，升级会员后无限制"
			return result
		}
	}

	person := Person{}
	err := this.Orm().
		QueryTable(this).
		Filter("id", id).
		Limit(1).
		One(&person)

	if err != nil {
		fmt.Println("person Detail err:", err)
	}

	personDetail := PersonDetail{
		WechatQrUrl: person.WechatQrUrl,
		Nickname:    person.Nickname,
		AvatarUrl:   person.AvatarUrl,
	}

	// 记录一下，已经被查看
	joinIn.Record(userId, "person", person.Id)

	result.Success = person.Id > 0
	result.Data = personDetail
	return result
}

func (this *Person) QueryConfig(userId int64) global.DataResultModel {
	result := global.GetDataResultModel()

	queryConfig := config.PersonListQueryConfigData{}
	//region := Region{}
	//queryConfig.Region = region.List()

	industry := Industry{}
	queryConfig.IndustryList = industry.List(userId).Data.([]config.IndustryList)

	result.Data = queryConfig
	result.Success = true
	return result
}

func (this *Person) List(
	userId int64,
	query config.PersonListQuery,
	pagination config.Pagination,
	targetUserId int64,
) global.DataResultModel {
	result := global.GetDataResultModel()
	outData := global.GetFuncListResult(userId)

	// 注意 sql 注入
	query.SearchValue = global.SafeString(query.SearchValue)

	list := []PersonList{}

	qb := this.ListDbWhere(userId, query, targetUserId).
		Limit(int(pagination.PageSize)).
		Offset(int((pagination.Current - 1) * pagination.PageSize))

	qbCount := this.ListDbWhere(userId, query, targetUserId)

	sql := qb.String()
	// 执行 SQL 语句
	_, err := this.Orm().Raw(sql).QueryRows(&list)
	if err != nil {
		fmt.Println("sql per son list:", err)
	}

	//fmt.Println("sql per son list sql:", sql )

	outData.Pagination = global.SqlPaginationBySql(qbCount, pagination)

	if err != nil {
		fmt.Println("per son  list:", err)
	}

	listData := []PersonListData{}
	for _, items := range list {
		itemData := PersonListData{
			IdKey:        global.IdEncrypt(items.Id),
			Nickname:     items.Nickname,
			AvatarUrl:    items.AvatarUrl,
			Gender:       items.Gender,
			ProvinceName: items.ProvinceName,
			CityName:     items.CityName,
			Introduction: items.Introduction,
			IsJoinIn:     items.JoinInId > 0,
		}
		listData = append(listData, itemData)
	}

	outData.List = listData
	result.Success = true
	result.Data = outData
	return result
}

func (this *Person) ListDbWhere(
	userId int64,
	query config.PersonListQuery,
	targetUserId int64,
) orm.QueryBuilder {
	//now := time.Now().Unix()
	qb := this.
		GetQB().
		Select(
			"`we_sns_person`.*",
			"we_sns_industry.name as industry_name",
			"region_province.name as province_name",
			"region_city.name as city_name",

			"we_sns_join_in.id as join_in_id",
		).
		From("`we_sns_person`").
		LeftJoin("we_sns_industry").
		On("we_sns_industry.id = `we_sns_person`.industry_id").
		LeftJoin("region_province").
		On("region_province.code = `we_sns_person`.province_id").
		LeftJoin("region_city").
		On("region_city.code = `we_sns_person`.city_id").
		LeftJoin("we_sns_join_in").
		On("we_sns_join_in.target_id = `we_sns_person`.id AND we_sns_join_in.user_id = " + global.Int64ToString(userId) + " AND we_sns_join_in.join_type = \"we_sns_person\" ").
		Where("`we_sns_person`.deleted = 0 ")

	if targetUserId > 0 {
		qb = qb.And("`we_sns_person`.user_id = " + global.Int64ToString(targetUserId))
	} else {

		if query.ProvinceId > 0 {
			qb = qb.And("`we_sns_person`.province_id = " + global.Int64ToString(query.ProvinceId))
		}

		if query.CityId > 0 {
			qb = qb.And("`we_sns_person`.city_id = " + global.Int64ToString(query.CityId))
		}

		if query.IndustryId > 0 {
			qb = qb.And("we_sns_person.industry_id = " + global.Int64ToString(query.IndustryId))
		}

		if query.Gender == 0 || query.Gender == 1 {
			qb = qb.And("we_sns_person.gender = " + global.IntToString(query.Gender))
		}

		//关键词
		if query.SearchValue != "" {
			qb = qb.And("( `we_sns_person`.nickname LIKE '%" + query.SearchValue + "%' OR `we_sns_person`.introduction LIKE '%\"+ query.SearchValue +\"%' )")
		}
	}

	qb = qb.OrderBy("`we_sns_person`.id").
		Desc()

	return qb
}

func (this *Person) MyList(
	userId int64,
	query config.PersonListQuery,
	pagination config.Pagination,
) global.DataResultModel {
	return this.List(userId, query, pagination, userId)
}
