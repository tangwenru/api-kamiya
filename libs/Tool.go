package libs

import (
	beego "github.com/beego/beego/v2/adapter"
)

type Tool struct {
	beeController beego.Controller
}

func main() {

}

func (this *Tool) QueryUrl(param string) string {

	con := beego.Controller{}

	return con.GetString(":" + param)
}

//func ( this *Tool ) QueryUrlInt32 ( param string ) int32 {
//	paramString := this.beeController.Get   ;//   .GetString( ":"+ param );
//	paramInt64, _ := strconv.ParseInt( paramString, 10, 32 );
//	paramInt32 := int32( paramInt64 );
//	return paramInt32;
//}
