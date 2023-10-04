package main

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/spider/spiderman"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dir = "./public"
)

func ensureDir(baseDir string) error {
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, 0755)
}

func init() {
	ensureDir(dir)
}

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
		var wg sync.WaitGroup
		for _, subItem := range *item.Videos {
			cover := ""
			if len(subItem.Pic) >= 1 {
				wg.Add(1)
				cover = createFilename(subItem.Pic)
				go imageDownload(subItem.Pic, cover, &wg)
			}
			var value = repos.VideoRepo{
				IVideo: repos.IVideo{
					SpiderType: "maccms",
					Title:      subItem.Name,
					Desc:       subItem.Desc,
					RealID:     subItem.Id,
					RealTime:   subItem.Last,
					RealCover:  subItem.Pic,
					Cover:      cover,
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
		wg.Wait()
	}
	if err := gplus.InsertBatch[repos.VideoRepo](list).Error; err != nil {
		panic(err)
	}
}

func createFilename(url string) string {
	ext := filepath.Ext(url)
	uuid := uuid.New()
	filename := uuid.String()
	filename += ext
	path := filepath.Join(dir, filename)
	return path
}

func imageDownload(url string, filename string, wg *sync.WaitGroup) error {
	defer wg.Done()
	resp, err := req.Get(url)
	if err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
