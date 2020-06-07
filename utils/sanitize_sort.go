package utils

import "strings"

func SanitizeSort(sort string) string {
	trimmed := strings.TrimSpace(sort)
	splitted := strings.Split(trimmed, " ")
	length := len(splitted)
	if length != 2 {
		return ""
	}
	keyword := "ASC"
	if strings.ToUpper(splitted[1]) == "DESC" {
		keyword = "DESC"
	}
	return strings.ToLower(Underscore(splitted[0])) + " " + keyword
}
