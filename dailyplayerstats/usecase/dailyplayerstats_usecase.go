package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/dailyplayerstats"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo dailyplayerstats.Repository
}

func New(repo dailyplayerstats.Repository) dailyplayerstats.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg dailyplayerstats.FetchConfig) ([]*models.DailyPlayerStats, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.DailyPlayerStatsFilter{}
	}

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > dailyplayerstats.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = dailyplayerstats.FetchLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}
