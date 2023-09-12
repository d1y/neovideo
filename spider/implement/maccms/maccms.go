package maccms

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/imroc/req/v3"
)

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

func (m *IMacCMS) isWhichTagWithXMLElement(e *etree.Element, tag string) bool {
	return e.Tag == tag
}

func (m *IMacCMS) isClassTagWithXMLElement(e *etree.Element) bool {
	return m.isWhichTagWithXMLElement(e, "class")
}
func (m *IMacCMS) isListTagWithXMLElement(e *etree.Element) bool {
	return m.isWhichTagWithXMLElement(e, "list")
}

func (m *IMacCMS) parseClassGetCategory(doc *etree.Element) []IMacCMSCategory {
	var category []IMacCMSCategory
	for _, tr := range doc.Child {
		if t, ok := tr.(*etree.Element); ok {
			var text = t.Text()
			var id int
			if t.Tag == "tr" {
				continue
			}
			for _, a := range t.Attr {
				if a.Key == "id" {
					id, _ = strconv.Atoi(a.Value)
					break
				} else {
					continue
				}
			}
			category = append(category, IMacCMSCategory{
				Text: text,
				Id:   id,
			})
		}
	}
	return category
}

func (m *IMacCMS) getURL2XMLDocument(url string) (*etree.Document, error) {
	res, err := req.Get(url)
	if err != nil {
		return &etree.Document{}, err
	}
	doc := etree.NewDocument()
	doc.ReadFrom(res.Body)
	return doc, nil
}

func (m *IMacCMS) getURL2XMLDocumentWithRoot(url string) (*etree.Element, error) {
	doc, err := m.getURL2XMLDocument(url)
	if err != nil {
		return &etree.Element{}, err
	}
	root := doc.Root()
	return root, nil
}

func (m *IMacCMS) parseList(doc *etree.Element) (IMacCMSListAttr, []IMacCMSListVideoItem, error) {
	var listAttr IMacCMSListAttr
	var videos []IMacCMSListVideoItem
	var parseVideo = func(ele *etree.Element) (IMacCMSListVideoItem, error) {
		var item IMacCMSListVideoItem
		for _, child := range ele.Child {
			if c, ok := child.(*etree.Element); ok {
				var text = c.Text()
				switch c.Tag {
				case "last":
					item.Last, _ = time.Parse(time.DateTime, text)
				case "id":
					item.Id, _ = strconv.Atoi(text)
				case "tid":
					item.Tid, _ = strconv.Atoi(text)
				case "name":
					item.Name = text
				case "type":
					item.Type = text
				case "dt":
					item.Dt = text
				case "note":
					item.Note = text
				case "des":
					item.Desc = text
				case "lang":
					item.Lang = text
				case "area":
					item.Area = text
				case "year":
					item.Year = text
				case "state":
					item.State = text
				case "actor":
					item.Actor = text
				case "director":
					item.Director = text
				case "dl":
					var dd []IMacCMSVideoDDTag
					for _, el := range c.Child {
						if e, ok := el.(*etree.Element); ok {
							if m.isWhichTagWithXMLElement(e, "dd") {
								var d IMacCMSVideoDDTag
								for _, attr := range e.Attr {
									if attr.Key == "flag" {
										d.Flag = attr.Value
										break
									}
								}
								d.RawURL = strings.TrimSpace(e.Text())
								dd = append(dd, d)
							}
						}
					}
					item.DD = dd
				}
			}
		}
		return item, nil
	}
	for _, a := range doc.Attr {
		val, _ := strconv.Atoi(a.Value)
		switch a.Key {
		case "page":
			listAttr.Page = val
		case "pagecount":
			listAttr.PageCount = val
		case "pagesize":
			listAttr.PageSize = val
		}
	}
	for _, child := range doc.Child {
		if e, ok := child.(*etree.Element); ok {
			if m.isWhichTagWithXMLElement(e, "video") {
				video, err := parseVideo(e)
				if err == nil {
					videos = append(videos, video)
				}
			} else {
				continue
			}
		}
	}
	return listAttr, videos, nil
}

func (m *IMacCMS) GetHome() (IMacCMSHomeData, error) {
	root, err := m.getURL2XMLDocumentWithRoot(m.ApiURL)
	if err != nil {
		return IMacCMSHomeData{}, err
	}
	var data IMacCMSHomeData
	for _, child := range root.Child {
		if c, ok := child.(*etree.Element); ok {
			if m.isClassTagWithXMLElement(c) {
				data.Category = m.parseClassGetCategory(c)
			} else if m.isListTagWithXMLElement(c) {
				listAttr, videos, _ := m.parseList(c)
				data.ListHeader = listAttr
				data.Videos = videos
			} else {
				continue
			}
		}
	}
	return data, nil
}

func (m *IMacCMS) GetCategory() ([]IMacCMSCategory, error) {
	root, err := m.getURL2XMLDocumentWithRoot(m.ApiURL)
	if err != nil {
		return []IMacCMSCategory{}, err
	}
	for _, child := range root.Child {
		if c, ok := child.(*etree.Element); ok {
			if m.isClassTagWithXMLElement(c) {
				return m.parseClassGetCategory(c), nil
			}
		}
	}
	return []IMacCMSCategory{}, nil
}

func (m *IMacCMS) GetSearch(keyword string, page int) (IMacCMSVideosAndHeader, error) {
	res, err := req.R().SetQueryParams(map[string]string{
		// "ac": "videolist",
		"pg": strconv.Itoa(page),
		"wd": keyword,
	}).Post(m.ApiURL)
	if err != nil {
		return IMacCMSVideosAndHeader{}, err
	}
	doc := etree.NewDocument()
	doc.ReadFrom(res.Body)
	for _, el := range doc.Root().Child {
		if e, ok := el.(*etree.Element); ok {
			if m.isListTagWithXMLElement(e) {
				a, b, c := m.parseList(e)
				if c != nil {
					return IMacCMSVideosAndHeader{}, err
				}
				return IMacCMSVideosAndHeader{
					ListHeader: a,
					Videos:     b,
				}, nil
			}
		}
	}
	return IMacCMSVideosAndHeader{}, errors.New("not found")
}

func (m *IMacCMS) GetDetail(id int) (IMacCMSListAttr, []IMacCMSListVideoItem, error) {
	res, err := req.R().SetQueryParams(map[string]string{
		"ac":  "videolist",
		"ids": strconv.Itoa(id),
	}).Get(m.ApiURL)
	if err != nil {
		return IMacCMSListAttr{}, []IMacCMSListVideoItem{}, err
	}
	doc := etree.NewDocument()
	doc.ReadFrom(res.Body)
	for _, el := range doc.Root().Child {
		if e, ok := el.(*etree.Element); ok {
			if m.isListTagWithXMLElement(e) {
				return m.parseList(e)
			}
		}
	}
	return IMacCMSListAttr{}, []IMacCMSListVideoItem{}, err
}