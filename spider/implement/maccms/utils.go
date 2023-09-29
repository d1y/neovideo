package maccms

import "strings"

func parseDDRawURL(raw string) []IMacCMSVideoDDTagWithURL {
	var result = make([]IMacCMSVideoDDTagWithURL, 0)
	for _, rawItem := range strings.Split(raw, "#") {
		item := strings.Split(rawItem, "$")
		if len(item) != 2 {
			continue
		}
		result = append(result, IMacCMSVideoDDTagWithURL{
			Name: item[0],
			URL:  item[1],
		})
	}
	return result
}
