package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/serverstats"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo serverstats.Repository
}

func New(repo serverstats.Repository) serverstats.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, server string, filter *models.ServerStatsFilter) ([]*models.ServerStats, int, error) {
	if filter == nil {
		filter = &models.ServerStatsFilter{}
	}
	if filter.Limit > serverstats.PaginationLimit || filter.Limit <= 0 {
		filter.Limit = serverstats.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	return ucase.repo.Fetch(ctx, server, filter)
}
