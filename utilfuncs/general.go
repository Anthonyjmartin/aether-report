package utilfuncs

import (
	"strings"
)

//StringInSlice return true if string is in slice else return false.
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a || strings.HasPrefix(a, b) {
			return true
		}
	}
	return false
}
