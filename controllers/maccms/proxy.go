package maccms

import (
	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/models/web"
	"d1y.io/neovideo/spider/implement/maccms"
	"d1y.io/neovideo/sqls"
	"github.com/kataras/iris/v12"
)

type requestForm struct {
	ResponseType string `json:"r_type" form:"r_type"`
	Page         int    `json:"page" form:"page"`
	Keyword      string `json:"keyword" form:"keyword"`
	Action       string `json:"action" form:"action"`
	Category     string `json:"category" form:"category"`
	Hour         int    `json:"hour" form:"hour"`
	Ids          []int  `json:"ids" form:"ids"`
}

type IMacCMSProxyController struct {
	// cc *cache.Cache
}

func (pc *IMacCMSProxyController) Register(u iris.Party) {
	// pc.cc = cache.New(6*time.Second, 24*time.Second)
	u.Post("/{id:int}", pc.request)
}

// func (pc *IMacCMSProxyController) getCacheWithID(id int) (repos.MacCMSRepo, bool) {
// 	k := strconv.Itoa(id)
// 	iface, ok := pc.cc.Get(k)
// 	if !ok {
// 		return repos.MacCMSRepo{}, false
// 	}
// 	v, ok := iface.(repos.MacCMSRepo)
// 	if ok {
// 		return v, true
// 	}
// 	return repos.MacCMSRepo{}, false
// }

// func (pc *IMacCMSProxyController) setCacheWithID(id int, data repos.MacCMSRepo) {
// 	k := strconv.Itoa(id)
// 	pc.cc.SetDefault(k, data)
// }

func (pc *IMacCMSProxyController) request(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	var form requestForm
	if err := ctx.ReadBody(&form); err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	var data repos.MacCMSRepo
	if err := sqls.DB().Model(&repos.MacCMSRepo{}).Where("id = ?", id).First(&data).Error; err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	cms := maccms.New(data.RespType, data.Api)
	h, e := cms.GetHome()
	if err != nil {
		web.NewError(e)
		return
	}
	web.NewData(h)
}
