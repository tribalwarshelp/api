package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/tribechange"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo tribechange.Repository
}

func New(repo tribechange.Repository) tribechange.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, server string, filter *models.TribeChangeFilter) ([]*models.TribeChange, int, error) {
	if filter == nil {
		filter = &models.TribeChangeFilter{}
	}
	if filter.Limit > tribechange.PaginationLimit || filter.Limit <= 0 {
		filter.Limit = tribechange.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	return ucase.repo.Fetch(ctx, tribechange.FetchConfig{
		Server: server,
		Filter: filter,
		Count:  true,
	})
}
