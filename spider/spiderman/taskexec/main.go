package main

import (
	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/spider/spiderman"
	"github.com/acmestack/gorm-plus/gplus"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("/Users/d1y/code/github/neovideo/db.sqlite3"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	gplus.Init(db)
	tasks, err := spiderman.Exec("xml", "http://156.249.29.8/inc/api.php")
	if err != nil {
		panic(err)
	}
	var list []*repos.VideoRepo
	for _, item := range tasks {
		if !item.Successful {
			continue
		}
		for _, subItem := range *item.Videos {
			var value = repos.VideoRepo{
				IVideo: repos.IVideo{
					SpiderType: "maccms",
					Title:      subItem.Name,
					Desc:       subItem.Desc,
					RealID:     subItem.Id,
					RealTime:   subItem.Last,
					RealCover:  subItem.Pic,
					Cover:      "",
					CategoryID: subItem.Tid,
					Lang:       subItem.Lang,
					Area:       subItem.Area,
					Year:       subItem.Year,
					State:      subItem.State,
					Actor:      subItem.Actor,
					Director:   subItem.Director,
				},
			}
			var videos = make([]repos.IVideoDataInfo, 0)
			for _, d := range subItem.DD {
				var vc = repos.IVideoDataInfo{
					Flag: d.Flag,
				}
				for _, dv := range d.Videos {
					vc.Videos = append(vc.Videos, repos.IVideoData{
						Name:  dv.Name,
						URL:   dv.URL,
						Embed: dv.Embed,
					})
				}
				videos = append(videos, vc)
			}
			value.Videos = datatypes.NewJSONSlice[repos.IVideoDataInfo](videos)
			list = append(list, &value)
		}
	}
	if err := gplus.InsertBatch[repos.VideoRepo](list).Error; err != nil {
		panic(err)
	}
}
