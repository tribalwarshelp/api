package usecase

import (
	"context"
	"fmt"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo village.Repository
}

func New(repo village.Repository) village.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg village.FetchConfig) ([]*models.Village, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.VillageFilter{}
	}

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > village.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = village.FetchLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*models.Village, error) {
	villages, _, err := ucase.repo.Fetch(ctx, village.FetchConfig{
		Filter: &models.VillageFilter{
			ID: []int{id},
		},
		Limit:  1,
		Server: server,
	})
	if err != nil {
		return nil, err
	}
	if len(villages) == 0 {
		return nil, fmt.Errorf("Village (ID: %d) not found.", id)
	}
	return villages[0], nil
}
