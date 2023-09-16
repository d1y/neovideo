package repos

type JiexiRepo struct {
	Name string `json:"name,omitempty" gorm:"name"`
	Url  string `json:"url,omitempty" gorm:"url"`
	Note string `json:"note,omitempty" gorm:"note"`
}
