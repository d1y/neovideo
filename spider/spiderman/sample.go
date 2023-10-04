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
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

var (
	ErrRepeatTask = errors.New("repeat task")
)

var pool *ants.Pool

type Counter struct {
	sm      sync.Mutex
	count   int
	current int
}

func (ct *Counter) Padding() {
	ct.sm.Lock()
	defer ct.sm.Unlock()
	ct.count = 0
	ct.current = 0
}

func (ct *Counter) SetCount(count int) {
	ct.sm.Lock()
	defer ct.sm.Unlock()
	ct.count = count
}

func (ct *Counter) Increase() {
	ct.sm.Lock()
	defer ct.sm.Unlock()
	ct.current++
}

func (ct *Counter) Reset() {
	ct.sm.Lock()
	defer ct.sm.Unlock()
	ct.count = -1
	ct.current = -1
}

func (ct *Counter) GetCount() int {
	ct.sm.Lock()
	defer ct.sm.Unlock()
	return ct.count
}

func (ct *Counter) GetCurrent() int {
	ct.sm.Lock()
	defer ct.sm.Unlock()
	return ct.current
}

func (ct *Counter) String() string {
	if ct.GetCurrent() == 0 && ct.GetCount() == 0 {
		return "任务准备中(或已失败)"
	}
	var msg = fmt.Sprintf("%d/%d", ct.GetCurrent(), ct.GetCount())
	return msg
}

func (ct *Counter) IsStart() bool {
	ct.sm.Lock()
	defer ct.sm.Unlock()
	if ct.count == -1 || ct.current == -1 {
		return false
	}
	if ct.count == 0 && ct.current == 0 {
		return true
	}
	return ct.count >= ct.current
}

var ct = Counter{
	count:   -1,
	current: -1,
	sm:      sync.Mutex{},
}

func GetTaskMsg() string {
	return ct.String()
}

type task struct {
	Successful bool
	Page       int
	Reason     string
	Videos     *[]maccms.IMacCMSListVideoItem
}

func (t *task) SetReason(msg string) {
	t.Reason = msg
}

func (t *task) SetVideos(v *[]maccms.IMacCMSListVideoItem) {
	t.Videos = v
}

func (t *task) SetSuccessful(v *[]maccms.IMacCMSListVideoItem) {
	t.Successful = true
	t.SetVideos(v)
}

func (t *task) SetFail(reason string) {
	t.Successful = false
	t.SetReason(reason)
}

func newTask(page int) *task {
	return &task{
		Page: page,
	}
}

func IsStart() bool {
	return ct.IsStart()
}

func Start(rtype, api string, id uint) (any, error) {
	if IsStart() {
		return nil, ErrRepeatTask
	}
	ct.Padding()
	defer ct.Reset()
	var err error
	pool, err = ants.NewPool(240)
	if err != nil {
		return nil, err
	}
	defer pool.Release()
	cms := maccms.New(rtype, api)
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
	wg.Add(count)
	for i := 0; i < count; i++ {
		idx := i
		pool.Submit(taskWrapper(&sm, &wg, idx, rtype, api, id))
	}
	wg.Wait()
	return nil, nil
}

func taskWrapper(sm *sync.Mutex, wg *sync.WaitGroup, idx int, rt string, api string, id uint) func() {
	return func() {
		sm.Lock()
		defer sm.Unlock()
		defer wg.Done()
		defer ct.Increase()
		logrus.Printf("[task] 开始爬取第%v任务", idx+1)
		cms := maccms.New(rt, api)
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
			c := maccms.New(rt, api)
			_, v, err := c.GetDetail(ids...)
			if err == nil {
				task.SetSuccessful(&v)
			} else {
				task.SetFail(err.Error())
			}
			insertData(task, id)
		}
	}
}

func insertData(item *task, sid uint) {
	if !item.Successful {
		idt := other.NewSpiderTask(sid, item.Page)
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
					cv := other.NewCoverTask(*url, *filename)
					gplus.Insert(cv)
				}
			}(&subItem.Pic, &cover, &wg)
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
		if err := gplus.Insert[repos.VideoRepo](&value).Error; err != nil { /* 忽略这里的错误 */
			logrus.Errorln(err)
		}
	}
	idt := other.NewSpiderTask(sid, item.Page)
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
