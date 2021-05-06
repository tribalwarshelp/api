package usecase

import (
	"context"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/server"
)

type usecase struct {
	repo server.Repository
}

func New(repo server.Repository) server.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg server.FetchConfig) ([]*twmodel.Server, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.ServerFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > server.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = server.FetchLimit
	}
	if len(cfg.Sort) > server.MaxOrders {
		cfg.Sort = cfg.Sort[0:server.MaxOrders]
	}
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByKey(ctx context.Context, key string) (*twmodel.Server, error) {
	servers, _, err := ucase.repo.Fetch(ctx, server.FetchConfig{
		Filter: &twmodel.ServerFilter{
			Key: []string{key},
		},
		Limit:  1,
		Count:  false,
		Select: true,
	})
	if err != nil {
		return nil, err
	}
	if len(servers) == 0 {
		return nil, errors.Errorf("Server (key: %s) not found.", key)
	}
	return servers[0], nil
}
