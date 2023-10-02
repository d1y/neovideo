package handler

import (
	"errors"
	"io"

	"d1y.io/neovideo/controllers/typekit"
	"github.com/imroc/req/v3"
	"github.com/kataras/iris/v12"
)

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
