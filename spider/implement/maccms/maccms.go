package maccms

import (
	"strconv"
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

type IMacCMSListVideoItem struct {
	Last time.Time `json:"last,omitempty"`
	Id   int       `json:"id,omitempty"`
	Tid  int       `json:"tid,omitempty"`
	Name string    `json:"name,omitempty"`
	Type string    `json:"type,omitempty"`
	Dt   string    `json:"dt,omitempty"`
	Note string    `json:"note,omitempty"`
}

type IMacCMSHomeData struct {
	Category   []IMacCMSCategory
	ListHeader IMacCMSListAttr
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

func (m *IMacCMS) parseList(doc *etree.Element) (IMacCMSListAttr, []IMacCMSListVideoItem, error) {
	var listAttr IMacCMSListAttr
	var videos []IMacCMSListVideoItem
	var parseVideo = func(ele *etree.Element) (IMacCMSListVideoItem, error) {
		var item IMacCMSListVideoItem
		for _, child := range ele.Child {
			if c, ok := child.(*etree.Element); ok {
				switch c.Tag {
				case "last":
					item.Last, _ = time.Parse(time.DateTime, c.Text())
				case "id":
					item.Id, _ = strconv.Atoi(c.Text())
				case "tid":
					item.Tid, _ = strconv.Atoi(c.Text())
				case "name":
					item.Name = c.Text()
				case "type":
					item.Type = c.Text()
				case "dt":
					item.Dt = c.Text()
				case "note":
					item.Note = c.Text()
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
			if e.Tag == "video" {
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

func (m *IMacCMS) GetHome() (IMacCMSHomeData, error) {
	root, err := m.getURL2XMLDocumentWithRoot(m.ApiURL)
	if err != nil {
		return IMacCMSHomeData{}, err
	}
	var data IMacCMSHomeData
	for _, child := range root.Child {
		if c, ok := child.(*etree.Element); ok {
			if c.Tag == "class" {
				data.Category = m.parseClassGetCategory(c)
			} else if c.Tag == "list" {
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
			if c.Tag == "class" {
				return m.parseClassGetCategory(c), nil
			}
		}
	}
	return []IMacCMSCategory{}, nil
}
