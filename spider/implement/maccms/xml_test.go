package maccms_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/spider/implement/maccms"
	"github.com/beevik/etree"
)

var cms = maccms.New(maccms.MacCMSReponseTypeXML, "https://ikunzyapi.com/api.php/provide/vod/at/xml")

func readXML(file string) []byte {
	f := filepath.Join("./testdata/xml", file+".xml")
	buf, _ := os.ReadFile(f)
	return buf
}

func readXML2ETree(file string) *etree.Element {
	homeBuf := readXML(file)
	homeXML := etree.NewDocument()
	homeXML.ReadFromBytes(homeBuf)
	return homeXML.Root()
}

func findIdByCategory(raw []repos.IMacCMSCategory, id int) (repos.IMacCMSCategory, error) {
	for i := 0; i < len(raw); i++ {
		if raw[i].Id == id {
			return raw[i], nil
		}
	}
	return repos.IMacCMSCategory{}, errors.New("not found")
}

func TestXMLHomeResponse(t *testing.T) {
	homeData, err := cms.XMLGetHomeWithEtreeRoot(readXML2ETree("home"))
	if err != nil {
		t.Fail()
	}
	if len(homeData.Videos) != 20 {
		t.Fail()
	} else {
		item, err := findIdByCategory(homeData.Category, homeData.Videos[0].Tid)
		if err != nil {
			t.Fail()
		}
		if item.Text != "大陆综艺" {
			t.Fail()
		}
	}
	if len(homeData.Category) != 39 && homeData.Category[0].Text != "电影" {
		t.Fail()
	}
	lh := homeData.ListHeader
	if lh.RecordCount != 46729 || lh.PageCount != 2337 || lh.Page != 1 || lh.PageSize != 20 {
		t.Fail()
	}
}

func TestXMLSearchResponse(t *testing.T) {
	searchData, err := cms.XMLGetSearchWithEtreeRoot(readXML2ETree("search"))
	if err != nil {
		t.Fail()
	}
	lh := searchData.ListHeader
	if lh.RecordCount != 5 || lh.PageCount != 1 || lh.Page != 1 || lh.PageSize != 20 {
		t.Fail()
	}
	v := searchData.Videos
	if len(v) != 5 || v[len(v)-1].Id != 6876 {
		t.Fail()
	}
}

func TestXMLDetailResponse(t *testing.T) {
	_, v, e := cms.XMLGetDetailWithEtreeRoot(readXML2ETree("detail"))
	if e != nil {
		t.Fail()
	}
	if len(v) != 1 || len(v[0].DD) != 1 || v[0].DD[0].Flag != "ikm3u8" {
		t.Fail()
	}
}
