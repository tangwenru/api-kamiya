package modelsWeSns

import (
	"api-kamiya/global"
	"api-kamiya/models"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type Crawler struct {
	models.Base
}

func init() {

}

type CrawlerResult struct {
	Success bool `json:"success"`
	Data    struct {
		GroupList  []CrawlerGroupList
		PersonList []CrawlerPersonList
	} `json:"data"`
}

type CrawlerGroupList struct {
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
}

type CrawlerPersonList struct {
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
}

func (this *Crawler) Task() global.DataResultModel {
	result := global.GetDataResultModel()
	url := models.AppConfig.String("WeSns.CrawlerUrl")
	req := httplib.Get(url)

	bytesResult, err := req.Bytes()
	if err != nil {
		fmt.Println("Task 0:", err)

	}

	query := CrawlerResult{}
	errJson := json.Unmarshal(bytesResult, &query)
	if errJson != nil {
		fmt.Println("Task 2:", errJson)
		fmt.Println("Task 3:", string(bytesResult))

	}

	fmt.Println("detailGroup --------------", len(query.Data.GroupList))

	now := time.Now().Unix()

	// Group 插入数据库
	// 开始记录
	group := Group{}
	sqlQuery := group.GetQueryTable()
	i, _ := sqlQuery.PrepareInsert()
	for _, items := range query.Data.GroupList {
		// 检查
		detailGroup := Group{}
		detailErr := group.GetQueryTable().
			Filter("from_web_name", items.FromWebName).
			Filter("from_group_id", items.FromGroupId).
			Limit(1).
			One(&detailGroup)

		if detailErr != nil {
			//fmt.Println("detailErr 1", detailErr )
		}

		// 更新
		if detailGroup.Id > 0 {
			group.GetQueryTable().
				Filter("id", detailGroup.Id).
				Limit(1).
				Update(orm.Params{
					"name":          items.Name,
					"people_number": items.PeopleNumber,
					"updated":       items.Updated,
				})
		} else {
			insertData := Group{
				Name:         items.Name,
				Introduction: items.Introduction,
				PeopleNumber: items.PeopleNumber,
				ProvinceId:   items.ProvinceId,
				CityId:       items.CityId,
				IndustryId:   items.IndustryId,
				GroupType:    items.GroupType,
				GroupAvatar:  items.GroupAvatar,
				GroupUrl:     items.GroupUrl,
				UserId:       items.UserId,
				FromWebName:  items.FromWebName,
				FromGroupId:  items.FromGroupId,
				Created:      items.Created,
				Updated:      items.Updated,
				Deleted:      items.Deleted,
			}

			_, insertError := i.Insert(&insertData)
			if insertError != nil {
				fmt.Println("crawler order create:", insertError)
			}
		}
	}
	defer i.Close()

	// Person 插入数据库
	// 开始记录
	person := Person{}
	iPerson, _ := person.GetQueryTable().PrepareInsert()
	for _, items := range query.Data.PersonList {
		// 检查
		detailPerson := Person{}
		detailErr := person.GetQueryTable().
			Filter("from_web_name", items.FromWebName).
			Filter("from_person_id", items.FromPersonId).
			Limit(1).
			One(&detailPerson)

		if detailErr != nil {
			//fmt.Println("detailErr 1", detailErr )
		}

		// 更新
		if detailPerson.Id > 0 {
			//person.GetQueryTable().
			//	Filter("id", detailGroup.Id ).
			//	Limit(1).
			//	Update( orm.Params{
			//		"name": items.Name,
			//		"people_number": items.PeopleNumber,
			//		"updated": items.Updated,
			//	})
		} else {
			//插入
			insertData := Person{
				FromWebName:  items.FromWebName,
				FromPersonId: items.Id,
				Nickname:     items.Nickname,
				AvatarUrl:    items.AvatarUrl,
				Gender:       items.Gender,
				CityId:       0,
				Introduction: items.Introduction,
				Created:      now,
			}
			_, insertError := iPerson.Insert(&insertData)
			if insertError != nil {
				fmt.Println("person create", insertError)
			}
		}
	}
	defer iPerson.Close()

	result.Success = true
	result.Data = query
	return result
}