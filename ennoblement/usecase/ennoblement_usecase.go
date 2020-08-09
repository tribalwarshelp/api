package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/ennoblement"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo ennoblement.Repository
}

func New(repo ennoblement.Repository) ennoblement.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, server string, filter *models.EnnoblementFilter) ([]*models.Ennoblement, int, error) {
	if filter == nil {
		filter = &models.EnnoblementFilter{}
	}
	if !middleware.MayExceedLimit(ctx) && (filter.Limit > ennoblement.PaginationLimit || filter.Limit <= 0) {
		filter.Limit = ennoblement.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	return ucase.repo.Fetch(ctx, ennoblement.FetchConfig{
		Server: server,
		Filter: filter,
		Count:  true,
	})
}
