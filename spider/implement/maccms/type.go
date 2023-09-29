package maccms

import "time"

const (
	MacCMSReponseTypeXML  = "XML"
	MacCMSReponseTypeJSON = "JSON"
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

// Deprecated: RawURL need parse, so use IMacCMSVideoDDTag
type IMacCMSVideoRawDDTag struct {
	Flag   string `json:"flag,omitempty"`
	RawURL string `json:"raw_url,omitempty"`
}

type IMacCMSVideoDDTag struct {
	Flag   string                     `json:"flag,omitempty"`
	Videos []IMacCMSVideoDDTagWithURL `json:"videos,omitempty"`
}

type IMacCMSVideoDDTagWithURL struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"` /* TODO: check url is [m3u8|mp4] */
}

type IMacCMSListVideoItem struct {
	Last     time.Time           `json:"last,omitempty"`
	Id       int                 `json:"id,omitempty"`
	Tid      int                 `json:"tid,omitempty"`
	Name     string              `json:"name,omitempty"`
	Type     string              `json:"type,omitempty"`
	Dt       string              `json:"dt,omitempty"`
	Note     string              `json:"note,omitempty"`
	Desc     string              `json:"desc,omitempty"`
	Lang     string              `json:"lang,omitempty"`
	Area     string              `json:"area,omitempty"`
	Year     string              `json:"year,omitempty"`
	State    string              `json:"state,omitempty"`
	Actor    string              `json:"actor,omitempty"`
	Director string              `json:"director,omitempty"`
	DD       []IMacCMSVideoDDTag `json:"dd,omitempty"`
}

type IMacCMSVideosAndHeader struct {
	ListHeader IMacCMSListAttr        `json:"list_header,omitempty"`
	Videos     []IMacCMSListVideoItem `json:"videos,omitempty"`
}

type IMacCMSHomeData struct {
	ListHeader IMacCMSListAttr        `json:"list_header,omitempty"`
	Category   []IMacCMSCategory      `json:"category,omitempty"`
	Videos     []IMacCMSListVideoItem `json:"videos,omitempty"`
}
