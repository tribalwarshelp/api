// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"github.com/tribalwarshelp/shared/models"
)

type LangVersionsList struct {
	Items []*models.LangVersion `json:"items"`
	Total int                   `json:"total"`
}

type PlayersList struct {
	Items []*models.Player `json:"items"`
	Total int              `json:"total"`
}

type ServersList struct {
	Items []*models.Server `json:"items"`
	Total int              `json:"total"`
}

type TribesList struct {
	Items []*models.Tribe `json:"items"`
	Total int             `json:"total"`
}

type VillagesList struct {
	Items []*models.Village `json:"items"`
	Total int               `json:"total"`
}
