package usecase

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/dailyplayerstats"
	"github.com/tribalwarshelp/api/middleware"
)

type usecase struct {
	repo dailyplayerstats.Repository
}

func New(repo dailyplayerstats.Repository) dailyplayerstats.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg dailyplayerstats.FetchConfig) ([]*twmodel.DailyPlayerStats, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.DailyPlayerStatsFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > dailyplayerstats.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = dailyplayerstats.FetchLimit
	}
	return ucase.repo.Fetch(ctx, cfg)
}
