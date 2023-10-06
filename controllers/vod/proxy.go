package vod

import (
	"errors"
	"sort"
	"sync"

	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/models/web"
	"d1y.io/neovideo/spider/implement/maccms"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/kataras/iris/v12"
)

const (
	homeRepoQueryKey = "$home"
	homeRenderKey    = "$hrender"
)

type homeItem struct {
	ID    uint                   `json:"id"`
	Api   string                 `json:"api"`
	Name  string                 `json:"name"`
	Data  maccms.IMacCMSHomeData `json:"data,omitempty"`
	Error string                 `json:"error,omitempty"`
}

func (vc *IVodController) queryRawCMS() ([]repos.MacCMSRepo, error) {
	vc.sm.Lock()
	defer vc.sm.Unlock()
	var result []repos.MacCMSRepo
	if val, ok := vc.cc.Get(homeRepoQueryKey); ok {
		if v, o := val.([]repos.MacCMSRepo); o {
			result = v
		}
	} else {
		cms, gb := gplus.SelectList[repos.MacCMSRepo](nil)
		if gb.Error != nil {
			return nil, gb.Error
		}
		for _, item := range cms {
			result = append(result, *item)
		}
		vc.cc.SetDefault(homeRepoQueryKey, result)
	}
	return result, nil
}

func (vc *IVodController) queryAndCMSFetchHome() ([]homeItem, error) {
	c, e := vc.queryRawCMS()
	if e != nil {
		return nil, e
	}
	if len(c) <= 0 {
		return nil, errors.New("maccms is empty")
	}
	var wg sync.WaitGroup
	var data []homeItem
	wg.Add(len(c))
	for _, item := range c {
		go func(item repos.MacCMSRepo) {
			defer wg.Done()
			var im = homeItem{
				Name: item.Name,
				ID:   item.ID,
				Api:  item.Api,
			}
			val, err := maccms.New(item.RespType, item.Api).GetHome(1)
			if err != nil {
				im.Error = err.Error()
			} else {
				im.Data = val
			}
			data = append(data, im)
		}(item)
	}
	wg.Wait()
	var result = make([]homeItem, 0)
	for _, item := range data {
		if item.Error == "" {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result, nil
}

func (vc *IVodController) renderHome(ctx iris.Context) {
	var result []homeItem
	if val, ok := vc.cc.Get(homeRenderKey); ok {
		if v, o := val.([]homeItem); o {
			result = v
		}
	} else {
		ims, err := vc.queryAndCMSFetchHome()
		if err != nil {
			web.NewError(err).Build(ctx)
			return
		}
		vc.cc.SetDefault(homeRenderKey, ims)
		result = ims
	}
	web.NewData(result).Build(ctx)
}
