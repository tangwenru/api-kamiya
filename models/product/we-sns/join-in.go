package modelsWeSns

import (
	"api-kamiya/models"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type JoinIn struct {
	Id       int64
	UserId   int64
	JoinType string
	TargetId int64
	Created  int64
	models.Base
}

func init() {
	orm.RegisterModel(new(JoinIn))
}

func (this *JoinIn) TableName() string {
	return models.TableName("we_sns_join_in")
}

func (this *JoinIn) GetQueryTable() orm.QuerySeter {
	return this.Orm().QueryTable(this)
}

func (this *JoinIn) Record(userId int64, joinType string, targetId int64) {
	//如果 有不理会，如果无插入
	detail := this.Detail(userId, joinType, targetId)
	if detail.Id > 0 {
		return
	}
	// 插入；
	sqlQuery := this.Orm().QueryTable(this)
	i, _ := sqlQuery.PrepareInsert()

	insertData := JoinIn{
		UserId:   userId,
		JoinType: joinType,
		TargetId: targetId,
		Created:  time.Now().Unix(),
	}
	_, insertErr := i.Insert(&insertData)

	i.Close()
	if insertErr != nil {
		fmt.Println("join-in record :", insertErr)
	}
}

func (this *JoinIn) Detail(userId int64, joinType string, targetId int64) JoinIn {
	joinIn := JoinIn{}
	err := this.GetQueryTable().
		Filter("user_id", userId).
		Filter("joinType", joinType).
		Filter("targetId", targetId).
		Limit(1).
		One(&joinIn)

	if err != nil {
		fmt.Println("JoinIn Detail err:", err)
	}

	return joinIn
}

// timeSecond: 往前推多长时间；
func (this *JoinIn) GetCreateCount(userId int64, joinType string, timeSecond int64) int64 {
	now := time.Now().Unix()
	count, _ := this.GetQueryTable().
		Filter("userId", userId).
		Filter("joinType", joinType).
		Filter("created__gt", now-timeSecond).
		Count()
	return count
}
