package usecase

import (
	"context"

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
	if cfg.Filter.Limit > 0 {
		cfg.Limit = cfg.Filter.Limit
	}
	if cfg.Filter.Offset > 0 {
		cfg.Offset = cfg.Filter.Offset
	}
	if cfg.Filter.Sort != "" {
		cfg.Sort = append(cfg.Sort, cfg.Filter.Sort)
	}

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > ennoblement.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = ennoblement.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSortExpressions(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}
