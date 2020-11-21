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

	if cfg.Filter.Limit > 0 {
		cfg.Limit = cfg.Filter.Limit
	}
	if cfg.Filter.Offset > 0 {
		cfg.Offset = cfg.Filter.Offset
	}
	if cfg.Filter.Sort != "" {
		cfg.Sort = append(cfg.Sort, cfg.Filter.Sort)
	}
	if cfg.Filter.PlayerFilter != nil {
		if cfg.Filter.PlayerFilter.Sort != "" {
			cfg.Sort = append(cfg.Sort, "player."+cfg.Filter.PlayerFilter.Sort)
		}
		if cfg.Filter.PlayerFilter.TribeFilter != nil && cfg.Filter.PlayerFilter.TribeFilter.Sort != "" {
			cfg.Sort = append(cfg.Sort, "tribe."+cfg.Filter.PlayerFilter.TribeFilter.Sort)
		}
	}

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > village.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = village.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSortExpressions(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*models.Village, error) {
	villages, _, err := ucase.repo.Fetch(ctx, village.FetchConfig{
		Filter: &models.VillageFilter{
			ID:    []int{id},
			Limit: 1,
		},
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
