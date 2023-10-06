package handler

import (
	"errors"
	"io"

	"d1y.io/neovideo/controllers/typekit"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/imroc/req/v3"
	"github.com/kataras/iris/v12"
)

type IPagination struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

func newPagination(pg ...int) *IPagination {
	page, limit := 1, 20
	if len(pg) >= 1 {
		page = pg[0]
		if len(pg) >= 2 {
			limit = pg[1]
		}
	}
	p := IPagination{
		Page:  page,
		Limit: limit,
	}
	return &p
}

func NewImportFormWithContext(ctx iris.Context) (*typekit.ImportDataForm, error) {
	var form typekit.ImportDataForm
	if err := ctx.ReadBody(&form); err != nil {
		return nil, err
	}
	return &form, nil
}

func NewImportDataWithContext(ctx iris.Context) (string, error) {
	// TODO: add file support
	idf, err := NewImportFormWithContext(ctx)
	if err != nil {
		return "", err
	}
	url, data := idf.URL, idf.Data
	if len(url) >= 1 { // url 的优先级比 data 高一点
		resp, err := req.Get(url) /* FIXME: verify url */
		if err == nil {
			b, e := io.ReadAll(resp.Body)
			if e == nil {
				data = string(b)
			}
		}
	}
	if len(data) <= 1 {
		return "", errors.New("data is empty")
	}
	return data, nil
}

func NewIDWithContext(ctx iris.Context) (string, bool) {
	id := ctx.Params().Get("id")
	return id, len(id) >= 1
}

func BuildPagination[T any](ctx iris.Context, q ...*gplus.QueryCond[T]) (*gplus.Page[T], error) {
	pg := newPagination()
	if err := ctx.ReadBody(pg); err != nil {
		return nil, err
	}
	var query *gplus.QueryCond[T]
	if len(q) >= 1 {
		query = q[0]
	} else {
		query, _ = gplus.NewQuery[T]()
	}
	page := gplus.NewPage[T](pg.Page, pg.Limit)
	page, gb := gplus.SelectPage(page, query)
	if gb.Error != nil {
		return nil, gb.Error
	}
	return page, nil
}
