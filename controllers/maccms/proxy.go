package maccms

import (
	"errors"
	"strconv"
	"time"

	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/models/web"
	"d1y.io/neovideo/pkgs/stringbuilder"
	"d1y.io/neovideo/spider/implement/maccms"
	"d1y.io/neovideo/sqls"
	"github.com/kataras/iris/v12"
	"github.com/patrickmn/go-cache"
)

const (
	idSep = "$$"
)

type sbIDFunc = func(val string) (*stringbuilder.SB, string)

func i2s(i int) string {
	return strconv.Itoa(i)
}

type IMacCMSProxyController struct {
	cc *cache.Cache
}

func (pc *IMacCMSProxyController) Register(u iris.Party) {
	pc.cc = cache.New(42*time.Second, 60*time.Second)
	u.Post("/{id:int}", pc.request).Name = "代理访问苹果CMS"
}

func (pc *IMacCMSProxyController) re2ID(reqAction int) sbIDFunc {
	sb := stringbuilder.New(12)
	sb.AppendArray(i2s(reqAction), idSep)
	fn := func(val string /* TODO: add int type support(string | int) */) (*stringbuilder.SB, string) {
		sb.AppendArray(val, idSep)
		return sb, sb.String()
	}
	return fn
}

func (pc *IMacCMSProxyController) getHomeCacheID(id int, page int, category int) string {
	sb, _ := pc.re2ID(proxyActionWithHome)(i2s(id))
	sb.AppendArray(strconv.Itoa(page), idSep)
	sb.AppendArray(strconv.Itoa(category), idSep)
	return sb.String()
}

func (pc *IMacCMSProxyController) getCategoryCacheID(id int) (s string) {
	_, s = pc.re2ID(proxyActionWithCategory)(i2s(id))
	return
}

func (pc *IMacCMSProxyController) getDetailCacheID(id, detailID int) string {
	s, _ := pc.re2ID(proxyActionWithDetail)(i2s(id))
	s.AppendArray(i2s(detailID), idSep)
	return s.String()
}

func (pc *IMacCMSProxyController) getSearchCacheID(id int, keyword string /*, page int*/) string {
	s, _ := pc.re2ID(proxyActionWithSearch)(i2s(id))
	s.AppendArray(keyword /* FIXME: keyword maybe use $$(idSep), need reclean */, idSep)
	return s.String()
}

func reUseCache[T any](cc *cache.Cache, k string, output *any) bool {
	if val, ok := cc.Get(k); ok {
		if v, o := val.(maccms.IMacCMSHomeData); o {
			*output = v
			return true
		}
	}
	return false
}

func (pc *IMacCMSProxyController) setResult2Cache(err error, k string, val any) {
	if err == nil {
		pc.cc.SetDefault(k, val)
	}
}

func (pc *IMacCMSProxyController) realRequest(data repos.MacCMSRepo, req maccms.XHRRequest, id int, alwayFetch bool) (any, error) {
	var result any
	var err error = nil
	cms := maccms.New(data.RespType, data.Api)
	page := req.Page
	if page == 0 {
		page = 1
	}
	switch req.RequestAction {
	case proxyActionWithHome:
		var category = req.Category
		k := pc.getHomeCacheID(id, page, category)
		if ok := reUseCache[maccms.IMacCMSHomeData](pc.cc, k, &result); !ok || alwayFetch {
			result, err = cms.GetHome(page, category)
			pc.setResult2Cache(err, k, result)
		}
	case proxyActionWithCategory:
		k := pc.getCategoryCacheID(id)
		if ok := reUseCache[[]repos.IMacCMSCategory](pc.cc, k, &result); !ok || alwayFetch {
			result, err = cms.GetCategory()
			pc.setResult2Cache(err, k, result)
		}
	case proxyActionWithDetail:
		ids := req.GetIDs2Slice()
		if len(ids) < 1 {
			err = errors.New("proxy fetch detail faild(ids)")
		} else {
			detailID := ids[0]
			k := pc.getDetailCacheID(id, detailID)
			if ok := reUseCache[[]maccms.IMacCMSListVideoItem](pc.cc, k, &result); !ok || alwayFetch {
				_, result, err = cms.GetDetail(detailID /* TODO: add multiple id */)
				pc.setResult2Cache(err, k, result)
			}
		}
	case proxyActionWithSearch:
		k := pc.getSearchCacheID(id, req.Keyword)
		if ok := reUseCache[maccms.IMacCMSVideosAndHeader](pc.cc, k, &result); !ok || alwayFetch {
			result, err = cms.GetSearch(req.Keyword /* FIXME: check keyword is not empty */, page)
			pc.setResult2Cache(err, k, result)
		}
	}
	return result, err
}

func (pc *IMacCMSProxyController) request(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	var form maccms.XHRRequest
	if err := ctx.ReadBody(&form); err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	var data repos.MacCMSRepo
	if err := sqls.DB().Model(&repos.MacCMSRepo{}).Where("id = ?", id).First(&data).Error; err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	result, err := pc.realRequest(data, form, id, form.ForceFetch)
	if err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	web.NewData(result).Build(ctx)
}
