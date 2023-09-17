package impl

import (
	"regexp"
	"strings"

	"d1y.io/neovideo/common/json"
	"github.com/tidwall/gjson"
)

var jiexiURLFuzzyReg = regexp.MustCompile(`(?i)(jiexi|jiexiurl|url)=`)
var jiexiURLReg = regexp.MustCompile(`^https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
var jiexiURLAndNoteReg = regexp.MustCompile(`^(\S*:?)\s*(https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]=)`)

func parseJiexiWithLines(raw string) []JiexiParse {
	var result = make([]JiexiParse, 0)
	lines := strings.Split(raw, "\n")
	for _, item := range lines {
		s := strings.TrimSpace(item)
		if len(s) <= 6 || !jiexiURLFuzzyReg.MatchString(s) {
			continue
		}
		var i JiexiParse
		if jiexiURLReg.MatchString(s) {
			pl := strings.Split(s, "")
			if pl[len(pl)-1] == "=" {
				i.URL = s
			}
		} else {
			subs := jiexiURLAndNoteReg.FindStringSubmatch(s)
			if len(subs) == 3 {
				n, u := subs[1], subs[2]
				if strings.HasSuffix(n, ":") {
					n = strings.TrimRight(n, ":")
				}
				i.URL = u
				i.Name = n
			}
		}
		result = append(result, i)
	}
	return result
}

func parseJiexiWithJSON(raw string) []JiexiParse {
	return nil
}

func ParseJiexi(raw string) []JiexiParse {
	// 1. 先检测是否是 json
	if json.VerifyStringIsJSON(raw) && gjson.Valid(raw) {
		return parseJiexiWithJSON(raw)
	}
	// 2. 不是 json 的话就按每行处理
	return parseJiexiWithLines(raw)
}

func ParseMaccms(raw string) {

}
