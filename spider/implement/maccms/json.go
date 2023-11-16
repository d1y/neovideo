package maccms

import (
	"errors"
	"time"

	typekkkit "d1y.io/neovideo/common/typekit"
	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/spider/axios"
	"github.com/tidwall/gjson"
)

func (m *IMacCMS) gjsonResult2Int(r gjson.Result, key string) int {
	return typekkkit.Int642Int(r.Get(key).Int())
}

func (m *IMacCMS) gjsonResult2Str(r gjson.Result, key string) string {
	return r.Get(key).String()
}

func (m *IMacCMS) gjsonResult2Time(r gjson.Result, key string) time.Time {
	t, _ := time.Parse(time.DateTime, r.Get(key).String())
	return t
}

func (m *IMacCMS) JsonParseBody(result *gjson.Result) (IMacCMSListAttr, []IMacCMSListVideoItem, []repos.IMacCMSCategory) {

	// NOTE(d1y): `result` 传递过来之前不能是 nil
	if result == nil {
		panic("maccms json parse body result is nil")
	}

	var attr = IMacCMSListAttr{}
	ints := typekkkit.Int64Slice2Int(result.Get("pagecount").Int(), result.Get("page").Int(), result.Get("total").Int(), result.Get("limit").Int())
	attr.PageCount = ints[0]
	attr.Page = ints[1]
	attr.RecordCount = ints[2]
	attr.PageSize = ints[3]

	_class := result.Get("class").Array()
	_list := result.Get("list").Array()

	var category = make([]repos.IMacCMSCategory, len(_class))
	var videos = make([]IMacCMSListVideoItem, len(_list))

	for idx, item := range _list {
		typeID := m.gjsonResult2Int(item, "type_id")
		id := m.gjsonResult2Int(item, "vod_id")
		if id == 0 {
			id = m.gjsonResult2Int(item, "id")
		}
		t := m.gjsonResult2Time(item, "vod_time")
		name := m.gjsonResult2Str(item, "vod_name")
		videos[idx] = IMacCMSListVideoItem{
			Last: t,
			Id:   id,
			Tid:  typeID,
			Name: name,
		}
	}
	for idx, item := range _class {
		id := m.gjsonResult2Int(item, "type_id")
		name := m.gjsonResult2Str(item, "type_name")
		category[idx] = repos.IMacCMSCategory{
			Text: name,
			Id:   id,
		}
	}
	return attr, videos, category
}

func (m *IMacCMS) byte2gjson(buf []byte) (*gjson.Result, error) {
	if !gjson.ValidBytes(buf) {
		return nil, errors.New("invalid json")
	}
	js := gjson.ParseBytes(buf)
	return &js, nil
}

func (m *IMacCMS) JSONGetHome(page int, tid ...int) (*IMacCMSHomeData, error) {
	res, err := m.qs.SetHome(page, tid).BuildGetRequest(m.ApiURL)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	result, err := m.byte2gjson(res)
	if err != nil {
		return nil, err
	}
	attr, videos, category := m.JsonParseBody(result)
	return &IMacCMSHomeData{
		ListHeader: attr,
		Videos:     videos,
		Category:   category,
	}, nil
}

func (m *IMacCMS) JSONGetCategory() ([]repos.IMacCMSCategory, error) {
	res, err := axios.Get(m.ApiURL)
	if err != nil {
		return nil, err
	}
	result, err := m.byte2gjson(res)
	if err != nil {
		return nil, err
	}
	_, _, category := m.JsonParseBody(result)
	return category, nil
}

func (m *IMacCMS) JSONGetSearch(keyword string, page int) (*IMacCMSVideosAndHeader, error) {
	res, err := m.qs.SetKeyword(keyword).SetPage(page).BuildGetRequest(m.ApiURL)
	if err != nil {
		return nil, err
	}
	result, err := m.byte2gjson(res)
	if err != nil {
		return nil, err
	}
	attr, videos, _ := m.JsonParseBody(result)
	return &IMacCMSVideosAndHeader{
		ListHeader: attr,
		Videos:     videos,
	}, nil
}

func (m *IMacCMS) JSONGetDetail(ids ...int) (*IMacCMSListAttr, []IMacCMSListVideoItem, error) {
	res, err := m.qs.SetDetailAction().SetIDS(ids...).BuildGetRequest(m.ApiURL)
	if err != nil {
		return nil, nil, err
	}
	result, err := m.byte2gjson(res)
	if err != nil {
		return nil, nil, err
	}
	attr, videos, _ := m.JsonParseBody(result)
	return &attr, videos, nil
}
