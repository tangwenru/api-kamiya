package configCrawler

type CreateKuaiShouVideoListQuery struct {
	Data struct {
		VisionProfilePhotoList struct {
			Feeds []struct {
				Photo struct {
					Id                  string      `json:"id"`
					Duration            int64       `json:"duration"`
					Title               string      `json:"caption"`
					ViewCount           string      `json:"viewCount"`
					LikeCount           int64       `json:"realLikeCount"`
					CoverUrl            string      `json:"coverUrl"`
					PhotoUrl            string      `json:"photoUrl"`
					CoverUrls           interface{} `json:"coverUrls"`
					Timestamp           int64       `json:"timestamp"`
					AnimatedCoverUrl    string      `json:"animatedCoverUrl"`
					Distance            interface{} `json:"distance"`
					VideoRatio          float64     `json:"videoRatio"`
					StereoType          int64       `json:"stereoType"`
					ProfileUserTopPhoto *bool       `json:"profileUserTopPhoto"`
				} `json:"photo"`
			} `json:"feeds"`
		} `json:"visionProfilePhotoList"`
	} `json:"data"`
}
