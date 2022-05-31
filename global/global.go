package global

import (
	"api-kamiya/config"
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	beego "github.com/beego/beego/v2/server/web"
)

type DataResultModel struct {
	Success    bool        `json:"success"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	ServerTime int64       `json:"serverTime"`
	Data       interface{} `json:"data"`
}

func GetDataResultModel() DataResultModel {
	result := DataResultModel{
		Code:       "0",
		Success:    false,
		ServerTime: time.Now().Unix(),
	}
	return result
}

func GetFuncListResult(userId int64) config.FuncListResult {
	result := config.FuncListResult{}
	return result
}

func DataResultSuccess(message string, data interface{}) (r DataResultModel) {
	result := DataResultModel{}

	result.Success = true
	result.Code = "0"
	result.Message = message
	result.Data = data
	result.ServerTime = time.Now().Unix()
	return result
}

func DataResultError(code string, message string, data interface{}) (r DataResultModel) {
	result := DataResultModel{}

	result.Success = code == "0"
	result.Code = code
	result.Message = message
	result.Data = data
	result.ServerTime = time.Now().Unix()
	return result
}

func EncodePassword(password string) string {
	sha := sha1.New()
	sha.Write([]byte("微信封号发麻的很" + password))
	return hex.EncodeToString(sha.Sum([]byte(nil)))
}

//截取字符串
func Slice(word string, startPos int, endPos int) string {
	wordLen := len(word)
	max := Min(int64(wordLen), int64(endPos))

	if startPos > wordLen {
		return ""
	}
	return word[startPos:max]
}

// 四舍五入
func Round(x float64) int64 {
	return int64(math.Floor(x + 0.5))
}

// 四舍五入
func RoundToFloat64(x float64) float64 {
	return float64(int64(math.Floor(x + 0.5)))
}

func String2Int64(str string) int64 {
	if str == "" {
		return 0
	}

	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		id = 0
	}
	return id
}

// string | int64
//func interface2Int64( str interface{}) int64 {
//	if str == "" {
//		return 0
//	}
//
//	id, err := strconv.ParseInt( str, 10, 64)
//	if err != nil {
//		id = 0
//	}
//	return id
//}

func String2Int8(str string) int8 {
	if str == "" {
		return 0
	}

	id, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		id = 0
	}
	return int8(id)
}

func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'g', -1, 64)
}

func Int64ToString(int int64) string {
	out := strconv.FormatInt(int, 10)
	return out
}
func IntToString(int int) string {
	out := strconv.Itoa(int)
	return out
}

func String2Float64(str string) float64 {
	id, err := strconv.ParseFloat(str, 64)
	if err != nil {
		id = 0
	}
	return id
}

func String13ToTime10(str string) int64 {
	return String2Int64(string([]rune(str)[:10]))
}

func Strings2Int64(str string) []int64 {
	tradeIdsStringArray := strings.Split(str, ",")
	var ids []int64
	for _, idString := range tradeIdsStringArray {
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			id = 0
		}
		ids = append(ids, id)
	}
	return ids
}

func String2Array(str string, split string) []string {
	if split == "" {
		split = ","
	}

	if str == "" {
		return []string{}
	}

	return strings.Split(str, split)
}

// 100-200 转换 数字区间
func String2NumberArea(str string, maxNumber int64) (int64, int64) {
	arrayData := String2Array(str, "-")
	// 全部；
	if len(arrayData) != 2 {
		return 0, maxNumber
	}

	min := String2Int64(arrayData[0])
	max := String2Int64(arrayData[1])
	if max == 0 {
		max = maxNumber
	}

	return min, max
}

func Cent2Yen(money int64) float64 {
	x := decimal.NewFromFloat(float64(money))
	var z = x.Div(decimal.NewFromInt32(100))
	moneyYen, _ := strconv.ParseFloat(z.StringFixed(2), 64)
	return moneyYen
}

func Yen2Cent(money float64) int64 {
	return int64(money * 100)
}

// float64 四目运算；
//加
func Add(a float64, b float64) float64 {
	x := decimal.NewFromFloat(a)
	y := decimal.NewFromFloat(b)

	sum := x.Add(y)

	result, err := strconv.ParseFloat(sum.StringFixed(2), 64)
	if err != nil {
		result = 0
	}
	return result
}

//减
func Sub(a float64, b float64) float64 {
	x := decimal.NewFromFloat(a)
	y := decimal.NewFromFloat(b)

	sum := x.Sub(y)

	result, err := strconv.ParseFloat(sum.StringFixed(2), 64)
	if err != nil {
		result = 0
	}
	return result
}

func SqlPagination(sqlQuery orm.QuerySeter, pagination config.Pagination) config.Pagination {
	count, _ := sqlQuery.Count()
	pagination.Total = count

	return pagination
}

// float 64 小数位
func ToFixed(num float64, round int32) float64 {
	result, _ := decimal.NewFromFloat(num).Round(round).Float64()
	return result
}

// float64 -> int64
func ParseInt64(num float64) int64 {
	result := decimal.NewFromFloat(num).IntPart()
	return result
}

func SqlPaginationBySql(
	qb orm.QueryBuilder,
	pagination config.Pagination,
) config.Pagination {

	//  SELECT video.*, video_hot.crawler_digg_rank  FROM video_hot LEFT JOIN video ON video_hot.video_id = video.video_id WHERE video_hot.cycle = "H48"

	//通过函数进行替换
	re3, _ := regexp.Compile("^SELECT[^FROM]+");
	sql := re3.ReplaceAllString(qb.String(), "SELECT COUNT(*) as `count` ")
	//fmt.Println("rep2rep2rep2", sql );

	var detail struct {
		Count int64 `json:"count"`
	}
	err := orm.NewOrm().Raw(sql).QueryRow(&detail)

	if err != nil {
		//log.Create( userId, "error", "SqlPaginationBySql", qb.String(), err )
	}
	pagination.Total = detail.Count
	return pagination
}

func SqlPaginationByRaw(orm orm.Ormer, sql string, pagination config.Pagination) config.Pagination {

	index := strings.Index(sql, " FROM ")
	content := sql[index:len(sql)]

	//把匹配的所有字符a替换成b
	sql = "SELECT COUNT(*) as count " + content

	// ORDER BY `topic`.id DESC LIMIT 20 OFFSET 20
	index = strings.Index(sql, " ORDER BY ")
	if index == -1 {
		index = len(sql)
	}
	sql = sql[0:index]

	var list []struct {
		Count int64
	}
	_, err := orm.Raw(sql).QueryRows(&list)
	if err != nil {
		fmt.Println("SqlPaginationByRaw Error:", sql)
	}
	fmt.Println("SqlPaginationByRaw sql:", sql)

	if len(list) == 0 {
		pagination.Total = 0
		return pagination
	}

	count := list[0].Count
	pagination.Total = count
	return pagination
}

//查找字符是否在数组中
func InArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

func GetStaffRoleName(role string) string {
	RoleNames := map[string]string{
		"admin":   "管理员",
		"operate": "操作员",
		"funds":   "财务",
		"other":   "其他",
	}

	return RoleNames[role]
}

// id 的 大于 0；
func IdEncrypt(id int64) string {
	if id <= 0 {
		return ""
	}
	aesKeyString, _ := beego.AppConfig.String("AesKey")
	var aesKey = []byte(aesKeyString)
	result, err := AesEncrypt([]byte(Int64ToString(id)), aesKey)
	if err != nil {
		return ""
	}
	return result
}

func IdDecrypt(idKey string) int64 {
	aesKeyString, _ := beego.AppConfig.String("AesKey")
	var aesKey = []byte(aesKeyString)
	result, err := AesDecrypt(idKey, aesKey)
	if err != nil {
		return 0
	}
	return String2Int64(string(result))
}

func Encrypt(text string) string {
	if text == "" {
		return ""
	}
	aesKeyString, _ := beego.AppConfig.String("AesKey")
	var aesKey = []byte(aesKeyString)
	result, err := AesEncrypt([]byte(text), aesKey)
	if err != nil {
		return ""
	}
	return result
}

func Decrypt(idKey string) string {
	aesKeyString, _ := beego.AppConfig.String("AesKey")
	var aesKey = []byte(aesKeyString)
	result, err := AesDecrypt(idKey, aesKey)
	if err != nil {
		return ""
	}
	return string(result)
}

func SafeString(str string) string {
	regCompile := regexp.MustCompile("['\\\\`\"~!@#$%^&*()-+=;.,?{}|]+")

	//把匹配的所有字符a替换成b
	return regCompile.ReplaceAllString(str, "")
}

func SafeSqlText(str string) string {
	out := strings.Replace(str, "'", " ", -1)
	out = strings.Replace(out, ";", " ", -1)
	out = strings.Replace(out, "\"", " ", -1)
	out = strings.Replace(out, "(", " ", -1)
	out = strings.Replace(out, ")", " ", -1)
	return out
}

func MobileSafeShow(mobile string) string {
	if mobile == "" {
		return ""
	}
	if len(mobile) < 3 {
		return mobile
	}
	mobileInfo := strings.Split(mobile, "")
	return mobileInfo[0] + mobileInfo[1] + mobileInfo[2] + "******" + mobileInfo[len(mobileInfo)-2] + mobileInfo[len(mobileInfo)-1]
}

func RandWord(len int, words string) string {
	maxLen := len
	var container string
	var strDict = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	if words != "" {
		strDict = words
	}
	b := bytes.NewBufferString(strDict)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < maxLen; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(strDict[randomInt.Int64()])
	}
	return container
}

// 10进制 到任意进制的转换
func ConvertToAnyBase2to36(num int64, base int64) string {
	var baseLen = base
	var sourceString = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "J", "I", "J", "K",
		"L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	code := "";
	for {
		if num <= 0 {
			break
		}
		mod := num % baseLen;
		num = (num - mod) / baseLen;
		code = sourceString[mod] + code;
	}
	return code
}

//整数
func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func ReadTxtFile(path string) string {
	b, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	str := string(b) // convert content to a 'string'
	return str
}