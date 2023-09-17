package json

import (
	"strings"
)

// copy by https://github.com/waifu-project/movie/blob/dev/lib/utils/json.dart

var MAGIC_START_SYMBOL = []string{
	"[",
	"{",
}

var MAGIC_END_SYMBOL = []string{
	"]",
	"}",
}

func VerifyStringIsJSON(vJSON string) bool {
	if len(vJSON) <= 3 {
		return false
	}
	target := strings.TrimSpace(vJSON)
	start := string(target[0])
	end := string(target[len(target)-1])
	for index := 0; index < 2; index++ {
		startFlag := MAGIC_START_SYMBOL[index] == start
		endFlag := MAGIC_END_SYMBOL[index] == end
		if startFlag && endFlag {
			return true
		}
	}
	return false
}
