package utils

import (
	"regexp"
	"strings"
)

var (
	sortexprRegex = regexp.MustCompile(`^[\p{L}\_\.]+$`)
)

func SanitizeSort(expr string) string {
	splitted := strings.Split(strings.TrimSpace(expr), " ")
	length := len(splitted)
	if length != 2 || !sortexprRegex.Match([]byte(splitted[0])) {
		return ""
	}
	table := ""
	column := splitted[0]
	if strings.Contains(splitted[0], ".") {
		columnAndTable := strings.Split(splitted[0], ".")
		table = Underscore(columnAndTable[0]) + "."
		column = columnAndTable[1]
	}
	keyword := "ASC"
	if strings.ToUpper(splitted[1]) == "DESC" {
		keyword = "DESC"
	}
	return strings.ToLower(table+Underscore(column)) + " " + keyword
}

func SanitizeSorts(sorts []string) []string {
	sanitized := []string{}
	for _, sort := range sorts {
		sanitizedSort := SanitizeSort(sort)
		if sanitizedSort != "" {
			sanitized = append(sanitized, sanitizedSort)
		}
	}
	return sanitized
}
