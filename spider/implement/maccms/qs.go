package maccms

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"d1y.io/neovideo/spider/axios"
	"github.com/imroc/req/v3"
	"golang.org/x/exp/slices"
)

const (
	MacCMSListAction   = "list"
	MacCMSDetailAction = "detail"
)

var (
	// Deprecated: remove this
	ignoreQueryGen = []string{
		"ids",
	}
)

type MaccmsQSBuilder struct {
	buildType string
	m         map[string]any
}

func (b *MaccmsQSBuilder) SetKeyword(keyword string) *MaccmsQSBuilder {
	b.m["wd"] = keyword
	return b
}

func (b *MaccmsQSBuilder) SetHome(page int, tid []int) *MaccmsQSBuilder {
	b.SetPage(page)
	if len(tid) >= 1 && tid[0] >= 1 {
		b.SetCategory(tid[0])
	}
	return b
}

func (b *MaccmsQSBuilder) SetPage(page int) *MaccmsQSBuilder {
	b.m["pg"] = page
	return b
}

func (b *MaccmsQSBuilder) SetAction(ac string) *MaccmsQSBuilder {
	b.m["ac"] = ac
	return b
}

func (b *MaccmsQSBuilder) SetListAction() *MaccmsQSBuilder {
	b.SetAction(MacCMSListAction)
	return b
}

func (b *MaccmsQSBuilder) SetDetailAction() *MaccmsQSBuilder {
	b.SetAction(MacCMSDetailAction)
	return b
}

func (b *MaccmsQSBuilder) SetCategory(id int) *MaccmsQSBuilder {
	b.m["t"] = id
	return b
}

func (b *MaccmsQSBuilder) SetHWithTime(h int) *MaccmsQSBuilder {
	b.m["h"] = h // 几小时内的数据
	return b
}

func (b *MaccmsQSBuilder) SetIDS(ids ...int) *MaccmsQSBuilder {
	strSlice := make([]string, len(ids))
	for i, id := range ids {
		strSlice[i] = strconv.Itoa(id)
	}
	realIDs := strings.Join(strSlice, ",")
	b.m["ids"] = realIDs
	return b
}

// FIXME: querystring 全部自己生成, 而不是使用 req 里 net/url 标准库生成的自动转码方式
//
// Deprecated: 不稳定方法..
func (b *MaccmsQSBuilder) Build(full ...bool) (map[string]string, map[string]string, bool) {
	var result = make(map[string]string)
	var customMap = make(map[string]string)
	var isFull = false
	if len(full) == 1 {
		isFull = true
	}
	for k, v := range b.m {
		var onceNeedCustom = slices.Contains(ignoreQueryGen, k) && !isFull
		if strV, ok := v.(string); ok {
			if onceNeedCustom {
				customMap[k] = strV
			} else {
				result[k] = strV
			}
		} else {
			if reflect.TypeOf(v).Kind() == reflect.Int {
				s := strconv.Itoa(v.(int))
				if onceNeedCustom {
					customMap[k] = s
				} else {
					result[k] = s
				}
			}
		}
	}
	return result, customMap, len(customMap) >= 1
}

func (b *MaccmsQSBuilder) wrapperRequestHeader(r *req.Request) *req.Request {
	r.SetHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	return r
}

// 该方法没有包装缓存机制, 不要使用
// Deprecated: use BuildGetRequest/BuildPostRequest
func (b *MaccmsQSBuilder) BuildRequest() *req.Request {
	sq, _, _ := b.Build(true)
	return b.wrapperRequestHeader(axios.Request().SetQueryParams(sq))
}

func (b *MaccmsQSBuilder) WrapperRealURL(api string) (string, map[string]string) {
	qs, unqs, ok := b.Build()
	var realApi = api
	if ok {
		// TODO: 或许应该使用标准的标准的 `net/url` 包生成 querystring
		var s = "?"
		for k, v := range unqs {
			s += fmt.Sprintf("%s=%s&", k, v)
		}
		s = s[:len(s)-1]
		realApi += s
	}
	return realApi, qs
}

func (b *MaccmsQSBuilder) BuildGetRequest(api string) ([]byte, error) {
	realApi, qs := b.WrapperRealURL(api)
	return axios.Get(realApi, qs)
}

func (b *MaccmsQSBuilder) BuildPostRequest(api string) ([]byte, error) {
	realApi, qs := b.WrapperRealURL(api)
	return axios.Post(realApi, qs)
}

func (b *MaccmsQSBuilder) String() (string, error) {
	q, _, _ := b.Build(true)
	val, err := json.Marshal(q)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (b *MaccmsQSBuilder) MustString() (s string) {
	s, _ = b.String()
	return
}

func NewMacCMSQSBuilder(bType string) *MaccmsQSBuilder {
	m := make(map[string]any)
	return &MaccmsQSBuilder{
		buildType: bType,
		m:         m,
	}
}

func NewMacCMSXMLQSBuilder() *MaccmsQSBuilder {
	return NewMacCMSQSBuilder(MacCMSReponseTypeXML)
}

func NewMacCMSJSONQSBuilder() *MaccmsQSBuilder {
	return NewMacCMSQSBuilder(MacCMSReponseTypeJSON)
}
