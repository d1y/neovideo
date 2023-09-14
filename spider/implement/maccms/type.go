package maccms

import "time"

const (
	MacCMSReponseTypeXML  = "xml"
	MacCMSReponseTypeJSON = "Json"
)

type IMacCMSCategory struct {
	Text string `json:"text,omitempty"`
	Id   int    `json:"id,omitempty"`
}

type IMacCMSListAttr struct {
	Page        int `json:"page,omitempty"`
	PageCount   int `json:"page_count,omitempty"`
	PageSize    int `json:"page_size,omitempty"`
	RecordCount int `json:"record_count,omitempty"`
}

type IMacCMSVideoDDTag struct {
	Flag   string `json:"flag,omitempty"`
	RawURL string `json:"raw_url,omitempty"`
}

type IMacCMSListVideoItem struct {
	Last     time.Time `json:"last,omitempty"`
	Id       int       `json:"id,omitempty"`
	Tid      int       `json:"tid,omitempty"`
	Name     string    `json:"name,omitempty"`
	Type     string    `json:"type,omitempty"`
	Dt       string    `json:"dt,omitempty"`
	Note     string    `json:"note,omitempty"`
	Desc     string    `json:"desc,omitempty"`
	Lang     string    `json:"lang,omitempty"`
	Area     string    `json:"area,omitempty"`
	Year     string    `json:"year,omitempty"`
	State    string    `json:"state,omitempty"`
	Actor    string    `json:"actor,omitempty"`
	Director string    `json:"director,omitempty"`
	DD       []IMacCMSVideoDDTag
}

type IMacCMSVideosAndHeader struct {
	ListHeader IMacCMSListAttr
	Videos     []IMacCMSListVideoItem
}

type IMacCMSHomeData struct {
	ListHeader IMacCMSListAttr
	Category   []IMacCMSCategory
	Videos     []IMacCMSListVideoItem
}
