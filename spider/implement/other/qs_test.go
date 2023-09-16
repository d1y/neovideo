package other_test

import (
	"testing"

	"d1y.io/neovideo/spider/implement/other"
)

func TestQs(t *testing.T) {
	qs := other.NewMacCMSXMLQSBuilder().SetPage(1).SetKeyword("真的出现了")
	realVal, _ := qs.String()
	if realVal != `{"pg":"1","wd":"真的出现了"}` {
		t.Fail()
	}
}
