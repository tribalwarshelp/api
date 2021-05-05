package usecase

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/tribechange"
)

type usecase struct {
	repo tribechange.Repository
}

func New(repo tribechange.Repository) tribechange.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg tribechange.FetchConfig) ([]*twmodel.TribeChange, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.TribeChangeFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > tribechange.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = tribechange.FetchLimit
	}
	return ucase.repo.Fetch(ctx, cfg)
}
