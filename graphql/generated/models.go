// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"github.com/tribalwarshelp/shared/models"
)

type LangVersionList struct {
	Items []*models.LangVersion `json:"items"`
	Total int                   `json:"total"`
}

type PlayerList struct {
	Items []*models.Player `json:"items"`
	Total int              `json:"total"`
}

type ServerList struct {
	Items []*models.Server `json:"items"`
	Total int              `json:"total"`
}

type TribeList struct {
	Items []*models.Tribe `json:"items"`
	Total int             `json:"total"`
}

type VillageList struct {
	Items []*models.Village `json:"items"`
	Total int               `json:"total"`
}
