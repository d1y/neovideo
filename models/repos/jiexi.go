package repos

import "gorm.io/gorm"

type JiexiRepo struct {
	Name string `json:"name,omitempty" gorm:"name" form:"name"`
	URL  string `json:"url,omitempty" gorm:"url;unique" form:"url"`
	Note string `json:"note,omitempty" gorm:"note" form:"note"`
}

type JiexiRepoWithModel struct {
	gorm.Model
	JiexiRepo
}
