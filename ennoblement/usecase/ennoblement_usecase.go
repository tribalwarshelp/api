package usecase

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/ennoblement"
	"github.com/tribalwarshelp/api/middleware"
)

type usecase struct {
	repo ennoblement.Repository
}

func New(repo ennoblement.Repository) ennoblement.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg ennoblement.FetchConfig) ([]*twmodel.Ennoblement, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.EnnoblementFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > ennoblement.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = ennoblement.FetchLimit
	}
	return ucase.repo.Fetch(ctx, cfg)
}
