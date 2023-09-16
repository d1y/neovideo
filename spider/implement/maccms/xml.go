package maccms

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/imroc/req/v3"
)

func (m *IMacCMS) xmlIsWhichTagWithXMLElement(e *etree.Element, tag string) bool {
	return e.Tag == tag
}

func (m *IMacCMS) xmlIsClassTagWithXMLElement(e *etree.Element) bool {
	return m.xmlIsWhichTagWithXMLElement(e, "class")
}
func (m *IMacCMS) xmlIsListTagWithXMLElement(e *etree.Element) bool {
	return m.xmlIsWhichTagWithXMLElement(e, "list")
}

func (m *IMacCMS) xmlParseClassGetCategory(doc *etree.Element) []IMacCMSCategory {
	var category []IMacCMSCategory
	for _, tr := range doc.Child {
		if t, ok := tr.(*etree.Element); ok {
			var text = t.Text()
			var id int
			if m.xmlIsWhichTagWithXMLElement(t, "tr") {
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

func (m *IMacCMS) xmlGetURL2XMLDocument(url string) (*etree.Document, error) {
	res, err := req.Get(url)
	if err != nil {
		return &etree.Document{}, err
	}
	doc := etree.NewDocument()
	doc.ReadFrom(res.Body)
	return doc, nil
}

func (m *IMacCMS) xmlGetURL2XMLDocumentWithRoot(url string) (*etree.Element, error) {
	doc, err := m.xmlGetURL2XMLDocument(url)
	if err != nil {
		return &etree.Element{}, err
	}
	root := doc.Root()
	return root, nil
}

func (m *IMacCMS) xmlParseList(doc *etree.Element) (IMacCMSListAttr, []IMacCMSListVideoItem, error) {
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
							if m.xmlIsWhichTagWithXMLElement(e, "dd") {
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
			if m.xmlIsWhichTagWithXMLElement(e, "video") {
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

func (m *IMacCMS) XMLGetHome() (IMacCMSHomeData, error) {
	root, err := m.xmlGetURL2XMLDocumentWithRoot(m.ApiURL)
	if err != nil {
		return IMacCMSHomeData{}, err
	}
	var data IMacCMSHomeData
	for _, child := range root.Child {
		if c, ok := child.(*etree.Element); ok {
			if m.xmlIsClassTagWithXMLElement(c) {
				data.Category = m.xmlParseClassGetCategory(c)
			} else if m.xmlIsListTagWithXMLElement(c) {
				listAttr, videos, _ := m.xmlParseList(c)
				data.ListHeader = listAttr
				data.Videos = videos
			} else {
				continue
			}
		}
	}
	return data, nil
}

func (m *IMacCMS) XMLGetCategory() ([]IMacCMSCategory, error) {
	root, err := m.xmlGetURL2XMLDocumentWithRoot(m.ApiURL)
	if err != nil {
		return []IMacCMSCategory{}, err
	}
	for _, child := range root.Child {
		if c, ok := child.(*etree.Element); ok {
			if m.xmlIsClassTagWithXMLElement(c) {
				return m.xmlParseClassGetCategory(c), nil
			}
		}
	}
	return []IMacCMSCategory{}, nil
}

func (m *IMacCMS) XMLGetSearch(keyword string, page int) (IMacCMSVideosAndHeader, error) {
	res, err := m.qs.SetPage(page).SetKeyword(keyword).BuildRequest().Post(m.ApiURL)
	if err != nil {
		return IMacCMSVideosAndHeader{}, err
	}
	doc := etree.NewDocument()
	doc.ReadFrom(res.Body)
	for _, el := range doc.Root().Child {
		if e, ok := el.(*etree.Element); ok {
			if m.xmlIsListTagWithXMLElement(e) {
				a, b, c := m.xmlParseList(e)
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

func (m *IMacCMS) XMLGetDetail(id int) (IMacCMSListAttr, []IMacCMSListVideoItem, error) {
	res, err := m.qs.SetAction("videolist").SetIDS(id).BuildRequest().Get(m.ApiURL)
	if err != nil {
		return IMacCMSListAttr{}, []IMacCMSListVideoItem{}, err
	}
	doc := etree.NewDocument()
	doc.ReadFrom(res.Body)
	for _, el := range doc.Root().Child {
		if e, ok := el.(*etree.Element); ok {
			if m.xmlIsListTagWithXMLElement(e) {
				return m.xmlParseList(e)
			}
		}
	}
	return IMacCMSListAttr{}, []IMacCMSListVideoItem{}, err
}
