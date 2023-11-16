package maccms

import (
	"testing"
)

func TestQs(t *testing.T) {
	qs := NewMacCMSXMLQSBuilder().SetPage(1).SetKeyword("真的出现了")
	realVal := qs.MustString()
	if realVal != `{"pg":"1","wd":"真的出现了"}` {
		t.Fail()
	}
}

func TestQsIDs(t *testing.T) {
	qs := NewMacCMSXMLQSBuilder().SetIDS(1, 2, 3, 4)
	realVal := qs.MustString()
	if realVal != `{"ids":"1,2,3,4"}` {
		t.Fail()
	}
}
