package maccms

import (
	"time"

	"d1y.io/neovideo/models/repos"
)

const (
	MacCMSReponseTypeXML  = "XML"
	MacCMSReponseTypeJSON = "JSON"
)

type IMacCMSListAttr struct {
	Page        int `json:"page"`
	PageCount   int `json:"page_count"`
	PageSize    int `json:"page_size"`
	RecordCount int `json:"record_count"`
}

// Deprecated: RawURL need parse, so use IMacCMSVideoDDTag
type IMacCMSVideoRawDDTag struct {
	Flag   string `json:"flag"`
	RawURL string `json:"raw_url"`
}

type IMacCMSVideoDDTag struct {
	Flag   string                     `json:"flag"`
	Videos []IMacCMSVideoDDTagWithURL `json:"videos"`
}

type IMacCMSVideoDDTagWithURL struct {
	Name string `json:"name"`
	URL  string `json:"url"` /* TODO: check url is [m3u8|mp4] */
}

type IMacCMSListVideoItem struct {
	Last     time.Time           `json:"last"`
	Id       int                 `json:"id"`
	Tid      int                 `json:"tid"`
	Name     string              `json:"name"`
	Type     string              `json:"type"`
	Dt       string              `json:"dt"`
	Note     string              `json:"note"`
	Desc     string              `json:"desc"`
	Lang     string              `json:"lang"`
	Area     string              `json:"area"`
	Year     string              `json:"year"`
	State    string              `json:"state"`
	Actor    string              `json:"actor"`
	Director string              `json:"director"`
	Pic      string              `json:"pic"`
	DD       []IMacCMSVideoDDTag `json:"dd"`
}

type IMacCMSVideosAndHeader struct {
	ListHeader IMacCMSListAttr        `json:"list_header"`
	Videos     []IMacCMSListVideoItem `json:"videos"`
}

type IMacCMSHomeData struct {
	ListHeader IMacCMSListAttr         `json:"list_header"`
	Category   []repos.IMacCMSCategory `json:"category"`
	Videos     []IMacCMSListVideoItem  `json:"videos"`
}
