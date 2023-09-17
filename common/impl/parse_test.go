package impl_test

import (
	"fmt"
	"os"
	"testing"

	"d1y.io/neovideo/common/impl"
)

func readJiexi(f string, ext string) string {
	b, _ := os.ReadFile(fmt.Sprintf("./testdata/%s.%s", f, ext))
	return string(b)
}

func TestParseJiexiJSON(t *testing.T) {
	jiexiArray := impl.ParseJiexi(readJiexi("jiexi_array", "json"))
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
	jiexi := impl.ParseJiexi(readJiexi("jiexi1", "txt"))
	if len(jiexi) != 5 {
		t.Log("parse len verification failed")
		t.Fail()
		return
	}
	if jiexi[0].Name != "白嫖线路" || jiexi[0].URL != "https://jiexi.dev/balabala/url=" {
		t.Log("jiexi1 parse [0] fail")
		t.Fail()
		return
	}
	if jiexi[1].Name != "调试线路" || jiexi[1].URL != "https://patch1.dev/llallalal/url=" {
		t.Log("jiexi1 parse [1] fail")
		t.Fail()
		return
	}
	for i := 2; i < 5; i++ {
		var item = jiexi[i]
		if len(item.URL) <= 6 {
			t.Fail()
		}
	}
}
