package maccms

import (
	"d1y.io/neovideo/common/json"
)

type IMacCMS struct {
	ResponseType string `json:"response_type,omitempty"`
	ApiURL       string `json:"api_url,omitempty"`
	qs           *MaccmsQSBuilder
}

func NewMacCMS(resType string, api string) *IMacCMS {
	qs := NewMacCMSQSBuilder(resType)
	return &IMacCMS{
		ResponseType: MacCMSReponseTypeJSON,
		ApiURL:       api,
		qs:           qs,
	}
}

func (m *IMacCMS) GetResponseType(raw string) string {
	if json.VerifyStringIsJSON(raw) {
		return MacCMSReponseTypeJSON
	}
	return MacCMSReponseTypeXML
}
