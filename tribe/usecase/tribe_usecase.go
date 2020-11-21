package usecase

import (
	"context"
	"fmt"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo tribe.Repository
}

func New(repo tribe.Repository) tribe.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg tribe.FetchConfig) ([]*models.Tribe, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.TribeFilter{}
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
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > tribe.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = tribe.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSortExpressions(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*models.Tribe, error) {
	tribes, _, err := ucase.repo.Fetch(ctx, tribe.FetchConfig{
		Filter: &models.TribeFilter{
			ID: []int{id},
		},
		Limit:  1,
		Server: server,
		Count:  false,
	})
	if err != nil {
		return nil, err
	}
	if len(tribes) == 0 {
		return nil, fmt.Errorf("Tribe (ID: %s) not found.", id)
	}
	return tribes[0], nil
}
