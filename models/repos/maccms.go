package repos

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type IMacCMS struct {
	Api         string                      `json:"api,omitempty" gorm:"api;unique"`
	Name        string                      `json:"name,omitempty" gorm:"name"`
	R18         bool                        `json:"r_18,omitempty" gorm:"r_18"`
	LastCheck   time.Time                   `json:"last_check,omitempty" gorm:"last_check"`
	Available   bool                        `json:"available,omitempty" gorm:"available"`
	RespType    string                      `json:"resp_type,omitempty" gorm:"resp_type"`
	Category    datatypes.JSONSlice[string] `json:"category,omitempty" gorm:"category"`
	JiexiURL    string                      `json:"jiexi_url,omitempty" gorm:"jiexi_url"`
	JiexiEnable bool                        `json:"jiexi_enable,omitempty" gorm:"jiexi_enable"`
}

type MacCMSRepo struct {
	gorm.Model
	IMacCMS
}
