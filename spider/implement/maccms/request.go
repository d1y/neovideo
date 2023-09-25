package maccms

type XHRRequest struct {
	// ResponseType string `json:"r_type" form:"r_type"`
	RequestAction int    `json:"request_action,omitempty"`
	Page          int    `json:"page" form:"page"`
	Keyword       string `json:"keyword" form:"keyword"`
	Action        string `json:"action" form:"action"`
	Category      int    `json:"category" form:"category"`
	Hour          int    `json:"hour" form:"hour"`
	Ids           []int  `json:"ids" form:"ids"`
}

// func (xhr *XHRRequest) SetPage(page int) {
// 	xhr.Page = page
// }

// func (xhr *XHRRequest) SetKeyword(kw string) {
// 	xhr.Keyword = kw
// }

// func (xhr *XHRRequest) SetAction(ac string) {
// 	xhr.Action = ac
// }

// func (xhr *XHRRequest) SetCategory(cat int) {
// 	xhr.Category = cat
// }

// func (xhr *XHRRequest) AddIDS(ids ...int) {
// 	xhr.Ids = append(xhr.Ids, ids...)
// }
