package global

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func Mock( productType string, subType string, apiName string ) string {
	data, err := GetFileContentAsStringLines( "./mock/"+ productType +"/"+ subType +"/"+ apiName +".json"  )
	if err != nil{
		return ""
	}
	return data
}


func GetFileContentAsStringLines(filePath string) ( string, error) {
	fmt.Println("get file content as lines: %v", filePath)
	result := []string{}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("read file: %v error: %v", filePath, err)
		return "", err
	}
	s := string(b)
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		result = append(result, lineStr)
	}
	fmt.Println("get file content as lines: %v, size: %v", filePath, len(result))
	return strings.Join( result, ""), nil
}