package spiderman

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"d1y.io/neovideo/spider/implement/maccms"
	"github.com/panjf2000/ants/v2"
)

type task struct {
	Successful bool
	Page       int
	Reason     string
	Vidoes     *[]maccms.IMacCMSListVideoItem
}

func (t *task) SetReason(msg string) {
	t.Reason = msg
}

func (t *task) SetVideos(v *[]maccms.IMacCMSListVideoItem) {
	t.Vidoes = v
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

func Exec(rtype string, api string) error {
	pool, err := ants.NewPool(240)
	if err != nil {
		panic(err)
	}
	defer pool.Release()
	// cms := maccms.New(rtype, api)
	// homeData, err := cms.GetHome(1)
	if err != nil {
		return err
	}
	// header := homeData.ListHeader
	var count = 50 //header.PageCount
	if count <= 1 {
		return errors.New("count <= 1, not need task exec")
	}
	var sm sync.Mutex
	var wg sync.WaitGroup
	var tasks = make([]*task, count)
	wg.Add(count)
	for i := 0; i < count; i++ {
		idx := i
		fmt.Println(idx)
		pool.Submit(taskWrapper(&sm, &wg, idx, rtype, api, tasks))
	}
	wg.Wait()
	b, e := json.Marshal(tasks)
	if e != nil {
		panic(e)
	}
	if err := os.WriteFile("1.json", b, 0755); err != nil {
		panic("写入文件失败")
	}
	return nil
}

func taskWrapper(sm *sync.Mutex, wg *sync.WaitGroup, idx int, rt string, api string, tasks []*task) func() {
	return func() {
		sm.Lock()
		defer sm.Unlock()
		defer wg.Done()
		fmt.Printf("[task] 开始爬取第%v任务\n", idx+1)
		cms := maccms.New(rt, api)
		task := newTask(idx)
		data, err := cms.GetHome(idx + 1)
		if err != nil {
			fmt.Printf("[task] 爬取第%v任务失败\n", idx+1)
			task.SetFail(err.Error())
		} else {
			fmt.Printf("爬取第%v任务成功, 共爬取到%v条数据\n", idx+1, len(data.Videos))
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
		}
		tasks[idx] = task
	}
}
