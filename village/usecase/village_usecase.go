package usecase

import (
	"context"
	"fmt"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo village.Repository
}

func New(repo village.Repository) village.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, server string, filter *models.VillageFilter) ([]*models.Village, int, error) {
	if filter == nil {
		filter = &models.VillageFilter{}
	}
	if !middleware.MayExceedLimit(ctx) && (filter.Limit > village.PaginationLimit || filter.Limit <= 0) {
		filter.Limit = village.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	if filter.PlayerFilter != nil {
		filter.PlayerFilter.Sort = utils.SanitizeSort(filter.PlayerFilter.Sort)
		if filter.PlayerFilter.TribeFilter != nil {
			filter.PlayerFilter.TribeFilter.Sort = utils.SanitizeSort(filter.PlayerFilter.TribeFilter.Sort)
		}
	}
	return ucase.repo.Fetch(ctx, village.FetchConfig{
		Server: server,
		Count:  true,
		Filter: filter,
	})
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*models.Village, error) {
	villages, _, err := ucase.repo.Fetch(ctx, village.FetchConfig{
		Filter: &models.VillageFilter{
			ID:    []int{id},
			Limit: 1,
		},
		Server: server,
	})
	if err != nil {
		return nil, err
	}
	if len(villages) == 0 {
		return nil, fmt.Errorf("Village (ID: %d) not found.", id)
	}
	return villages[0], nil
}
