package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/playerhistory"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo playerhistory.Repository
}

func New(repo playerhistory.Repository) playerhistory.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg playerhistory.FetchConfig) ([]*models.PlayerHistory, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.PlayerHistoryFilter{}
	}

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > playerhistory.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = playerhistory.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}
