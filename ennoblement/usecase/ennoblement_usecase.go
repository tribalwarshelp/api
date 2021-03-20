package usecase

import (
	"context"
	"time"

	"github.com/tribalwarshelp/api/ennoblement"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo ennoblement.Repository
}

func New(repo ennoblement.Repository) ennoblement.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg ennoblement.FetchConfig) ([]*models.Ennoblement, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.EnnoblementFilter{}
	}

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > ennoblement.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = ennoblement.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) FetchLiveEnnoblements(ctx context.Context, server string) ([]*models.LiveEnnoblement, error) {
	limit := 0
	if !middleware.CanExceedLimit(ctx) {
		limit = ennoblement.PaginationLimit
	}
	ennoblements, _, err := ucase.repo.Fetch(ctx, ennoblement.FetchConfig{
		Server: server,
		Count:  false,
		Filter: &models.EnnoblementFilter{
			EnnobledAtGTE: time.Now().Add(-1 * time.Hour),
		},
		Limit: limit,
		Sort:  []string{"ennobled_at ASC"},
	})
	if err != nil {
		return nil, err
	}
	return convertToLiveEnnoblements(ennoblements), nil
}

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
