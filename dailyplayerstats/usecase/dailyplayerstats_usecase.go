package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/dailyplayerstats"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo dailyplayerstats.Repository
}

func New(repo dailyplayerstats.Repository) dailyplayerstats.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, server string, filter *models.DailyPlayerStatsFilter) ([]*models.DailyPlayerStats, int, error) {
	if filter == nil {
		filter = &models.DailyPlayerStatsFilter{}
	}
	if filter.Limit > dailyplayerstats.PaginationLimit || filter.Limit <= 0 {
		filter.Limit = dailyplayerstats.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	return ucase.repo.Fetch(ctx, dailyplayerstats.FetchConfig{
		Server: server,
		Filter: filter,
		Count:  true,
	})
}
