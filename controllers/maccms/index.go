package maccms

import (
	"io"
	"sort"
	"sync"
	"time"

	"d1y.io/neovideo/common/impl"
	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/models/web"
	"d1y.io/neovideo/pkgs/safeset"
	"d1y.io/neovideo/sqls"

	"github.com/acmestack/gorm-plus/gplus"
	"github.com/imroc/req/v3"
	"github.com/kataras/iris/v12"
)

type IMacCMSController struct {
	sm sync.Mutex
}

func (im *IMacCMSController) getList(ctx iris.Context) {
	list, gb := gplus.SelectList[repos.MacCMSRepo](nil)
	if gb.Error != nil {
		web.NewError(gb.Error)
		return
	}
	web.NewData(list).Build(ctx)
}

func (im *IMacCMSController) create(ctx iris.Context) {
	var form repos.IMacCMS
	ctx.ReadBody(&form)
	var data = repos.MacCMSRepo{IMacCMS: form}
	if err := gplus.Insert[repos.MacCMSRepo](&data).Error; err != nil {
		web.NewError(err)
		return
	}
	web.NewData[repos.MacCMSRepo](data).Build(ctx)
}

func (im *IMacCMSController) delete(ctx iris.Context) {
	id := ctx.Params().Get("id")
	if err := gplus.DeleteById[repos.MacCMSRepo](id).Error; err != nil {
		web.NewError(err)
		return
	}
	web.NewData(id).SetMessage("删除成功").Build(ctx)
}

func (im *IMacCMSController) queryDBMaccms2Map() (map[string]*repos.MacCMSRepo, error) {
	dbMaccms, gb := gplus.SelectList[repos.MacCMSRepo](nil)
	if gb.Error != nil {
		return nil, gb.Error
	}
	var m = make(map[string]*repos.MacCMSRepo)
	for _, item := range dbMaccms {
		m[item.Api] = item
	}
	return m, nil
}

func (im *IMacCMSController) batchImport(ctx iris.Context) {
	url := ctx.FormValueDefault("url", "")
	importData := ctx.FormValueDefault("data", "")
	if len(url) >= 3 && len(importData) == 0 {
		resp, err := req.Get(url) /* FIXME: verify url */
		if err != nil {
			web.NewError(err).Build(ctx)
			return
		}
		b, e := io.ReadAll(resp.Body)
		if e != nil {
			web.NewError(e).Build(ctx)
			return
		}
		importData = string(b)
	}
	if len(importData) == 0 {
		web.NewMessage("导入数据为空").SetSuccessWithBool(false).Build(ctx)
		return
	}
	m, err := im.queryDBMaccms2Map()
	if err != nil {
		web.NewMessage("查询本地数据错误").SetSuccessWithBool(false).Build(ctx)
		return
	}
	var cs []*repos.MacCMSRepo
	ss := safeset.New()
	for _, item := range impl.ParseMaccms(importData) {
		api := item.Api
		if _, ok := m[api]; ok {
			continue
		}
		if ss.Contains(api) {
			continue
		}
		ss.Add(api)
		cs = append(cs, &repos.MacCMSRepo{
			IMacCMS: repos.IMacCMS{
				Api:         item.Api,
				Name:        item.Name,
				R18:         item.R18,
				RespType:    item.RespType,
				JiexiURL:    item.JiexiURL,
				JiexiEnable: item.JiexiParse,
				Category:    []string{},
			},
		})
	}
	l := len(cs)
	if l < 1 {
		web.NewMessage("数据已经全部导入过 :)").SetSuccessWithBool(false).Build(ctx)
		return
	}
	if err := gplus.InsertBatch[repos.MacCMSRepo](cs).Error; err != nil {
		web.NewMessage("保存到数据库失败").SetSuccessWithBool(false).Build(ctx)
		return
	}
	web.NewData(l).Build(ctx)
}

func (im *IMacCMSController) checkList(list []*repos.MacCMSRepo) []map[string]any {
	var wg sync.WaitGroup
	wg.Add(len(list))
	var mm = make([]map[string]any, 0)
	for _, item := range list {
		go func(i *repos.MacCMSRepo) {
			defer wg.Done()
			_, err := req.C().SetTimeout(4 * time.Second).R().Get(i.Api)
			m := map[string]any{
				"id":   i.ID,
				"name": i.Name,
			}
			m["successful"] = true
			if err != nil {
				m["message"] = err.Error()
				m["successful"] = false
			} else {
				m["message"] = "请求成功"
			}
			mm = append(mm, m)
		}(item)
	}
	wg.Wait()
	sort.Slice(mm, func(i, j int) bool {
		return mm[i]["id"].(uint) < mm[j]["id"].(uint)
	})
	return mm
}

func (im *IMacCMSController) handlerCheck(ctx iris.Context) ([]map[string]any, bool) {
	im.sm.Lock()
	defer im.sm.Unlock()
	list, gb := gplus.SelectList[repos.MacCMSRepo](nil)
	if gb.Error != nil {
		web.NewError(gb.Error).Build(ctx)
		return nil, false
	}
	l := len(list)
	if l <= 0 {
		web.NewMessage("数据为空").SetSuccessWithBool(false).Build(ctx)
		return nil, false
	}
	mm := im.checkList(list)
	return mm, true
}

func (im *IMacCMSController) check(ctx iris.Context) {
	mm, ok := im.handlerCheck(ctx)
	if !ok {
		return
	}
	web.NewData(mm).Build(ctx)
}

func (im *IMacCMSController) checkAndSync(ctx iris.Context) {
	mm, ok := im.handlerCheck(ctx)
	if !ok {
		return
	}
	now := time.Now()
	var errs = make([]map[string]any, 0)
	for _, item := range mm {
		id, ok := item["id"].(uint)
		if !ok {
			continue
		}
		successful := item["successful"].(bool)
		cols := map[string]any{
			"last_check": now,
			"available":  successful,
		}
		if err := sqls.DB().Model(&repos.MacCMSRepo{}).Where("id = ?", id).Updates(cols).Error; err != nil {
			errs = append(errs, map[string]any{
				"id":    id,
				"error": err.Error(),
			})
		}
	}
	web.NewData(errs).Build(ctx)
}

func (im *IMacCMSController) removeUnavailable(ctx iris.Context) {
	query, u := gplus.NewQuery[repos.MacCMSRepo]()
	query.Eq(&u.Available, false)
	unavailable, gb := gplus.SelectList[repos.MacCMSRepo](query)
	if gb.Error != nil {
		web.NewError(gb.Error).Build(ctx)
	}
	if len(unavailable) <= 0 {
		web.NewMessage("数据为空").SetSuccessWithBool(false).Build(ctx)
		return
	}
	var ids = make([]uint, len(unavailable))
	for idx, item := range unavailable {
		ids[idx] = item.ID
	}
	if err := gplus.DeleteByIds[repos.MacCMSRepo](ids).Error; err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	web.NewData(unavailable).SetMessage("删除成功").Build(ctx)
}

func Register(u iris.Party) {
	var imc IMacCMSController
	var px IMacCMSProxyController
	u.Get("/", imc.getList).Name = "获取苹果CMS列表"
	u.Post("/", imc.create).Name = "创建苹果CMS"
	u.Post("/batch_import", imc.batchImport).Name = "批量导入苹果CMS列表"
	u.Delete("/{id:int}", imc.delete).Name = "删除苹果CMS"
	u.Post("/allcheck", imc.check).Name = "检查苹果CMS"
	u.Post("/allcheck/sync", imc.checkAndSync).Name = "检查苹果CMS并同步到服务器"
	u.Delete("/allcheck/unavailable", imc.removeUnavailable).Name = "删除不可用的苹果CMS"
	u.PartyFunc("/proxy", px.Register)
}
