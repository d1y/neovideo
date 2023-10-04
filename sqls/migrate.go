package sqls

import (
	"d1y.io/neovideo/models/repos"
)

func AutoMigrate() {
	if err := db.AutoMigrate(&repos.MacCMSRepo{}, &repos.JiexiRepo{}, &repos.VideoRepo{}); err != nil {
		panic(err)
	}
}
