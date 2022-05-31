package models

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type BusinessLog struct {
	Id      int64  `json:"id"`
	UserId  int64  `json:"userId"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Error   string `json:"error"`
	Created int64  `json:"created"`
	Base
}

func init() {
	orm.RegisterModel(new(BusinessLog))
}

//func (this *BusinessLog) orm() (o orm.Ormer) {
//	return orm.NewOrm()
//}

func (this *BusinessLog) TableName() string {
	return "log"
}

func (this *BusinessLog) Create(userId int64, logType string, title string, content interface{}, err interface{}) {
	sqlQuery := this.orm().QueryTable(this)
	i, _ := sqlQuery.PrepareInsert()

	contentByte, _ := json.Marshal(content)
	errByte, _ := json.Marshal(err)

	insertData := BusinessLog{
		UserId:  userId,
		Type:    logType,
		Title:   title,
		Content: string(contentByte),
		Error:   string(errByte),
		Created: time.Now().Unix(),
	}
	_, insertError := i.Insert(&insertData)
	i.Close()
	if insertError != nil {
		// todo log
	}
}
