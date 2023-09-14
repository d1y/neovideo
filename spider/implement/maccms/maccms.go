package maccms

import "d1y.io/neovideo/common/json"

type IMacCMS struct {
	ResponseType string `json:"response_type,omitempty"`
	ApiURL       string `json:"api_url,omitempty"`
}

func NewMacCMS(resType string, api string) *IMacCMS {
	return &IMacCMS{
		ResponseType: MacCMSReponseTypeJSON,
		ApiURL:       api,
	}
}

func (m *IMacCMS) GetResponseType(raw string) string {
	if json.VerifyStringIsJSON(raw) {
		return MacCMSReponseTypeJSON
	}
	return MacCMSReponseTypeXML
}
