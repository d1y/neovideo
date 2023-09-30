package maccms

import (
	"d1y.io/neovideo/common/impl"
	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/models/web"

	"github.com/acmestack/gorm-plus/gplus"
	"github.com/kataras/iris/v12"
)

type IMacCMSController struct {
}

func (im *IMacCMSController) getList(ctx iris.Context) {
	qs := ctx.Request().URL.Query()
	list, _ := gplus.SelectList[repos.MacCMSRepo](gplus.BuildQuery[repos.MacCMSRepo](qs))
	ctx.JSON(list)
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
	importData := ctx.FormValueDefault("data", "")
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
	for _, item := range impl.ParseMaccms(importData) {
		if _, ok := m[item.Api]; ok {
			continue
		}
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
		web.NewMessage("保存到数据库失败").SetSuccessWithBool(false).Build(ctx)
		return
	}
	web.NewData(l).Build(ctx)
}

func Register(u iris.Party) {
	var imc IMacCMSController
	var px IMacCMSProxyController
	u.Get("/", imc.getList).Name = "获取苹果CMS列表"
	u.Post("/", imc.create).Name = "创建苹果CMS"
	u.Post("/batch_import", imc.batchImport).Name = "批量导入苹果CMS列表"
	u.Delete("/{id:int}", imc.delete).Name = "删除苹果CMS"
	u.PartyFunc("/proxy", px.Register)
}
