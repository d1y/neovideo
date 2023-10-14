package sqls

import (
	"d1y.io/neovideo/models/other"
	"d1y.io/neovideo/models/repos"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AutoMigrate() {
	if err := db.AutoMigrate(&repos.MacCMSRepo{}, &repos.JiexiRepo{}, &repos.VideoRepo{}, &repos.VideoCategoryRepo{}, &other.ImageCoverTask{}, &other.SpiderTask{}); err != nil {
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
	},
	// {
	// 	ID: "202310062000",
	// 	Migrate: func(tx *gorm.DB) error {
	// 		tx.Model(&repos.MacCMSRepo{}).Where("1 = 1").UpdateColumn("r18", true)
	// 		tx.Model(&repos.VideoRepo{}).Where("1 = 1").UpdateColumn("r18", true)
	// 		return nil
	// 	},
	// },
	// {
	// 	ID: "202310141306",
	// 	Migrate: func(tx *gorm.DB) error {
	// 		var videos []repos.VideoRepo
	// 		if err := tx.Model(&repos.VideoRepo{}).Where("1 = 1").Find(&videos).Error; err != nil {
	// 			return err
	// 		}
	// 		for _, video := range videos {
	// 			cover := strings.Replace(video.Cover, "public/", "", -1)
	// 			fmt.Println("video cover current replace to ", cover)
	// 			tx.Model(&video).UpdateColumn("cover", cover)
	// 		}
	// 		return nil
	// 	},
	// },
	// {
	// 	ID: "202310141639",
	// 	Migrate: func(tx *gorm.DB) error {
	// 		var videos []repos.VideoRepo
	// 		if err := tx.Debug().Model(&repos.VideoRepo{}).Where("1 = 1").Find(&videos).Error; err != nil {
	// 			return err
	// 		}
	// 		return tx.Transaction(func(tx *gorm.DB) error {
	// 			for _, video := range videos {
	// 				maccmsID, rtype, r18 := video.Mid, video.RealType, video.R18
	// 				var category repos.VideoCategoryRepo
	// 				var sk = tx.Debug().Model(&repos.VideoCategoryRepo{}).Where("name = ?", rtype).First(&category)
	// 				if err := sk.Error; err != nil {
	// 					if errors.Is(err, gorm.ErrRecordNotFound) {
	// 						cate := repos.VideoCategoryRepo{
	// 							IVideoCategory: repos.IVideoCategory{
	// 								Name:    rtype,
	// 								R18:     r18,
	// 								Sources: datatypes.NewJSONSlice([]uint{maccmsID}),
	// 							},
	// 						}
	// 						tx.Debug().Save(&cate)
	// 					}
	// 				} else {
	// 					var fd = false
	// 					for _, source := range category.Sources {
	// 						if source == maccmsID {
	// 							fd = true
	// 						}
	// 					}
	// 					if !fd {
	// 						category.Sources = append(category.Sources, maccmsID)
	// 						tx.Debug().Save(&category)
	// 					}
	// 				}
	// 			}
	// 			return nil
	// 		})
	// 	},
	// },
	}).Migrate()
}
