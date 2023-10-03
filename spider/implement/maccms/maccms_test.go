package maccms_test

import (
	"testing"

	"d1y.io/neovideo/spider/implement/maccms"
)

var u = "https://ikunzyapi.com/api.php/provide/vod/at/xml"

func TestMaccms(t *testing.T) {
	var cms1 = maccms.New(maccms.MacCMSReponseTypeXML, u)
	_, err := cms1.GetHome(1)
	if err != nil {
		t.FailNow()
	}
}
