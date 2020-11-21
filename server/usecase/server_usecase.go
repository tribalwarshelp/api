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
	if cfg.Filter.Limit > 0 {
		cfg.Limit = cfg.Filter.Limit
	}
	if cfg.Filter.Offset > 0 {
		cfg.Offset = cfg.Filter.Offset
	}
	if cfg.Filter.Sort != "" {
		cfg.Sort = append(cfg.Sort, cfg.Filter.Sort)
	}

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > server.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = server.PaginationLimit
	}
	if len(cfg.Filter.LangVersionTag) > 0 {
		cfg.Filter.VersionCode = append(cfg.Filter.VersionCode, cfg.Filter.LangVersionTag...)
	}
	if len(cfg.Filter.LangVersionTagNEQ) > 0 {
		cfg.Filter.VersionCodeNEQ = append(cfg.Filter.VersionCode, cfg.Filter.LangVersionTagNEQ...)
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
