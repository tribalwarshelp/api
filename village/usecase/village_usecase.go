package usecase

import (
	"context"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/village"
)

type usecase struct {
	repo village.Repository
}

func New(repo village.Repository) village.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg village.FetchConfig) ([]*twmodel.Village, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.VillageFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > village.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = village.FetchLimit
	}
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*twmodel.Village, error) {
	villages, _, err := ucase.repo.Fetch(ctx, village.FetchConfig{
		Filter: &twmodel.VillageFilter{
			ID: []int{id},
		},
		Limit:  1,
		Server: server,
		Select: true,
		Count:  false,
	})
	if err != nil {
		return nil, err
	}
	if len(villages) == 0 {
		return nil, errors.Errorf("Village (ID: %d) not found.", id)
	}
	return villages[0], nil
}
