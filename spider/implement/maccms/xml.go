package maccms

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"d1y.io/neovideo/models/repos"
	"d1y.io/neovideo/spider/axios"
	"github.com/beevik/etree"
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

func (m *IMacCMS) xmlParseClassGetCategory(doc *etree.Element) []repos.IMacCMSCategory {
	var category []repos.IMacCMSCategory
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
			category = append(category, repos.IMacCMSCategory{
				Text: text,
				Id:   id,
			})
		}
	}
	return category
}

// TODO: 增强该方法, 传递
// 1. method 方法(get / post)
// 2. qs 参数
func (m *IMacCMS) xmlGetURL2XMLDocument(url string) (*etree.Document, error) {
	res, err := axios.Get(url)
	if err != nil {
		return nil, err
	}
	doc := etree.NewDocument()
	doc.ReadFromBytes(res)
	return doc, nil
}

func (m *IMacCMS) xmlGetURL2XMLDocumentWithRoot(url string) (*etree.Element, error) {
	doc, err := m.xmlGetURL2XMLDocument(url)
	if err != nil {
		return nil, err
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
				case "pic":
					item.Pic = text
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
								d.Videos = parseDDRawURL(strings.TrimSpace(e.Text()))
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
		case "recordcount":
			listAttr.RecordCount = val
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

func (m *IMacCMS) XMLGetHomeWithEtreeRoot(root *etree.Element) (IMacCMSHomeData, error) {
	var data IMacCMSHomeData
	if root == nil {
		return data, errors.New("root is nil")
	}
	for _, child := range root.Child {
		if c, ok := child.(*etree.Element); ok {
			if m.xmlIsClassTagWithXMLElement(c) {
				data.Category = m.xmlParseClassGetCategory(c)
			} else if m.xmlIsListTagWithXMLElement(c) {
				listAttr, videos, err := m.xmlParseList(c)
				if err != nil {
					return data, err
				}
				data.ListHeader = listAttr
				data.Videos = videos
			} else {
				continue
			}
		}
	}
	return data, nil
}

func (m *IMacCMS) XMLGetCategoryWithEtreeRoot(root *etree.Element) []repos.IMacCMSCategory {
	for _, child := range root.Child {
		if c, ok := child.(*etree.Element); ok {
			if m.xmlIsClassTagWithXMLElement(c) {
				return m.xmlParseClassGetCategory(c)
			}
		}
	}
	return []repos.IMacCMSCategory{}
}

func (m *IMacCMS) XMLGetSearchWithEtreeRoot(root *etree.Element) (IMacCMSVideosAndHeader, error) {
	for _, el := range root.Child {
		if e, ok := el.(*etree.Element); ok {
			if m.xmlIsListTagWithXMLElement(e) {
				a, b, c := m.xmlParseList(e)
				if c != nil {
					return IMacCMSVideosAndHeader{}, c
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

func (m *IMacCMS) XMLGetDetailWithEtreeRoot(root *etree.Element) (IMacCMSListAttr, []IMacCMSListVideoItem, error) {
	for _, el := range root.Child {
		if e, ok := el.(*etree.Element); ok {
			if m.xmlIsListTagWithXMLElement(e) {
				return m.xmlParseList(e)
			}
		}
	}
	return IMacCMSListAttr{}, []IMacCMSListVideoItem{}, errors.New("")
}

func (m *IMacCMS) XMLGetHome(page int, tid ...int) (IMacCMSHomeData, error) {
	res, err := m.qs.SetHome(page, tid).BuildGetRequest(m.ApiURL)
	if err != nil {
		return IMacCMSHomeData{}, err
	}
	doc := etree.NewDocument()
	doc.ReadFromBytes(res)
	return m.XMLGetHomeWithEtreeRoot(doc.Root())
}

func (m *IMacCMS) XMLGetCategory() ([]repos.IMacCMSCategory, error) {
	root, err := m.xmlGetURL2XMLDocumentWithRoot(m.ApiURL)
	if err != nil {
		return []repos.IMacCMSCategory{}, err
	}
	return m.XMLGetCategoryWithEtreeRoot(root), nil
}

func (m *IMacCMS) XMLGetSearch(keyword string, page int) (IMacCMSVideosAndHeader, error) {
	res, err := m.qs.SetPage(page).SetKeyword(keyword).BuildPostRequest(m.ApiURL)
	if err != nil {
		return IMacCMSVideosAndHeader{}, err
	}
	doc := etree.NewDocument()
	doc.ReadFromBytes(res)
	return m.XMLGetSearchWithEtreeRoot(doc.Root())
}

func (m *IMacCMS) XMLGetDetail(id int) (IMacCMSListAttr, []IMacCMSListVideoItem, error) {
	res, err := m.qs.SetAction("videolist").SetIDS(id).BuildPostRequest(m.ApiURL)
	if err != nil {
		return IMacCMSListAttr{}, []IMacCMSListVideoItem{}, err
	}
	doc := etree.NewDocument()
	doc.ReadFromBytes(res)
	return m.XMLGetDetailWithEtreeRoot(doc.Root())
}
