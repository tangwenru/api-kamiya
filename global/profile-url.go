package global

func MakeUserAvatarUrl( userId int64 ) string {
	key := "1-"+ Int64ToString( userId % 19 )
	avatarUrl := "https://cai-yu-static.oss-cn-beijing.aliyuncs.com/resource/avatar/" + key + ".png"
	return avatarUrl
}
