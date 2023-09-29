package maccms

import (
	"testing"
)

func TestParseDDRawURL(t *testing.T) {
	var d1 = "第1集$https://c3.monidai.com/20230922/SXJpJgy8/index.m3u8#第2集$https://c3.monidai.com/20230922/ldO8mJfc/index.m3u8#第3集$https://c3.monidai.com/20230922/gFua4Ofu/index.m3u8#第4集$https://c3.monidai.com/20230922/dz205WN2/index.m3u8#第5集$https://c3.monidai.com/20230922/nQSCz9lb/index.m3u8#第6集$https://c3.monidai.com/20230922/KFq5BiNv/index.m3u8#第7集$https://c3.monidai.com/20230922/NfXZ0nfP/index.m3u8#第8集$https://c3.monidai.com/20230922/W5D0YBSt/index.m3u8#第9集$https://c3.monidai.com/20230922/IDjFlncW/index.m3u8#第10集$https://c3.monidai.com/20230922/e7yVirZN/index.m3u8#第11集$https://c3.monidai.com/20230923/77cneAbE/index.m3u8#第12集$https://c3.monidai.com/20230923/JUOOtx90/index.m3u8#第13集$https://c3.monidai.com/20230924/e1Fgw4CV/index.m3u8#第14集$https://c3.monidai.com/20230924/9w2s6xuz/index.m3u8#第15集$https://c3.monidai.com/20230925/DIUA4kYX/index.m3u8#第16集$https://c3.monidai.com/20230926/zE2CfhdC/index.m3u8#第17集$https://c3.monidai.com/20230927/5xgSa9VG/index.m3u8#第18集$https://c3.monidai.com/20230928/08MOp9qU/index.m3u8#第19集$https://c3.monidai.com/20230929/2MkW5H9J/index.m3u8"
	var dds = parseDDRawURL(d1)
	if len(dds) != 19 {
		t.FailNow()
	}
}
