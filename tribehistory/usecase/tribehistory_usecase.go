package usecase

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/tribehistory"
)

type usecase struct {
	repo tribehistory.Repository
}

func New(repo tribehistory.Repository) tribehistory.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg tribehistory.FetchConfig) ([]*twmodel.TribeHistory, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.TribeHistoryFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > tribehistory.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = tribehistory.FetchLimit
	}
	if len(cfg.Sort) > tribehistory.MaxOrders {
		cfg.Sort = cfg.Sort[0:tribehistory.MaxOrders]
	}
	return ucase.repo.Fetch(ctx, cfg)
}
