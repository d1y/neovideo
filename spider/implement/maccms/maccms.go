package maccms

import (
	"d1y.io/neovideo/common/json"
)

type IMacCMS struct {
	ResponseType string `json:"response_type,omitempty"`
	ApiURL       string `json:"api_url,omitempty"`
	qs           *MaccmsQSBuilder
}

func New(resType string, api string) *IMacCMS {
	qs := NewMacCMSQSBuilder(resType)
	return &IMacCMS{
		ResponseType: resType,
		ApiURL:       api,
		qs:           qs,
	}
}

// json: [MacCMSReponseTypeJSON]
//
// xml: [MacCMSReponseTypeXML]
func GetResponseType(raw string) string {
	if json.VerifyStringIsJSON(raw) {
		return MacCMSReponseTypeJSON
	}
	return MacCMSReponseTypeXML
}

func (m *IMacCMS) GetHome() (IMacCMSHomeData, error) {
	if m.ResponseType == MacCMSReponseTypeJSON {
		return m.JSONGetHome()
	}
	return m.XMLGetHome()
}

func (m *IMacCMS) GetCategory() ([]IMacCMSCategory, error) {
	if m.ResponseType == MacCMSReponseTypeJSON {
		return m.JSONGetCategory()
	}
	return m.XMLGetCategory()
}

func (m *IMacCMS) GetDetail(id int) (IMacCMSListAttr, []IMacCMSListVideoItem, error) {
	if m.ResponseType == MacCMSReponseTypeJSON {
		return m.JSONGetDetail(id)
	}
	return m.XMLGetDetail(id)
}

func (m *IMacCMS) GetSearch(keyword string, p int) (IMacCMSVideosAndHeader, error) {
	if m.ResponseType == MacCMSReponseTypeJSON {
		return m.JSONGetSearch(keyword, p)
	}
	return m.XMLGetSearch(keyword, p)
}
