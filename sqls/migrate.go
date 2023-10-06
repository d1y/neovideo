package sqls

import (
	"d1y.io/neovideo/models/other"
	"d1y.io/neovideo/models/repos"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AutoMigrate() {
	if err := db.AutoMigrate(&repos.MacCMSRepo{}, &repos.JiexiRepo{}, &repos.VideoRepo{}, &other.ImageCoverTask{}, &other.SpiderTask{}); err != nil {
		panic(err)
	}
}

func MigrateBatch() {
	gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{{
		ID: "20231006",
		Migrate: func(tx *gorm.DB) error {
			var list []repos.VideoRepo
			if err := tx.Model(&repos.VideoRepo{}).Find(&list).Error; err != nil {
				return err
			}
			var m = make(map[int]uint)
			for _, item := range list {
				if item.RealID <= 0 {
					continue
				}
				m[item.RealID] = item.ID
			}
			for _, item := range list {
				_, ok := m[item.RealID]
				if !ok {
					logrus.Printf("[db] 删除重复的视频(%d)", item.ID)
					if err := tx.Model(&repos.VideoRepo{}).Where("id = ?", item.ID).Delete(nil).Error; err != nil {
						return err
					}
				}
			}
			return nil
		},
	}}).Migrate()
}
