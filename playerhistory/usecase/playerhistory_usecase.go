package usecase

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/playerhistory"
)

type usecase struct {
	repo playerhistory.Repository
}

func New(repo playerhistory.Repository) playerhistory.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg playerhistory.FetchConfig) ([]*twmodel.PlayerHistory, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.PlayerHistoryFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > playerhistory.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = playerhistory.FetchLimit
	}
	if len(cfg.Sort) > playerhistory.MaxOrders {
		cfg.Sort = cfg.Sort[0:playerhistory.MaxOrders]
	}
	return ucase.repo.Fetch(ctx, cfg)
}
