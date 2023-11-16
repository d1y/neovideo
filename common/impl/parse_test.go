package impl

import (
	"fmt"
	"os"
	"testing"

	"d1y.io/neovideo/spider/implement/maccms"
)

func readFile(f string, ext string) string {
	b, _ := os.ReadFile(fmt.Sprintf("./testdata/%s.%s", f, ext))
	return string(b)
}

func TestParseMaccmsWithLines(t *testing.T) {
	maccmsArray := ParseMaccms(readFile("maccms1", "txt"))
	if len(maccmsArray) != 2 {
		t.FailNow()
	}
	m1, m2 := maccmsArray[0], maccmsArray[1]
	if m1.R18 != true || m1.RespType != maccms.MacCMSReponseTypeXML || m1.Api != "https://x.com" {
		t.FailNow()
	}
	if m2.Name != "低端影视" || m2.Api != "http://v.io/xml" {
		t.FailNow()
	}
}

func TestParseMaccmsWithJSON(t *testing.T) {
	m := ParseMaccms(readFile("maccms_array", "json"))
	if len(m) != 2 {
		t.FailNow()
	}
	m1, m2 := m[0], m[1]
	if m1.Name != "test" || m1.RespType != maccms.MacCMSReponseTypeXML || !m1.R18 || !m1.JiexiParse {
		t.FailNow()
	}
	if m2.Name != "hh" || m2.Api != "https://hh.h" || !m2.R18 {
		t.FailNow()
	}
}

func TestParseJiexiJSON(t *testing.T) {
	jiexiArray := ParseJiexi(readFile("jiexi_array", "json"))
	if len(jiexiArray) != 3 {
		t.Fail()
		return
	}
	if jiexiArray[0].Name != "" || jiexiArray[0].URL != "https://1.io/jiexi_url=" {
		t.Fail()
		return
	}
	if jiexiArray[1].Name != "" || jiexiArray[1].URL != "https://2.io/url=" {
		t.Fail()
		return
	}
	if jiexiArray[2].Name != "test" || jiexiArray[2].URL != "https://d1y.io/jiexi_url=" {
		t.Fail()
	}
}

func TestParseJiexiText(t *testing.T) {
	jiexi := ParseJiexi(readFile("jiexi1", "txt"))
	if len(jiexi) != 7 {
		t.Log("parse len verification failed")
		t.FailNow()
	}
	if jiexi[0].Name != "白嫖线路" || jiexi[0].URL != "https://jiexi.dev/balabala/url=" {
		t.Log("jiexi1 parse [0] fail")
		t.FailNow()
	}
	if jiexi[1].Name != "调试线路" || jiexi[1].URL != "https://patch1.dev/llallalal/url=" {
		t.Log("jiexi1 parse [1] fail")
		t.FailNow()
	}
	for i := 2; i < 7; i++ {
		var item = jiexi[i]
		if len(item.URL) <= 6 {
			t.Fail()
		}
	}
}
