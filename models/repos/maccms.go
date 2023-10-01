package repos

import (
	"time"

	"d1y.io/neovideo/models"
	"gorm.io/datatypes"
)

type IMacCMS struct {
	Api         string                      `json:"api" gorm:"api;unique"`
	Name        string                      `json:"name" gorm:"name"`
	R18         bool                        `json:"r_18" gorm:"r_18"`
	LastCheck   time.Time                   `json:"last_check" gorm:"last_check"`
	Available   bool                        `json:"available" gorm:"available"`
	RespType    string                      `json:"resp_type" gorm:"resp_type"`
	Category    datatypes.JSONSlice[string] `json:"category" gorm:"category"`
	JiexiURL    string                      `json:"jiexi_url" gorm:"jiexi_url"`
	JiexiEnable bool                        `json:"jiexi_enable" gorm:"jiexi_enable"`
}

type MacCMSRepo struct {
	models.BaseRepo
	IMacCMS
}

func (m *MacCMSRepo) TableName() string {
	return "t_maccms"
}
