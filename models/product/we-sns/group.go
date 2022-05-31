package modelsWeSns

import (
	"api-kamiya/config"
	"api-kamiya/global"
	"api-kamiya/models"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type Group struct {
	Id           int64
	Name         string
	Introduction string
	PeopleNumber int64
	ProvinceId   int64
	CityId       int64
	IndustryId   int64
	GroupType    string
	GroupAvatar  string
	GroupUrl     string
	UserId       int64
	FromWebName  string
	FromGroupId  int64
	Created      int64
	Updated      int64
	Deleted      int64
	models.Base
}

type GroupList struct {
	Id           int64 `json:"id"`
	Name         string
	Introduction string
	PeopleNumber int64
	ProvinceId   int64
	ProvinceName string
	CityId       int64
	CityName     string
	IndustryName string
	GroupAvatar  string
	GroupType    string
	JoinInId     int64
	Created      int64
	Updated      int64
}

type GroupListData struct {
	IdKey        string `json:"idKey"`
	Name         string `json:"name"`
	Introduction string `json:"introduction"`
	PeopleNumber int64  `json:"peopleNumber"`
	RegionInfo   struct {
		//ProvinceId   int64 `json:"provinceId"`
		ProvinceName string `json:"provinceName"`
		//CityId       int64 `json:"cityId"`
		CityName string `json:"cityName"`
	} `json:"regionInfo"`
	IndustryName string `json:"industryName"`
	GroupAvatar  string `json:"groupAvatar"`
	GroupType    string `json:"groupType"`
	IsJoinIn     bool   `json:"isJoinIn"`
	Created      int64  `json:"created"`
	Updated      int64  `json:"updated"`
}

type GroupDetail struct {
	Name        string `json:"name"`
	GroupUrl    string `json:"groupUrl"`
	GroupAvatar string `json:"groupAvatar"`
}

func init() {
	orm.RegisterModel(new(Group))
}

func (this *Group) GetOrm() orm.Ormer {
	return this.Orm()
}

func (this *Group) TableName() string {
	return models.TableName("we_sns_group")
}

func (this *Group) GetQueryTable() orm.QuerySeter {
	return this.Orm().QueryTable(this)
}

func (this *Group) Detail(userId int64, idKey string, appQuery config.AppClientQuery) global.DataResultModel {
	result := global.GetDataResultModel()
	groupId := global.IdDecrypt(idKey)

	//u := models.User{}
	//userInfo := u.GetInfo( userId )
	//if userInfo.Id <= 0 {
	//	return global.NotLogin()
	//}

	//检查是否开启收费模式
	// 必须登录
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
		dayMaxCount := systemConfig.Get("we-sns.group-day-max-count").(int64)
		count := joinIn.GetCreateCount(userId, "group", 86400)
		if count > dayMaxCount {
			result.Message = "免费用户每天只能查看 " + global.Int64ToString(dayMaxCount) + " 个，升级会员后无限制"
			return result
		}
	}

	group := Group{}
	err := this.Orm().
		QueryTable(this).
		Filter("id", groupId).
		Limit(1).
		One(&group)

	if err != nil {
		fmt.Println("group Detail err:", err)
	}

	groupDetail := GroupDetail{
		Name:        group.Name,
		GroupUrl:    group.GroupUrl,
		GroupAvatar: this.GetGroupAvatar(group.Id),
	}

	// 记录一下，已经被查看
	joinIn.Record(userId, "group", group.Id)

	result.Success = group.Id > 0
	result.Data = groupDetail
	return result
}

func (this *Group) QueryConfig(userId int64) global.DataResultModel {
	result := global.GetDataResultModel()

	queryConfig := config.GroupListQueryConfigData{}
	//region := Region{}
	//queryConfig.Region = region.List()

	industry := Industry{}
	queryConfig.IndustryList = industry.List(userId).Data.([]config.IndustryList)

	result.Data = queryConfig
	result.Success = true
	return result
}

func (this *Group) List(
	userId int64,
	query config.GroupListQuery,
	pagination config.Pagination,
	targetUserId int64,
) global.DataResultModel {
	result := global.GetDataResultModel()
	outData := global.GetFuncListResult(userId)
	// 注意 sql 注入
	query.SearchValue = global.SafeString(query.SearchValue)

	list := []GroupList{}

	qb := this.ListDbWhere(userId, query, targetUserId).
		Limit(int(pagination.PageSize)).
		Offset(int((pagination.Current - 1) * pagination.PageSize))

	qbCount := this.ListDbWhere(userId, query, targetUserId)

	sql := qb.String()
	// 执行 SQL 语句
	_, err := this.Orm().Raw(sql).QueryRows(&list)
	if err != nil {
		fmt.Println("sql group list:", err)
	}

	//fmt.Println("sql group list sql:", sql )

	outData.Pagination = global.SqlPaginationBySql(qbCount, pagination)

	if err != nil {
		fmt.Println("group  list:", err)
	}

	listData := []GroupListData{}
	for _, items := range list {

		itemData := GroupListData{
			IdKey:        global.IdEncrypt(items.Id),
			Name:         items.Name,
			Introduction: items.Introduction,
			PeopleNumber: items.PeopleNumber,
			IndustryName: items.IndustryName,
			GroupAvatar:  this.GetGroupAvatar(items.Id), //items.GroupAvatar,
			GroupType:    items.GroupType,
			IsJoinIn:     items.JoinInId > 0,
			Created:      items.Created,
			Updated:      items.Updated,
		}
		//itemData.RegionInfo.ProvinceId = items.ProvinceId
		itemData.RegionInfo.ProvinceName = items.ProvinceName
		//itemData.RegionInfo.CityId = items.CityId
		itemData.RegionInfo.CityName = items.CityName

		listData = append(listData, itemData)
	}

	outData.List = listData
	result.Success = true
	result.Data = outData
	return result
}

func (this *Group) ListDbWhere(
	userId int64,
	query config.GroupListQuery,
	targetUserId int64,
) orm.QueryBuilder {
	now := time.Now().Unix()
	qb := this.
		GetQB().
		Select(
			"we_sns_group.*",
			"we_sns_industry.name as industry_name",
			"region_province.name as province_name",
			"region_city.name as city_name",

			"we_sns_join_in.id as join_in_id",

		).
		From("we_sns_group").
		LeftJoin("we_sns_industry").
		On("we_sns_industry.id = we_sns_group.industry_id").
		LeftJoin("region_province").
		On("region_province.code = we_sns_group.province_id").
		LeftJoin("region_city").
		On("region_city.code = we_sns_group.city_id").
		LeftJoin("we_sns_join_in").
		On("we_sns_join_in.target_id = we_sns_group.id AND we_sns_join_in.user_id = " + global.Int64ToString(userId) + " AND we_sns_join_in.join_type = \"group\" ").
		Where("we_sns_group.created > " + global.Int64ToString(now-864000*7)).
		And("we_sns_group.deleted = 0 ").
		And("we_sns_group.group_url != '' ")

	if targetUserId > 0 {
		qb = qb.And("we_sns_group.user_id = " + global.Int64ToString(targetUserId))
	}

	if query.ProvinceId > 0 {
		qb = qb.And("we_sns_group.province_id = " + global.Int64ToString(query.ProvinceId))
	}

	if query.CityId > 0 {
		qb = qb.And("we_sns_group.city_id = " + global.Int64ToString(query.CityId))
	}

	if query.IndustryId > 0 {
		qb = qb.And("we_sns_group.industry_id = " + global.Int64ToString(query.IndustryId))
	}

	if query.GroupType == "company" {
		qb = qb.And("we_sns_group.group_type = 'company' ")
	} else if query.GroupType == "wechat" {
		qb = qb.And("we_sns_group.group_type = 'wechat' ").
			And("we_sns_group.people_number < 200")
	} else if query.GroupType == "hundred" {
		qb = qb.And("we_sns_group.group_type = 'wechat'").
			And("we_sns_group.people_number >= 100").
			And("we_sns_group.people_number < 200")
	}

	//关键词
	if query.SearchValue != "" {
		qb = qb.And("we_sns_group.name LIKE '%" + query.SearchValue + "%' ")
	}

	qb = qb.OrderBy("we_sns_group.updated").
		Desc()

	return qb
}

func (this *Group) GetGroupAvatar(groupId int64) string {
	index := global.Int64ToString(groupId % 150)
	return "https://wechat-groups.oss-cn-hangzhou.aliyuncs.com/group-avatar/" + index + ".png"
}

