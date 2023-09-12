package typekit

type ISpider interface {
	GetCategory()
	GetHome()
	GetSearch(keyword string, page int)
	GetDetail()
}
