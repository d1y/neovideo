package typekit

import (
	"strconv"
	"unsafe"
)

// raw string to int64
func Str2Int64(raw string) (int64, error) {
	return strconv.ParseInt(raw, 10, 64)
}

// unsafe int64 to int
// enjoy happy
func Int642Int(raw int64) int {
	ptr := (*int)(unsafe.Pointer(&raw))
	return *ptr
}

func Int64Slice2Int(opt ...int64) []int {
	res := make([]int, len(opt))
	for i, v := range opt {
		res[i] = Int642Int(v)
	}
	return res
}
