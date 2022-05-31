package global

import (
	"api-kamiya/config"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"strconv"
	"strings"
)

type SystemConfig struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	ValueType   string `json:"valueType"`
	ValueConfig string `json:"valueConfig"`
	Value       string `json:"value"`
}

func init() {
	orm.RegisterModel(new(SystemConfig))
}

func (this *SystemConfig) GetQueryTable() orm.QuerySeter {
	return this.orm().QueryTable(this)
}

func (m *SystemConfig) TableName() string {
	return "system_config"
}

func (b *SystemConfig) orm() (o orm.Ormer) {
	return orm.NewOrm()
}

//生成一个 sid
func (this *SystemConfig) Get(name string) interface{} {
	systemConfig := SystemConfig{}

	sqlQuery := this.
		orm().
		QueryTable(this).
		Filter("name", name)

	err :=
		sqlQuery.
			Limit(1).
			One(&systemConfig)

	// todo err
	if err != nil {
		fmt.Println("SystemConfig get err:", err)
	}
	//fmt.Println( "SystemConfig get:", name,  systemConfig )

	if systemConfig.Id <= 0 {
		return ""
	}

	return this.FormatValue(systemConfig.ValueConfig, systemConfig.ValueType, systemConfig.Value)
}

func (this *SystemConfig) FormatValue(valueConfig string, valueType string, value string) interface{} {
	//ENUM('整数', '小数', '文本', '对象数组', '网址', '下拉框', '是否', '多选', '时间区间', '时间', '')
	var data interface{}
	switch valueType {
	case "整数":
		d, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			data = 0
		} else {
			data = int64(d)
		}

	case "小数":
		d, err := strconv.ParseFloat(value, 64)
		if err != nil {
			data = 0
		} else {
			data = d
		}
	case "多选":
		dataInfo := strings.Split(value, ",")

		// []
		valueConfigResult := &[]config.SystemConfigKeyValue{}
		errData := json.Unmarshal([]byte(valueConfig), valueConfigResult)
		if errData != nil {
			valueConfigResult = &[]config.SystemConfigKeyValue{}
		}

		dataArray := []config.SystemConfigKeyValue{}
		for _, items := range *valueConfigResult {
			if InArray(items.Key, dataInfo) {
				dataArray = append(dataArray, items)
			}
		}
		data = dataArray
	case "是否":
		data = value == "是"

	case "网址", "文本", "下拉框", "时间":
		data = value

	case "时间区间":
		data = strings.Split(value, "-")
	default:
		data = value

	}
	return data
}
