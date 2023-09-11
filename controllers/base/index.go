package base

import (
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
)

type BaseController struct {
	iris.Context
}

func (c *BaseController) ping(ctx iris.Context) {
	var t = time.Now().Unix()
	ctx.Text(fmt.Sprintf("pong %d", t))
}

func Register(u iris.Party) {
	var bc BaseController
	u.Get("/ping", bc.ping)
}
