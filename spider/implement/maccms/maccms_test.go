package maccms

import (
	"testing"
)

var u = "https://ikunzyapi.com/api.php/provide/vod/at/xml"

func TestMaccms(t *testing.T) {
	var testCMS = New(MacCMSReponseTypeXML, u)
	_, err := testCMS.GetHome(1)
	if err != nil {
		t.FailNow()
	}
}
