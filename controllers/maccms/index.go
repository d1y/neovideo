package maccms

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"sync"
	"time"

	"d1y.io/neovideo/common/impl"
	"d1y.io/neovideo/controllers/handler"
	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/models/web"
	"d1y.io/neovideo/pkgs/safeset"
	"d1y.io/neovideo/spider/implement/maccms"
	"d1y.io/neovideo/sqls"
	"gorm.io/datatypes"

	"github.com/acmestack/gorm-plus/gplus"
	"github.com/beevik/etree"
	"github.com/imroc/req/v3"
	"github.com/kataras/iris/v12"
	"github.com/tidwall/gjson"
)

// 域名过期返回html格式(可能不准确)
//
// 参考返回值: https://www.77zy.vip/inc/m3u8.php
var malwaredomainlistHTMLRegexp = regexp.MustCompile(`^<!doctype html>[\s\S]*data-adblockkey[\s\S]*>window\.park =`)

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
	web.NewData[repos.MacCMSRepo](data).SetMessage("创建成功").Build(ctx)
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
	importData, err := handler.NewImportDataWithContext(ctx)
	if err != nil {
		web.NewError(err).Build(ctx)
		return
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
			},
		})
	}
	l := len(cs)
	if l < 1 {
		web.NewMessage("数据已经全部导入过 :)").SetSuccessWithBool(false).Build(ctx)
		return
	}
	if err := gplus.InsertBatch[repos.MacCMSRepo](cs).Error; err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	web.NewData(l).SetMessage(fmt.Sprintf("导入成功(%d)", l)).Build(ctx)
}

func (im *IMacCMSController) checkOnce(i *repos.MacCMSRepo) map[string]any {
	resp, err := req.C().SetTimeout(4 * time.Second).R().Get(i.Api)
	sfu := true
	if err == nil {
		sfu = resp.StatusCode == iris.StatusOK // 状态只要不是200就是错误!
	}
	m := map[string]any{
		"id":   i.ID,
		"name": i.Name,
		"api":  i.Api,
	}
	m["successful"] = sfu
	m["type"] = "unknown"
	if err != nil || !sfu {
		if err != nil {
			m["message"] = err.Error()
		} else {
			m["message"] = "response status is not 200!"
		}
		m["successful"] = false
	} else {
		b, e := io.ReadAll(resp.Body)
		s := string(b)
		if e != nil || malwaredomainlistHTMLRegexp.MatchString(s) {
			if e != nil {
				m["message"] = e.Error()
			} else {
				m["message"] = "response parse fail(domain expired)"
			}
			m["successful"] = false
		} else {
			rt := maccms.GetResponseType(s)
			av := maccms.NewWithApi(i.Api)
			if rt.IsJSON() {
				if gjson.Valid(s) {
					gp := gjson.Parse(s)
					_, _, categorys := av.JsonParseBody(gp)
					if len(categorys) >= 1 {
						m["category"] = datatypes.NewJSONSlice(categorys)
					}
				}
			} else {
				doc := etree.NewDocument()
				if err := doc.ReadFromString(s); err == nil {
					root := doc.Root()
					if root != nil {
						if categorys := av.XMLGetCategoryWithEtreeRoot(root); len(categorys) >= 1 {
							m["category"] = datatypes.NewJSONSlice(categorys)
						}
					}
				}
			}
			m["type"] = string(rt)
			m["message"] = "请求成功"
		}
	}
	return m
}

func (im *IMacCMSController) checkList(list []*repos.MacCMSRepo) []map[string]any {
	var wg sync.WaitGroup
	wg.Add(len(list))
	var mm = make([]map[string]any, 0)
	for _, item := range list {
		go func(i *repos.MacCMSRepo) {
			defer wg.Done()
			mm = append(mm, im.checkOnce(i))
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
	id, _ := handler.NewIDWithContext(ctx)
	val, gb := gplus.SelectById[repos.MacCMSRepo](id)
	if gb.Error != nil {
		web.NewError(gb.Error).Build(ctx)
		return
	}
	web.NewData(im.checkOnce(val)).SetMessage("执行成功").Build(ctx)
}

func (im *IMacCMSController) allcheck(ctx iris.Context) {
	mm, ok := im.handlerCheck(ctx)
	if !ok {
		return
	}
	web.NewData(mm).Build(ctx)
}

func (im *IMacCMSController) allcheckAndSync(ctx iris.Context) {
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
		resType := item["type"].(string)
		cols := map[string]any{
			"last_check": now,
			"available":  successful,
		}
		if resType == maccms.MacCMSReponseTypeJSON {
			cols["resp_type"] = maccms.MacCMSReponseTypeJSON
		} else if resType == maccms.MacCMSReponseTypeXML {
			cols["resp_type"] = maccms.MacCMSReponseTypeXML
		}
		if c, ok := item["category"]; ok {
			if j, o := c.(datatypes.JSONSlice[repos.IMacCMSCategory]); o {
				cols["category"] = j
			}
		}
		if err := sqls.DB().Model(&repos.MacCMSRepo{}).Where("id = ?", id).Updates(cols).Error; err != nil {
			errs = append(errs, map[string]any{
				"id":    id,
				"error": err.Error(),
			})
		}
	}
	web.NewData(errs).SetMessage("同步成功").Build(ctx)
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
	u.Post("/check/{id:int}", imc.check)
	u.Post("/allcheck", imc.allcheck).Name = "检查苹果CMS"
	u.Post("/allcheck/sync", imc.allcheckAndSync).Name = "检查苹果CMS并同步到服务器"
	u.Delete("/allcheck/unavailable", imc.removeUnavailable).Name = "删除不可用的苹果CMS"
	u.PartyFunc("/proxy", px.Register)
}
