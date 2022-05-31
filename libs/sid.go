package libs

import (
	"api-kamiya/global"
	beego "github.com/beego/beego/v2/adapter"
	"strings"
)

//roleType = "u" | "s" | "c" | "w"; user, staff, client, web
type Sid struct {
	beeController beego.Controller
}

func GetRoleId(roleType string, sid string) int64 {
	text := global.Decrypt(sid)
	textData := strings.Split(text, "_")
	if len(textData) <= 1 || textData[1] == "" {
		return 0
	}

	roleInfo := strings.Split(textData[1], ":")
	// 不符合预期
	if roleInfo[0] != roleType {
		return 0
	}
	return global.String2Int64(roleInfo[1])
}

func MakeSid(roleType string, roleId int64) string {
	// 999_u:123_999 , 随机数，是为了解决 AES 加密部分不变的 "不友好"
	randNum1 := global.RandWord(5, "")
	randNum2 := global.RandWord(5, "")
	text := randNum1 + "_" + roleType + ":" + global.Int64ToString(roleId) + "_" + randNum2
	return global.Encrypt(text)
}
