package utils

import (
	"github.com/tribalwarshelp/shared/models"
)

func LanguageTagFromServerKey(key string) models.LanguageTag {
	if len(key) < 2 {
		return ""
	}
	return models.LanguageTag(key[0:2])
}
