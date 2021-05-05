package usecase

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/dailytribestats"
	"github.com/tribalwarshelp/api/middleware"
)

type usecase struct {
	repo dailytribestats.Repository
}

func New(repo dailytribestats.Repository) dailytribestats.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg dailytribestats.FetchConfig) ([]*twmodel.DailyTribeStats, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.DailyTribeStatsFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > dailytribestats.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = dailytribestats.FetchLimit
	}
	return ucase.repo.Fetch(ctx, cfg)
}
