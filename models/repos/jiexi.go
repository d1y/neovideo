package repos

import "d1y.io/neovideo/models"

type IJiexi struct {
	Name string `json:"name,omitempty" gorm:"name" form:"name"`
	URL  string `json:"url,omitempty" gorm:"url;unique" form:"url"`
	Note string `json:"note,omitempty" gorm:"note" form:"note"`
}

type JiexiRepo struct {
	models.BaseRepo
	IJiexi
}

func (j *JiexiRepo) TableName() string {
	return "t_jiexi"
}
