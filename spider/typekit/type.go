package typekit

type ISpider interface {
	GetCategory()
	GetHome(page int, tid ...int)
	GetSearch(keyword string, page int)
	GetDetail(id int)
}
