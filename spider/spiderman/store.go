package spiderman

import (
	"errors"
	"fmt"
	"sync"

	"d1y.io/neovideo/spider/implement/maccms"
	"github.com/panjf2000/ants/v2"
)

var (
	ErrRepeatTask = errors.New("repeat task")
)

var pool *ants.Pool

// 目前来说,实际上没啥用, 因为都是等上个任务完成再继续的
// 它的数量决定任务暂停的速度
// 已知超时检测时间为 3秒
// 则最多要 poolSize * 3 秒任务才会暂停 :)
const poolSize = 6

type Counter struct {
	sm      sync.RWMutex
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
	ct.sm.RLock()
	defer ct.sm.RUnlock()
	return ct.count
}

func (ct *Counter) GetCurrent() int {
	ct.sm.RLock()
	defer ct.sm.RUnlock()
	return ct.current
}

func (ct *Counter) String() string {
	if ct.GetCurrent() == 0 && ct.GetCount() == 0 {
		return "任务准备中(或已失败)"
	}
	var msg = fmt.Sprintf("%d/%d(协程数量%d)", ct.GetCurrent(), ct.GetCount(), pool.Running())
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
	sm:      sync.RWMutex{},
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

func newPool() (err error) {
	pool, err = ants.NewPool(poolSize)
	return
}
