package usecase

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/serverstats"
)

type usecase struct {
	repo serverstats.Repository
}

func New(repo serverstats.Repository) serverstats.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg serverstats.FetchConfig) ([]*twmodel.ServerStats, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.ServerStatsFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > serverstats.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = serverstats.FetchLimit
	}
	return ucase.repo.Fetch(ctx, cfg)
}
