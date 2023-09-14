package repos

import (
	"time"

	"gorm.io/gorm"
)

type MacCMSRepo struct {
	gorm.Model
	Api       string    `json:"api,omitempty" gorm:"api"`
	Name      string    `json:"name,omitempty" gorm:"name"`
	R18       bool      `json:"r_18,omitempty" gorm:"r_18"`
	LastCheck time.Time `json:"last_check,omitempty" gorm:"last_check"`
	Available bool      `json:"available,omitempty" gorm:"available"`
	RespType  int       `json:"resp_type,omitempty" gorm:"resp_type"`
}
