package repos

import (
	"time"

	"d1y.io/neovideo/models"
	"gorm.io/datatypes"
)

type IVideoData struct {
	URL   string `json:"url" gorm:"url"`
	Name  string `json:"name" gorm:"name"`
	Embed bool   `json:"embed"`
}

type IVideoDataInfo struct {
	Flag   string       `json:"flag" gorm:"flag"`
	Videos []IVideoData `json:"videos" gorm:"videos"`
}

type IVideo struct {
	SpiderType string                              `json:"spider_type" gorm:"spider_type"` // 爬虫类型(TODO: 1=> maccms)
	Title      string                              `json:"title" gorm:"title"`             // 标题
	Desc       string                              `json:"desc" gorm:"desc"`               // 描述
	Mid        uint                                `json:"mid" gorm:"mid"`                 // [maccms] id(方便关联起来)
	RealType   string                              `json:"real_type" gorm:"real_type"`     // 分类名称
	RealID     int                                 `json:"real_id" gorm:"real_id"`         // 真实的id
	RealTime   time.Time                           `json:"real_time" gorm:"real_time"`     // 真实的更新(创建)时间
	RealCover  string                              `json:"real_cover" gorm:"real_cover"`   // 真实的封面
	Cover      string                              `json:"cover" gorm:"cover"`             // 贩卖
	CategoryID int                                 `json:"category_id" gorm:"category_id"` // 分类id
	Videos     datatypes.JSONSlice[IVideoDataInfo] `json:"videos" gorm:"videos"`           // 视频
	Lang       string                              `json:"lang" gorm:"lang"`               // 语言
	Area       string                              `json:"area" gorm:"area"`               // 国家
	Year       string                              `json:"year" gorm:"year"`               // 年份
	State      string                              `json:"state" gorm:"state"`             // 状态
	Actor      string                              `json:"actor" gorm:"actor"`             // 演员
	Director   string                              `json:"director" gorm:"director"`       // 导演
	R18        bool                                `json:"r18" gorm:"r18"`                 // 18+
}

type VideoRepo struct {
	models.BaseRepo
	IVideo
}

func (vr *VideoRepo) TableName() string {
	return "t_videos"
}
