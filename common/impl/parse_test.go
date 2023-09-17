package impl_test

import (
	"fmt"
	"os"
	"testing"

	"d1y.io/neovideo/common/impl"
)

func readJiexi(f string) string {
	b, _ := os.ReadFile(fmt.Sprintf("./testdata/%s.txt", f))
	return string(b)
}

func TestParseJiexiText(t *testing.T) {
	jiexi := impl.ParseJiexi(readJiexi("jiexi1"))
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
