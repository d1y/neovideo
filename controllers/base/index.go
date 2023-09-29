package base

import (
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type routeMeta struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

type siteInfo struct {
	Route []routeMeta `json:"route_meta,omitempty"`
}

type BaseController struct {
	iris.Context
}

func (c *BaseController) ping(ctx iris.Context) {
	var t = time.Now().Unix()
	ctx.Text(fmt.Sprintf("pong %d", t))
}

func Register(u iris.Party) {
	var bc BaseController
	u.Get("/ping", bc.ping).Name = "调试"
}

func RenderSiteInfo(p iris.Context, meta []context.RouteReadOnly) {
	var result = make([]routeMeta, len(meta))
	for idx, item := range meta {
		result[idx] = routeMeta{
			Name: item.Name(),
			Path: item.Path(),
		}
	}
	p.JSON(siteInfo{
		Route: result,
	})
}
