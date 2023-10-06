package other

import "d1y.io/neovideo/models"

type ISpiderTask struct {
	SpiderType string `json:"spider_type" gorm:"spider_type"` // 爬虫类型(maccms)
	Mid        uint   `json:"mid" gorm:"mid"`                 // 爬虫(maccms)的真实的id
	Page       int    `json:"page" gorm:"page"`               // 爬取页数
	Successful bool   `json:"successful" gorm:"successful"`   // 是否成功
	Message    string `json:"message" gorm:"message"`         // 消息
}

type SpiderTask struct {
	models.BaseRepo
	ISpiderTask
}

func NewSpiderTask(mid uint, page int) *SpiderTask {
	return &SpiderTask{
		ISpiderTask: ISpiderTask{
			SpiderType: "maccms",
			Mid:        mid,
			Page:       page,
		},
	}
}

func (st *SpiderTask) TableName() string {
	return "t_spider_task"
}

func (st *SpiderTask) SetPage(pg int) *SpiderTask {
	st.Page = pg
	return st
}

func (st *SpiderTask) SetSuccessful(msg string) {
	st.Successful = true
	st.Message = msg
}

func (st *SpiderTask) SetFailed(reason string) {
	st.Successful = false
	st.Message = reason
}
