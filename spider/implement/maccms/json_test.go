package maccms

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tidwall/gjson"
)

var cms1 = New(MacCMSReponseTypeJSON, "https://ikunzyapi.com/api.php/provide/vod/at/xml")

func readJSON(file string) []byte {
	f := filepath.Join("./testdata/json", file+".json")
	buf, _ := os.ReadFile(f)
	return buf
}

func readJSON2gjson(file string) *gjson.Result {
	var js = gjson.ParseBytes(readJSON(file))
	return &js
}

func TestJSONGetHome(t *testing.T) {
	a, v, _ := cms1.JsonParseBody(readJSON2gjson("home"))
	if a.Page != 1 || a.PageCount != 2337 || a.RecordCount != 46730 || a.PageSize != 20 {
		t.Fail()
	}
	if len(v) <= 0 {
		t.Fail()
	}
	item := v[0]
	if item.Name != "只是朋友" || item.Id != 49380 || item.Tid != 30 {
		t.Fail()
	}
}

func TestJSONGetCategory(t *testing.T) {
	_, _, c := cms1.JsonParseBody(readJSON2gjson("home"))
	if len(c) != 42 || c[0].Id != 1 || c[0].Text != "电影" {
		t.Fail()
	}
}

func TestJSONGetSearch(t *testing.T) {
	a, v, _ := cms1.JsonParseBody(readJSON2gjson("search"))
	if a.Page != 1 || a.PageCount != 3 || a.RecordCount != 47 || a.PageSize != 20 {
		t.Fail()
	}
	if len(v) <= 0 {
		t.Fail()
	}
	item := v[0]
	if item.Name != "疾速营救" || item.Id != 49766 || item.Tid != 6 {
		t.Fail()
	}
}

func TestJSONGetDetail(t *testing.T) {
	_, v, _ := cms1.JsonParseBody(readJSON2gjson("detail"))
	if len(v) != 1 || v[0].Id != 48691 || v[0].Name != "惊天营救2" || v[0].Tid != 6 {
		t.Fail()
	}
}
