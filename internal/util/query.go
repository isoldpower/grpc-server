package util

import (
	"net/url"
	"strconv"
)

func GetQueryUint64(params url.Values, key string) *uint64 {
	var value *uint64 = nil

	if parsed, err := strconv.ParseUint(params.Get(key), 10, 64); err == nil && parsed != 0 {
		value = &parsed
	}
	return value
}
