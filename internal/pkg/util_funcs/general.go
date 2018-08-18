package util_funcs

import (
	"strings"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a || strings.HasPrefix(a, b) {
			return true
		}
	}
	return false
}
