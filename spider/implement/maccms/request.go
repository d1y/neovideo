package maccms

import (
	"strconv"
	"strings"
)

type XHRRequest struct {
	// ResponseType string `json:"r_type" form:"r_type"`
	RequestAction int    `json:"request_action,omitempty" form:"request_action"`
	ForceFetch    bool   `json:"force_fetch" form:"force_fetch"`
	Page          int    `json:"page" form:"page"`
	Keyword       string `json:"keyword" form:"keyword"`
	Action        string `json:"action" form:"action"`
	Category      int    `json:"category" form:"category"`
	Hour          int    `json:"hour" form:"hour"`
	Ids           string `json:"ids" form:"ids"`
}

func (xh *XHRRequest) GetIDs2Slice() []int {
	var result = make([]int, 0)
	ids := strings.Split(xh.Ids, ",")
	for _, id := range ids {
		id, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		result = append(result, id)
	}
	return result
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
