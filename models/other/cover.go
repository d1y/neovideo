package other

import (
	"d1y.io/neovideo/models"
	"d1y.io/neovideo/models/repos"
)

// 需要下载图片(封面)的任务表
//
// 一般在入库的时候就会下载的, 但是下载的时候失败了(maybe)
//
// 需要把任务放到这里来, 在某个时刻下载

type IImageCover struct {
	URL   string `json:"url"`
	Video repos.VideoRepo
}

type ImageCoverDownload struct {
	models.BaseRepo
	IImageCover
}

func (icd *ImageCoverDownload) TableName() string {
	return "t_cover_download"
}
