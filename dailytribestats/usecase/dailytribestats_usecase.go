package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/dailytribestats"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo dailytribestats.Repository
}

func New(repo dailytribestats.Repository) dailytribestats.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg dailytribestats.FetchConfig) ([]*models.DailyTribeStats, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.DailyTribeStatsFilter{}
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
	if cfg.Filter.TribeFilter != nil {
		if cfg.Filter.TribeFilter.Sort != "" {
			cfg.Sort = append(cfg.Sort, "tribe."+cfg.Filter.TribeFilter.Sort)
		}
	}

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > dailytribestats.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = dailytribestats.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSortExpressions(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}
