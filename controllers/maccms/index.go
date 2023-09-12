package maccms

import (
	"d1y.io/neovideo/models/repos"

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

func Register(u iris.Party) {
	var imc IMacCMSController
	u.Get("/", imc.getList)
}
