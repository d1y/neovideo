package other

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"d1y.io/neovideo/spider/implement/maccms"
)

const (
	MacCMSListAction   = "list"
	MacCMSDetailAction = "detail"
)

type maccmsQSBuilder struct {
	buildType string
	m         map[string]any
}

func (b *maccmsQSBuilder) SetKeyword(keyword string) *maccmsQSBuilder {
	b.m["wd"] = keyword
	return b
}

func (b *maccmsQSBuilder) SetPage(page int) *maccmsQSBuilder {
	b.m["pg"] = page
	return b
}

func (b *maccmsQSBuilder) SetAction(ac string) *maccmsQSBuilder {
	b.m["ac"] = ac
	return b
}

func (b *maccmsQSBuilder) SetListAction() *maccmsQSBuilder {
	b.SetAction(MacCMSListAction)
	return b
}

func (b *maccmsQSBuilder) SetDetailAction() *maccmsQSBuilder {
	b.SetAction(MacCMSDetailAction)
	return b
}

func (b *maccmsQSBuilder) SetCategory(id int) *maccmsQSBuilder {
	b.m["t"] = id
	return b
}

func (b *maccmsQSBuilder) SetHWithTime(h int) *maccmsQSBuilder {
	b.m["h"] = h // 几小时内的数据
	return b
}

func (b *maccmsQSBuilder) SetIDS(ids ...int) *maccmsQSBuilder {
	strSlice := make([]string, len(ids))
	for i, id := range ids {
		strSlice[i] = strconv.Itoa(id)
	}
	realIDs := strings.Join(strSlice, ",")
	b.m["ids"] = realIDs
	return b
}

func (b *maccmsQSBuilder) Build() map[string]string {
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

func (b *maccmsQSBuilder) String() (string, error) {
	val, err := json.Marshal(b.Build())
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func NewMacCMSQSBuilder(bType string) *maccmsQSBuilder {
	m := make(map[string]any)
	return &maccmsQSBuilder{
		buildType: bType,
		m:         m,
	}
}

func NewMacCMSXMLQSBuilder() *maccmsQSBuilder {
	return NewMacCMSQSBuilder(maccms.MacCMSReponseTypeXML)
}

func NewMacCMSJSONQSBuilder() *maccmsQSBuilder {
	return NewMacCMSQSBuilder(maccms.MacCMSReponseTypeJSON)
}
