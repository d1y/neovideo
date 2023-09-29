package handler

import (
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
)

func Ping(ctx iris.Context) {
	var t = time.Now().Unix()
	ctx.Text(fmt.Sprintf("pong %d", t))
}
