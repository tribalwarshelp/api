package utils

import (
	"regexp"
	"strings"
)

var (
	sortexprRegex = regexp.MustCompile(`^[\p{L}\_\.]+$`)
)

func SanitizeSortExpression(expr string) string {
	trimmed := strings.TrimSpace(expr)
	splitted := strings.Split(trimmed, " ")
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

func SanitizeSortExpressions(exprs []string) []string {
	filtered := []string{}
	for _, expr := range exprs {
		sanitized := SanitizeSortExpression(expr)
		if sanitized != "" {
			filtered = append(filtered, sanitized)
		}
	}
	return filtered
}
