package utils

import "strings"

//IsEmpty check is string is empty
func IsEmpty(parameter string) bool {
	return len(strings.TrimSpace(parameter)) == 0
}
