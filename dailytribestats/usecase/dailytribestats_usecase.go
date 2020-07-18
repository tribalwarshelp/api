package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/dailytribestats"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo dailytribestats.Repository
}

func New(repo dailytribestats.Repository) dailytribestats.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, server string, filter *models.DailyTribeStatsFilter) ([]*models.DailyTribeStats, int, error) {
	if filter == nil {
		filter = &models.DailyTribeStatsFilter{}
	}
	if filter.Limit > dailytribestats.PaginationLimit || filter.Limit <= 0 {
		filter.Limit = dailytribestats.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	if filter.TribeFilter != nil {
		filter.TribeFilter.Sort = utils.SanitizeSort(filter.TribeFilter.Sort)
	}
	return ucase.repo.Fetch(ctx, dailytribestats.FetchConfig{
		Server: server,
		Filter: filter,
		Count:  true,
	})
}
