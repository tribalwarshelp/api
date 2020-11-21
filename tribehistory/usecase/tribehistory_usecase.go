package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/tribehistory"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo tribehistory.Repository
}

func New(repo tribehistory.Repository) tribehistory.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg tribehistory.FetchConfig) ([]*models.TribeHistory, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.TribeHistoryFilter{}
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
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > tribehistory.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = tribehistory.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSortExpressions(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}
