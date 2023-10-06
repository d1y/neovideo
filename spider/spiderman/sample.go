package spiderman

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"d1y.io/neovideo/config/constant"
	"d1y.io/neovideo/models/other"
	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/spider/implement/maccms"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

func GetTaskMsg() string {
	return ct.String()
}

func Stop() int {
	if pool == nil {
		return 0
	}
	pool.Release()
	return pool.Running()
}

func IsStart() bool {
	return pool.Running() >= 1 || ct.IsStart()
}

func Start(cs *repos.MacCMSRepo) (any, error) {
	if len(cs.RespType) <= 0 || len(cs.Api) <= 0 {
		return nil, errors.New("maccms repo api or type is empty")
	}
	if IsStart() {
		return nil, ErrRepeatTask
	}
	ct.Padding()
	defer ct.Reset()
	err := newPool()
	if err != nil {
		return nil, err
	}
	defer pool.Release()
	cms := maccms.New(cs.RespType, cs.Api)
	homeData, err := cms.GetHome(1)
	if err != nil {
		return nil, err
	}
	header := homeData.ListHeader
	var count = header.PageCount
	if count <= 1 {
		return nil, errors.New("count <= 1, not need task exec")
	}
	ct.SetCount(count)
	var sm sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		idx := i
		wg.Add(1)
		pool.Submit(taskWrapper(&sm, &wg, idx, cs))
	}
	wg.Wait()
	return nil, nil
}

func taskWrapper(sm *sync.Mutex, wg *sync.WaitGroup, idx int, cs *repos.MacCMSRepo) func() {
	return func() {
		sm.Lock()
		defer sm.Unlock()
		defer wg.Done()
		defer ct.Increase()
		logrus.Printf("[task] 开始爬取第%v任务", idx+1)
		cms := maccms.New(cs.RespType, cs.Api)
		task := newTask(idx)
		data, err := cms.GetHome(idx + 1)
		if err != nil {
			logrus.Errorf("[task] 爬取第%v任务失败", idx+1)
			task.SetFail(err.Error())
		} else {
			logrus.Printf("[task] 爬取第%v任务成功, 共爬取到%v条数据", idx+1, len(data.Videos))
			ids := make([]int, len(data.Videos))
			for idx, item := range data.Videos {
				ids[idx] = item.Id
			}
			c := maccms.New(cs.RespType, cs.Api)
			_, v, err := c.GetDetail(ids...)
			if err == nil {
				task.SetSuccessful(&v)
			} else {
				task.SetFail(err.Error())
			}
			insertData(task, cs)
		}
	}
}

func insertData(item *task, cs *repos.MacCMSRepo) {
	if !item.Successful {
		idt := other.NewSpiderTask(cs.ID, item.Page)
		idt.SetFailed(item.Reason)
		gplus.Insert[other.SpiderTask](idt)
		logrus.Errorf("[task] 插入爬虫任务错误(页数: %d) 错误信息: %s", item.Page, item.Reason)
		return
	}
	var wg sync.WaitGroup
	for _, subItem := range *item.Videos {
		cover := ""
		if len(subItem.Pic) >= 1 {
			wg.Add(1)
			cover = createFilename(subItem.Pic)
			go func(url *string, filename *string, wg *sync.WaitGroup) {
				if err := imageDownload(*url, *filename, wg); err != nil {
					cv := other.NewCoverTask(*url, *filename, err)
					gplus.Insert(cv)
				}
			}(&subItem.Pic, &cover, &wg)
		}
		var value = repos.VideoRepo{
			IVideo: repos.IVideo{
				SpiderType: "maccms",
				Title:      subItem.Name,
				Desc:       subItem.Desc,
				Mid:        cs.ID,
				R18:        cs.R18,
				RealType:   subItem.Type,
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
		if err := gplus.Insert[repos.VideoRepo](&value).Error; err != nil { /* 忽略这里的错误 */
			logrus.Errorln(err)
		}
	}
	idt := other.NewSpiderTask(cs.ID, item.Page)
	msg := fmt.Sprintf("[task] 本次任务插入%d条数据 当前页数%d", len(*item.Videos), item.Page)
	idt.SetSuccessful(msg)
	gplus.Insert[other.SpiderTask](idt)
	logrus.Info(msg)
	wg.Wait()
}

func createFilename(url string) string {
	ext := filepath.Ext(url)
	uuid := uuid.New()
	filename := uuid.String()
	filename += ext
	path := filepath.Join(constant.Public, filename)
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
