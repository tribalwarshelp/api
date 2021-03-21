package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/serverstats"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo serverstats.Repository
}

func New(repo serverstats.Repository) serverstats.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg serverstats.FetchConfig) ([]*models.ServerStats, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.ServerStatsFilter{}
	}

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > serverstats.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = serverstats.FetchLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}
