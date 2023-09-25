package maccms

const (
	proxyActionWithHome = iota
	proxyActionWithCategory
	proxyActionWithDetail
	proxyActionWithSearch
)

// type proxyAction struct {
// 	ac int
// }

// func newProxyAction(ac int) *proxyAction {
// 	return &proxyAction{ac: ac}
// }

// func (pa *proxyAction) String() string {
// 	switch pa.ac {
// 	case proxyActionWithHome:
// 		return "home"
// 	case proxyActionWithCategory:
// 		return "category"
// 	case proxyActionWithDetail:
// 		return "detail"
// 	case proxyActionWithSearch:
// 		return "search"
// 	default:
// 		return "unknow"
// 	}
// }

// func (pa *proxyAction) Ac() int {
// 	return pa.ac
// }
