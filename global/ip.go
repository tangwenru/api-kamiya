package global

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"

)

type ThirdIpResult struct {
	Result struct{
		Ip string `json:"ip"`
	} `json:"result"`
}


func GetIp () string {
	req := httplib.Get("https://apis.map.qq.com/ws/location/v1/ip")
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.Param("key","37KBZ-ECT6X-2DX4P-ZO7SH-JZYNS-VOBFI")
	req.Header("referer", "https://qq.com")

	bytesResult, err := req.Bytes()
	if err != nil {
		fmt.Println("get ip :", err )
	}

	query := ThirdIpResult{}
	errJson := json.Unmarshal( bytesResult, &query )
	if errJson != nil {
		fmt.Println("err 21:", errJson )
		return ""
	}
	return query.Result.Ip
}

