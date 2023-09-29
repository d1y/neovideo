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

func (im *IMacCMSController) batchImport(ctx iris.Context) {
	importData := ctx.FormValueDefault("data", "")
	if len(importData) == 0 {
		web.NewMessage("导入数据为空").SetSuccessWithBool(false).Build(ctx)
		return
	}
	impl.ParseMaccms(importData)
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
