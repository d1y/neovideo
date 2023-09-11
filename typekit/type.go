package typekit

import (
	"github.com/kataras/iris/v12"
)

type Registerable interface {
	Register(u iris.Party)
}
