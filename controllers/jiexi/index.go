package jiexi

import (
	"fmt"
	"io"

	"d1y.io/neovideo/common/impl"
	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/models/web"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/imroc/req/v3"
	"github.com/kataras/iris/v12"
)

type JiexiController struct {
}

func (jx *JiexiController) getList(ctx iris.Context) {
	jiexi, db := gplus.SelectList[repos.JiexiRepo](nil)
	if db.Error != nil {
		web.NewError(db.Error).Build(ctx)
		return
	}
	web.NewData(jiexi).Build(ctx)
}

func (jx *JiexiController) create(ctx iris.Context) {
	var jiexiForm repos.IJiexi
	if err := ctx.ReadBody(&jiexiForm); err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	var insertData = repos.JiexiRepo{
		IJiexi: jiexiForm,
	}
	err := gplus.Insert[repos.JiexiRepo](&insertData).Error
	if err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	web.NewData(insertData)
}

func (jx *JiexiController) delete(ctx iris.Context) {
	id := ctx.Params().Get("id")
	err := gplus.DeleteById[repos.JiexiRepo](id).Error
	if err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	web.NewData(id).SetMessage("删除成功").Build(ctx)
}

func (jx *JiexiController) list2map(raw []*repos.JiexiRepo) map[string]*repos.JiexiRepo {
	var m = make(map[string]*repos.JiexiRepo)
	for _, item := range raw {
		if len(item.URL) >= 1 {
			m[item.URL] = item
		}
	}
	return m
}

func (jx *JiexiController) batchImport(ctx iris.Context) {
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
	list := impl.ParseJiexi(importData)
	if len(list) == 0 {
		web.NewMessage("导入数据为空").SetSuccessWithBool(false).Build(ctx)
		return
	}
	jiexiRepos, g := gplus.SelectList[repos.JiexiRepo](nil)
	if g.Error != nil {
		web.NewError(g.Error).Build(ctx)
		return
	}
	mp := jx.list2map(jiexiRepos)
	var importJiexi = make([]*repos.JiexiRepo, 0)
	for _, item := range list {
		if _, ok := mp[item.URL]; ok {
			continue
		}
		importJiexi = append(importJiexi, &repos.JiexiRepo{
			IJiexi: repos.IJiexi{
				Name: item.Name,
				URL:  item.URL,
			},
		})
	}
	if len(importJiexi) <= 0 {
		web.NewMessage("导入数据为空").SetSuccessWithBool(true).SetData(0).Build(ctx)
		return
	}
	err := gplus.InsertBatch(importJiexi).Error
	if err != nil {
		web.NewError(err).Build(ctx)
		return
	}
	importLen := len(importJiexi)
	web.NewMessage(fmt.Sprintf("新增成功(%d条)", importLen)).SetData(importLen).SetSuccessWithBool(true).Build(ctx)
}

func Register(u iris.Party) {
	var jx JiexiController
	u.Get("/", jx.getList).Name = "获取解析列表"
	u.Post("/", jx.create).Name = "创建解析"
	u.Delete("/{id:int}", jx.delete).Name = "删除解析"
	u.Post("/batch_import", jx.batchImport).Name = "批量导入解析列表"
}
