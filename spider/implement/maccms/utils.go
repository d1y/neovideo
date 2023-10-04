package maccms

import (
	"regexp"
	"strings"
)

const (
	unknownPlayName = "未命名"
)

var (
	// 据我所知应该只有这几种格式是真实播放链接
	realPlayExt = []string{
		".m3u8",
		".mp4",
		".flv",
	}
)

var mustURLReg = regexp.MustCompile(`^https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)

func isEmbed(u string) bool {
	for _, ext := range realPlayExt {
		if strings.HasSuffix(u, ext) {
			return false
		}
	}
	return true
}

func parseDDRawURL(raw string) []IMacCMSVideoDDTagWithURL {
	var result = make([]IMacCMSVideoDDTagWithURL, 0)
	if mustURLReg.MatchString(raw) {
		// 某些播放链接就是这么特殊, 就纯纯是个m3u8/mp4播放链接
		// ^_^
		result = append(result, IMacCMSVideoDDTagWithURL{
			URL:   raw,
			Name:  unknownPlayName,
			Embed: isEmbed(raw),
		})
		return result
	}
	for _, rawItem := range strings.Split(raw, "#") {
		item := strings.Split(rawItem, "$")
		if len(item) <= 1 {
			continue
		}
		result = append(result, IMacCMSVideoDDTagWithURL{
			Name:  item[0],
			URL:   item[1],
			Embed: isEmbed(item[1]),
		})
	}
	return result
}
