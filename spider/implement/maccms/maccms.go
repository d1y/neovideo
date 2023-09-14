package maccms

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
