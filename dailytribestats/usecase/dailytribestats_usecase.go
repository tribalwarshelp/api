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

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > dailytribestats.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = dailytribestats.FetchLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}
