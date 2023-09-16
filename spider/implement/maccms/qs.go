package maccms

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

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

func (b *MaccmsQSBuilder) BuildRequest() *req.Request {
	return req.R().SetQueryParams(b.Build())
}

func (b *MaccmsQSBuilder) String() (string, error) {
	val, err := json.Marshal(b.Build())
	if err != nil {
		return "", err
	}
	return string(val), nil
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
