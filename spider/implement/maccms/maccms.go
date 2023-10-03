package maccms

import (
	"d1y.io/neovideo/common/json"
	"d1y.io/neovideo/models/repos"
)

type IMacCMS struct {
	ResponseType string `json:"response_type"`
	ApiURL       string `json:"api_url"`
	qs           *MaccmsQSBuilder
}

type respType string

func (rt respType) IsJSON() bool {
	return rt == MacCMSReponseTypeJSON
}
func (rt respType) IsXML() bool {
	return rt == MacCMSReponseTypeXML
}

func NewWithApi(api string) *IMacCMS {
	return New("", api)
}

func New(resType string, api string) *IMacCMS {
	qs := NewMacCMSQSBuilder(resType)
	return &IMacCMS{
		ResponseType: resType,
		ApiURL:       api,
		qs:           qs,
	}
}

func (m *IMacCMS) SetReponseType(t string) {
	m.ResponseType = t
}

func (m *IMacCMS) SetJSONResponseType() {
	m.SetReponseType(MacCMSReponseTypeJSON)
}

func (m *IMacCMS) SetXMLReponseType() {
	m.SetReponseType(MacCMSReponseTypeXML)
}

func GetResponseTypeWithByte(b []byte) respType {
	return GetResponseType(string(b))
}

// json: [MacCMSReponseTypeJSON]
//
// xml: [MacCMSReponseTypeXML]
func GetResponseType(raw string) respType {
	if json.VerifyStringIsJSON(raw) {
		return MacCMSReponseTypeJSON
	}
	return MacCMSReponseTypeXML
}

func (m *IMacCMS) Request(xhr XHRRequest) {
}

func (m *IMacCMS) GetHome(page int, tid ...int) (IMacCMSHomeData, error) {
	if m.ResponseType == MacCMSReponseTypeJSON {
		return m.JSONGetHome(page, tid...)
	}
	return m.XMLGetHome(page, tid...)
}

func (m *IMacCMS) GetCategory() ([]repos.IMacCMSCategory, error) {
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
