package maccms_test

import (
	"testing"

	"d1y.io/neovideo/spider/implement/maccms"
)

func TestQs(t *testing.T) {
	qs := maccms.NewMacCMSXMLQSBuilder().SetPage(1).SetKeyword("真的出现了")
	realVal, _ := qs.String()
	if realVal != `{"pg":"1","wd":"真的出现了"}` {
		t.Fail()
	}
}

func TestQsIDs(t *testing.T) {
	qs := maccms.NewMacCMSXMLQSBuilder().SetIDS(1, 2, 3, 4)
	realVal, _ := qs.String()
	if realVal != `{"ids":"1,2,3,4"}` {
		t.Fail()
	}
}
