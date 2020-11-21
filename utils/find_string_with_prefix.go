package utils

import "strings"

func FindStringWithPrefix(sl []string, prefix string) string {
	for _, s := range sl {
		if strings.HasPrefix(s, prefix) {
			return s
		}
	}
	return ""
}
