package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/playerhistory"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo playerhistory.Repository
}

func New(repo playerhistory.Repository) playerhistory.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, server string, filter *models.PlayerHistoryFilter) ([]*models.PlayerHistory, int, error) {
	if filter == nil {
		filter = &models.PlayerHistoryFilter{}
	}
	if !middleware.MayExceedLimit(ctx) && (filter.Limit > playerhistory.PaginationLimit || filter.Limit <= 0) {
		filter.Limit = playerhistory.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	return ucase.repo.Fetch(ctx, playerhistory.FetchConfig{
		Server: server,
		Filter: filter,
		Count:  true,
	})
}
