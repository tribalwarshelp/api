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

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > tribehistory.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = tribehistory.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}
