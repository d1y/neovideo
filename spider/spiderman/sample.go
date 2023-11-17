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
	"d1y.io/neovideo/sqls"
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
	if pool != nil {
		return pool.Running() >= 1
	}
	return ct.IsStart()
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
		taskID, suf, ok /*理论来说 ok 是不应该存在的 */ := querySpiderTask(cs.ID, idx)
		if ok && suf {
			logrus.Printf("[task] 该任务已运行过(第%v页)", idx+1)
			return
		}
		logrus.Printf("[task] 开始爬取第%v页任务", idx+1)
		cms := maccms.New(cs.RespType, cs.Api)
		spiderTask := newTask(idx)
		data, err := cms.GetHome(idx + 1)
		if err != nil {
			logrus.Errorf("[task] 爬取第%v页任务失败", idx+1)
			spiderTask.SetFail(err.Error())
		} else {
			logrus.Printf("[task] 爬取第%v页任务成功, 共爬取到%v条数据", idx+1, len(data.Videos))
			ids := make([]int, len(data.Videos))
			for idx, item := range data.Videos {
				ids[idx] = item.Id
			}
			c := maccms.New(cs.RespType, cs.Api)
			_, v, err := c.GetDetail(ids...)
			if err == nil {
				spiderTask.SetSuccessful(&v)
			} else {
				spiderTask.SetFail(err.Error())
			}
			skipInsertSpiderTask := ok && !suf
			insertData(spiderTask, cs, skipInsertSpiderTask, taskID)
		}
	}
}

func insertData(item *task, cs *repos.MacCMSRepo, skipInsertSpiderTask bool, taskID uint) {
	if !item.Successful && !skipInsertSpiderTask {
		idt := other.NewSpiderTask(cs.ID, item.Page)
		idt.SetFailed(item.Reason)
		gplus.Insert[other.SpiderTask](idt)
		logrus.Errorf("[task] 插入爬虫任务错误(页数: %d) 错误信息: %s", item.Page, item.Reason)
		return
	}
	var wg sync.WaitGroup
	for _, subItem := range *item.Videos {
		var video repos.VideoRepo
		if err := sqls.DB().Where(&repos.VideoRepo{
			IVideo: repos.IVideo{
				Mid:    cs.ID,
				RealID: subItem.Id,
			},
		}).Find(&video).Error; err == nil {
			video.Videos = maccmsDD2Videos(subItem.DD)
			if err := sqls.DB().Save(&video).Error; err != nil {
				logrus.Errorf("[task] 更新数据失败(%d): %v", video.ID, err)
			} else {
				logrus.Printf("[task] 更新数据成功(%d)", video.ID)
			}
			continue
		}
		cover := ""
		if len(subItem.Pic) >= 1 {
			wg.Add(1)
			cover = createFilename(subItem.Pic)
			go func(url *string, filename *string, wg *sync.WaitGroup) {
				fullpath := filepath.Join(constant.Public, *filename)
				if err := imageDownload(*url, fullpath, wg); err != nil {
					cv := other.NewCoverTask(*url, *filename, err)
					gplus.Insert(cv)
				}
			}(&subItem.Pic, &cover, &wg)
		}
		v := maccmsItem2video(cs, &subItem, cover)
		if err := gplus.Insert[repos.VideoRepo](v).Error; err != nil {
			logrus.Errorln(err)
		}
	}
	if !skipInsertSpiderTask {
		idt := other.NewSpiderTask(cs.ID, item.Page)
		msg := fmt.Sprintf("[task] 本次任务插入%d条数据 当前页数%d", len(*item.Videos), item.Page)
		idt.SetSuccessful(msg)
		gplus.Insert[other.SpiderTask](idt)
		logrus.Info(msg)
	} else {
		if err := updateSpiderTaskWithSuccess(taskID); err != nil {
			logrus.Errorf("[task] 更新爬虫任务状态(%d)失败%v", taskID, err)
		} else {
			logrus.Infof("[task] 更新爬虫任务状态(%d)成功", taskID)
		}
	}
	wg.Wait()
}

func maccmsItem2video(cs *repos.MacCMSRepo, subItem *maccms.IMacCMSListVideoItem, cover string) *repos.VideoRepo {
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
	value.Videos = maccmsDD2Videos(subItem.DD)
	return &value
}

func maccmsDD2Videos(dd []maccms.IMacCMSVideoDDTag) datatypes.JSONSlice[repos.IVideoDataInfo] {
	var videos = make([]repos.IVideoDataInfo, len(dd))
	for index, d := range dd {
		var vc = repos.IVideoDataInfo{
			Flag:   d.Flag,
			Videos: make([]repos.IVideoData, len(d.Videos)),
		}
		for idx, dv := range d.Videos {
			vc.Videos[idx] = repos.IVideoData{
				Name:  dv.Name,
				URL:   dv.URL,
				Embed: dv.Embed,
			}
		}
		videos[index] = vc
	}
	return datatypes.NewJSONSlice[repos.IVideoDataInfo](videos)
}

func updateSpiderTaskWithSuccess(mid uint /*, page int*/) error {
	return sqls.DB().Model(&other.SpiderTask{}).Where("id = ?", mid).UpdateColumn("successful", true).Error
}

func querySpiderTask(mid uint, idx int) (taskID uint, successful bool, ok bool) {
	q, st := gplus.NewQuery[other.SpiderTask]()
	q.Eq(&st.Mid, mid)
	q.Eq(&st.Page, idx)
	spiderTask, gb := gplus.SelectOne(q)
	if gb.Error != nil { /* ignore error stack */
		ok = false
		return
	}
	ok = true
	taskID = spiderTask.ID
	successful = spiderTask.Successful
	return
}

func createFilename(url string) string {
	ext := filepath.Ext(url)
	guuid := uuid.New()
	file := guuid.String()
	file += ext
	return file
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
