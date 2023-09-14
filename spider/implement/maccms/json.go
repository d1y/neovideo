package maccms

import (
	"errors"
	"io"
	"time"

	typekkkit "d1y.io/neovideo/common/typekit"
	"github.com/imroc/req/v3"
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

func (m *IMacCMS) jsonParseBody(result gjson.Result) (IMacCMSListAttr, []IMacCMSListVideoItem, []IMacCMSCategory) {

	var attr IMacCMSListAttr
	ints := typekkkit.Int64Slice2Int(result.Get("pagecount").Int(), result.Get("page").Int(), result.Get("total").Int())
	attr.PageCount = ints[0]
	attr.Page = ints[1]
	attr.RecordCount = ints[2]

	_class := result.Get("class").Array()
	_list := result.Get("list").Array()

	var category = make([]IMacCMSCategory, len(_class))
	var videos = make([]IMacCMSListVideoItem, len(_list))

	for _, item := range _list {
		typeID := m.gjsonResult2Int(item, "type_id")
		id := m.gjsonResult2Int(item, "id")
		t := m.gjsonResult2Time(item, "vod_time")
		name := m.gjsonResult2Str(item, "vod_name")
		videos = append(videos, IMacCMSListVideoItem{
			Last: t,
			Id:   id,
			Tid:  typeID,
			Name: name,
		})
	}
	for _, item := range _class {
		id := m.gjsonResult2Int(item, "type_id")
		name := m.gjsonResult2Str(item, "type_name")
		category = append(category, IMacCMSCategory{
			Text: name,
			Id:   id,
		})
	}
	return attr, videos, category
}

func (m *IMacCMS) JSONGetHome() (IMacCMSVideosAndHeader, error) {
	res, err := req.Get(m.ApiURL)
	if err != nil {
		return IMacCMSVideosAndHeader{}, err
	}
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return IMacCMSVideosAndHeader{}, err
	}
	if !gjson.ValidBytes(buf) {
		return IMacCMSVideosAndHeader{}, errors.New("invalid json")
	}
	result := gjson.ParseBytes(buf)
	attr, videos, _ := m.jsonParseBody(result)
	return IMacCMSVideosAndHeader{
		ListHeader: attr,
		Videos:     videos,
	}, nil
}

func (m *IMacCMS) JSONGetCategory() ([]IMacCMSCategory, error) {
	res, err := req.Get(m.ApiURL)
	if err != nil {
		return []IMacCMSCategory{}, err
	}
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return []IMacCMSCategory{}, err
	}
	if !gjson.ValidBytes(buf) {
		return []IMacCMSCategory{}, errors.New("invalid json")
	}
	result := gjson.ParseBytes(buf)
	_, _, category := m.jsonParseBody(result)
	return category, nil
}

func (m *IMacCMS) JSONGetSearch(keyword string, page int) {

}

func (m *IMacCMS) JSONGetDetail(id int) {

}
