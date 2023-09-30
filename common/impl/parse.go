package impl

import (
	"errors"
	"regexp"
	"strings"

	"d1y.io/neovideo/common/json"
	"d1y.io/neovideo/spider/implement/maccms"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/slices"
)

var jiexiURLFuzzyReg = regexp.MustCompile(`(?i)(jiexi|jiexiurl|url)?=`)
var mustUrlReg = regexp.MustCompile(`^https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
var urlReg = regexp.MustCompile(`https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
var jiexiURLAndNoteReg = regexp.MustCompile(`^(\S*:?)\s*(https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]=)`)
var ignoreReg = regexp.MustCompile(`^(//|;)`)
var symbolReg = regexp.MustCompile(`[,.:：，。]*$`)

var group18 = []string{
	"18禁",
	"18+",
}

func gjsonGGGetString(g gjson.Result, s ...string) (string, error) {
	for _, v := range s {
		s := g.Get(v).String()
		if len(s) >= 1 {
			return s, nil
		}
	}
	return "", errors.New("not match")
}

func parseJiexiWithLines(raw string) []JiexiParse {
	var result = make([]JiexiParse, 0)
	lines := strings.Split(raw, "\n")
	for _, item := range lines {
		s := strings.TrimSpace(item)
		if len(s) <= 6 || !jiexiURLFuzzyReg.MatchString(s) || ignoreReg.MatchString(s) {
			continue
		}
		var i JiexiParse
		if mustUrlReg.MatchString(s) {
			pl := strings.Split(s, "")
			if pl[len(pl)-1] == "=" {
				i.URL = s
			}
		} else {
			subs := jiexiURLAndNoteReg.FindStringSubmatch(s)
			if len(subs) == 3 {
				n, u := subs[1], subs[2]
				n = symbolReg.ReplaceAllString(n, "")
				i.URL = u
				i.Name = n
			}
		}
		result = append(result, i)
	}
	return result
}

func parseJiexiWithJSON(raw string) []JiexiParse {
	var result = make([]JiexiParse, 0)
	gs := gjson.Parse(raw)
	if gs.IsArray() {
		for _, item := range gs.Array() {
			if item.IsObject() {
				t, e1 := gjsonGGGetString(item, "name", "title")
				u, e2 := gjsonGGGetString(item, "url", "jiexi_url", "jiexi_url", "jiexiUrl") /* 不校验 url */
				if e1 != nil && e2 != nil {
					continue
				}
				result = append(result, JiexiParse{URL: u, Name: t})
			} else {
				if mustUrlReg.MatchString(item.String()) {
					result = append(result, JiexiParse{URL: item.String()})
				}
			}
		}
	} /* else if gs.IsObject() {
	}*/
	return result
}

func ParseJiexi(raw string) []JiexiParse {
	var m = make(map[string]JiexiParse)
	if json.VerifyStringIsJSON(raw) && gjson.Valid(raw) {
		for _, item := range parseJiexiWithJSON(raw) {
			if len(item.URL) >= 1 {
				m[item.URL] = item
			}
		}
	} else {
		for _, item := range parseJiexiWithLines(raw) {
			if len(item.URL) >= 1 {
				m[item.URL] = item
			}
		}
	}
	var result = make([]JiexiParse, 0)
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func ParseMaccms(raw string) []MacCMSParse {
	if json.VerifyStringIsJSON(raw) && gjson.Valid(raw) {
		return parseMaccmsWithJSON(raw)
	} else {
		return parseMaccmsWithLines(raw)
	}
}

func deducedResType(url string) string {
	// TODO: 更智能检测机制
	if strings.HasSuffix(url, "/xml") {
		return maccms.MacCMSReponseTypeXML
	}
	return ""
}

func parseMaccmsWithLines(raw string) []MacCMSParse {
	lines := strings.Split(raw, "\n")
	var result = make([]MacCMSParse, 0)
	for _, item := range lines {
		if len(item) <= 6 || ignoreReg.MatchString(item) || !urlReg.MatchString(item) {
			continue
		}
		seps := strings.Split(item, ",")
		if len(seps) < 2 {
			continue
		}
		var mp MacCMSParse
		mp.Name = strings.TrimSpace(seps[0])
		mp.Api = strings.TrimSpace(seps[1])
		checkResType := true
		if len(seps) >= 3 {
			for i := 1; i < len(seps); i++ {
				n := strings.ToLower(strings.TrimSpace(seps[i]))
				switch n {
				case "nsfw":
					mp.R18 = true
				case "不解析":
					mp.JiexiParse = false
				case "解析":
					mp.JiexiParse = true
				case "xml":
					mp.RespType = maccms.MacCMSReponseTypeXML
					checkResType = false
				case "json":
					mp.RespType = maccms.MacCMSReponseTypeJSON
					checkResType = false
				}
			}
		}
		if checkResType {
			mp.RespType = deducedResType(mp.Api)
		}
		result = append(result, mp)
	}
	return result
}

func parseMaccmsWithJSON(raw string) []MacCMSParse {
	var result = make([]MacCMSParse, 0)
	gs := gjson.Parse(raw)
	if gs.IsArray() {
		for _, item := range gs.Array() {
			if item.IsObject() {
				var mp MacCMSParse
				var checkResType = true
				mp.Name, _ = gjsonGGGetString(item, "name", "title")
				mp.Api, _ = gjsonGGGetString(item, "api", "url")
				mp.JiexiParse = item.Get("jiexi_parse").Bool()
				mp.JiexiURL, _ = gjsonGGGetString(item, "jiexi_url")
				group := item.Get("group").String()
				if slices.Contains(group18, group) {
					mp.R18 = true
				} else {
					mp.R18 = item.Get("nsfw").Bool()
				}
				rType, _ := gjsonGGGetString(item, "res_type", "type")
				rType = strings.ToLower(strings.TrimSpace(rType))
				if len(rType) >= 3 {
					switch rType {
					case "xml":
						mp.RespType = maccms.MacCMSReponseTypeXML
						checkResType = false
					case "json":
						mp.RespType = maccms.MacCMSReponseTypeJSON
						checkResType = false
					}
				}
				if checkResType {
					mp.RespType = deducedResType(mp.Api)
				}
				result = append(result, mp)
			}
		}
	}
	return result
}
