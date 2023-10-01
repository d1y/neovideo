package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type routeMeta struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type siteInfo struct {
	Route []routeMeta `json:"route_meta"`
}

func Siteinfo(p iris.Context, meta []context.RouteReadOnly) {
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
