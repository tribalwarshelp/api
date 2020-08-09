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

func (ucase *usecase) Fetch(ctx context.Context, filter *models.ServerFilter) ([]*models.Server, int, error) {
	if filter == nil {
		filter = &models.ServerFilter{}
	}
	if !middleware.MayExceedLimit(ctx) && (filter.Limit > server.PaginationLimit || filter.Limit <= 0) {
		filter.Limit = server.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	return ucase.repo.Fetch(ctx, server.FetchConfig{
		Count:  true,
		Filter: filter,
	})
}

func (ucase *usecase) GetByKey(ctx context.Context, key string) (*models.Server, error) {
	servers, _, err := ucase.repo.Fetch(ctx, server.FetchConfig{
		Filter: &models.ServerFilter{
			Key:   []string{key},
			Limit: 1,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("Server (key: %s) not found.", key)
	}
	return servers[0], nil
}
