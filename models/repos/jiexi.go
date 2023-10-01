package repos

import "gorm.io/gorm"

type IJiexi struct {
	Name string `json:"name,omitempty" gorm:"name" form:"name"`
	URL  string `json:"url,omitempty" gorm:"url;unique" form:"url"`
	Note string `json:"note,omitempty" gorm:"note" form:"note"`
}

type JiexiRepo struct {
	gorm.Model
	IJiexi
}

func (j *JiexiRepo) TableName() string {
	return "t_jiexi"
}
