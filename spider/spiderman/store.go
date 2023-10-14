package spiderman

import (
	"errors"
	"fmt"

	"d1y.io/neovideo/common/safeint"
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

type ICounter struct {
	current *safeint.SafeInt
	count   *safeint.SafeInt
}

func (c *ICounter) Padding() {
	c.current.SetZero()
	c.count.SetZero()
}

func (c *ICounter) SetCount(count int) {
	c.count.Set(count)
}

func (c *ICounter) Increase() {
	c.current.Add(1)
}

func (c *ICounter) Reset() {
	c.count.Set(-1)
	c.current.Set(-1)
}

func (c *ICounter) GetCount() int {
	return c.count.Get()
}

func (c *ICounter) GetCurrent() int {
	return c.current.Get()
}

func (c *ICounter) String() string {
	if c.GetCurrent() == 0 && c.GetCount() == 0 {
		return "任务准备中(或已失败)"
	}
	var msg = fmt.Sprintf("%d/%d(协程数量%d)", c.GetCurrent(), c.GetCount(), pool.Running())
	return msg
}

func (c *ICounter) IsStart() bool {
	if c.count.Get() == -1 || c.current.Get() == -1 {
		return false
	}
	if c.count.Get() == 0 && c.current.Get() == 0 {
		return true
	}
	return c.count.Get() >= c.current.Get()
}

var ct = &ICounter{
	current: safeint.New(-1),
	count:   safeint.New(-1),
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
