package sqls

import (
	"d1y.io/neovideo/models/other"
	"d1y.io/neovideo/models/repos"
)

func AutoMigrate() {
	if err := db.AutoMigrate(&repos.MacCMSRepo{}, &repos.JiexiRepo{}, &repos.VideoRepo{}, &other.ImageCoverTask{}, &other.SpiderTask{}); err != nil {
		panic(err)
	}
}
