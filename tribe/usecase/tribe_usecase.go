package usecase

import (
	"context"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/tw/twmodel"
	"strings"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/tribe"
)

type usecase struct {
	repo tribe.Repository
}

func New(repo tribe.Repository) tribe.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg tribe.FetchConfig) ([]*twmodel.Tribe, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.TribeFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > tribe.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = tribe.FetchLimit
	}
	if len(cfg.Sort) > tribe.MaxOrders {
		cfg.Sort = cfg.Sort[0:tribe.MaxOrders]
	}
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*twmodel.Tribe, error) {
	tribes, _, err := ucase.repo.Fetch(ctx, tribe.FetchConfig{
		Filter: &twmodel.TribeFilter{
			ID: []int{id},
		},
		Limit:  1,
		Server: server,
		Count:  false,
		Select: true,
	})
	if err != nil {
		return nil, err
	}
	if len(tribes) == 0 {
		return nil, errors.Errorf("Tribe (ID: %d) not found.", id)
	}
	return tribes[0], nil
}

func (ucase *usecase) SearchTribe(ctx context.Context, cfg tribe.SearchTribeConfig) ([]*twmodel.FoundTribe, int, error) {
	if "" == strings.TrimSpace(cfg.Version) {
		return nil, 0, errors.New("Version is required.")
	}
	if "" == strings.TrimSpace(cfg.Query) {
		return nil, 0, errors.New("Query is too ambiguous. You must specify the variable 'query'.")
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > tribe.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = tribe.FetchLimit
	}
	if len(cfg.Sort) > tribe.MaxOrders {
		cfg.Sort = cfg.Sort[0:tribe.MaxOrders]
	}
	return ucase.repo.SearchTribe(ctx, cfg)
}
