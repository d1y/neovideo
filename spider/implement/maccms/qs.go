package maccms

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"d1y.io/neovideo/spider/axios"
	"github.com/imroc/req/v3"
)

const (
	MacCMSListAction   = "list"
	MacCMSDetailAction = "detail"
)

type MaccmsQSBuilder struct {
	buildType string
	m         map[string]any
}

func (b *MaccmsQSBuilder) SetKeyword(keyword string) *MaccmsQSBuilder {
	b.m["wd"] = keyword
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

func (b *MaccmsQSBuilder) Build() map[string]string {
	var result = make(map[string]string)
	for k, v := range b.m {
		if strV, ok := v.(string); ok {
			result[k] = strV
		} else {
			if reflect.TypeOf(v).Kind() == reflect.Int {
				result[k] = strconv.Itoa(v.(int))
			}
		}
	}
	return result
}

func (b *MaccmsQSBuilder) wrapperRequestHeader(r *req.Request) *req.Request {
	r.SetHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	return r
}

// 该方法没有包装缓存机制, 不要使用
// Deprecated: use BuildGetRequest/BuildPostRequest
func (b *MaccmsQSBuilder) BuildRequest() *req.Request {
	return b.wrapperRequestHeader(axios.Request().SetQueryParams(b.Build()))
}

func (b *MaccmsQSBuilder) BuildGetRequest(api string) ([]byte, error) {
	return axios.Get(api, b.Build())
}

func (b *MaccmsQSBuilder) BuildPostRequest(api string) ([]byte, error) {
	return axios.Post(api, b.Build())
}

func (b *MaccmsQSBuilder) String() (string, error) {
	val, err := json.Marshal(b.Build())
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (b *MaccmsQSBuilder) MustString() string {
	val, err := json.Marshal(b.Build())
	if err != nil {
		return ""
	}
	return string(val)
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
