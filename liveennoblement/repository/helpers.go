package repository

import (
	"github.com/tribalwarshelp/shared/models"
)

func convertToLiveEnnoblements(ennoblements []*models.Ennoblement) []*models.LiveEnnoblement {
	lv := []*models.LiveEnnoblement{}
	for _, e := range ennoblements {
		lv = append(lv, &models.LiveEnnoblement{
			VillageID:  e.VillageID,
			NewOwnerID: e.NewOwnerID,
			OldOwnerID: e.OldOwnerID,
			EnnobledAt: e.EnnobledAt,
		})
	}
	return lv
}
