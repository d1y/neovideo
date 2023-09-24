package maccms_test

import (
	"testing"

	"d1y.io/neovideo/spider/implement/maccms"
	"github.com/stretchr/testify/assert"
)

var u = "https://ikunzyapi.com/api.php/provide/vod/at/xml"

func TestMaccms(t *testing.T) {
	if !assert.Equal(t, maccms.GetResponseType(u), maccms.MacCMSReponseTypeXML, "get response type error") {
		t.FailNow()
	}
	var cms1 = maccms.NewMacCMS(maccms.MacCMSReponseTypeXML, u)
	_, err := cms1.GetHome()
	if err != nil {
		t.FailNow()
	}
}
