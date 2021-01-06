package usecase

import (
	"context"
	"fmt"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/server"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo server.Repository
}

func New(repo server.Repository) server.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg server.FetchConfig) ([]*models.Server, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.ServerFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > server.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = server.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSortExpressions(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByKey(ctx context.Context, key string) (*models.Server, error) {
	servers, _, err := ucase.repo.Fetch(ctx, server.FetchConfig{
		Filter: &models.ServerFilter{
			Key: []string{key},
		},
		Limit: 1,
		Count: false,
	})
	if err != nil {
		return nil, err
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("Server (key: %s) not found.", key)
	}
	return servers[0], nil
}
