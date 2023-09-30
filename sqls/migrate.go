package sqls

import (
	"d1y.io/neovideo/models/repos"
)

func AutoMigrate() {
	if err := db.AutoMigrate(&repos.MacCMSRepo{}, &repos.JiexiRepo{}); err != nil {
		panic(err)
	}
}
