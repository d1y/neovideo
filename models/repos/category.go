package repos

import (
	"d1y.io/neovideo/models"
	"gorm.io/datatypes"
)

type IVideoCategory struct {
	Name    string                    `json:"name" gorm:"name,unique"`
	R18     bool                      `json:"r18" gorm:"r18"`
	Sources datatypes.JSONSlice[uint] `json:"sources" gorm:"sources"` // 来源(maccms) || 或许不需要这种关联操作了!
}

type VideoCategoryRepo struct {
	models.BaseRepo
	IVideoCategory
}

func (vid *VideoCategoryRepo) TableName() string {
	return "t_category"
}
