package other

import (
	"d1y.io/neovideo/models"
)

// 需要下载图片(封面)的任务表
//
// 一般在入库的时候就会下载的, 但是下载的时候失败了(maybe)
//
// 需要把任务放到这里来, 在某个时刻下载

type IImageCover struct {
	URL      string `json:"url" gorm:"url"`
	Filename string `json:"filename" gorm:"filename"`
	Reason   string `json:"reason" gorm:"reason"`
}

type ImageCoverTask struct {
	models.BaseRepo
	IImageCover
}

func (icd *ImageCoverTask) TableName() string {
	return "t_cover_task"
}

func NewCoverTask(img string, filename string, err error) *ImageCoverTask {
	return &ImageCoverTask{
		IImageCover: IImageCover{
			URL:      img,
			Filename: filename,
			Reason:   err.Error(),
		},
	}
}
