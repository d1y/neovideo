package repos

import "gorm.io/gorm"

type MacCMSRepo struct {
	gorm.Model
	Api  string `json:"api,omitempty" gorm:"api"`
	Name string `json:"name,omitempty" gorm:"name"`
	R18  bool   `json:"r_18,omitempty" gorm:"r_18"`
}
