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

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > dailyplayerstats.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = dailyplayerstats.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSortExpressions(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}
