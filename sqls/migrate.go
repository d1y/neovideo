package sqls

import (
	"d1y.io/neovideo/models/repos"
)

func AutoMigrate() {
	db.AutoMigrate(&repos.MacCMSRepo{}, &repos.JiexiRepo{})
}
