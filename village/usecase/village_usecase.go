package usecase

import (
	"context"
	"fmt"

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
	if filter.Limit > village.PaginationLimit || filter.Limit <= 0 {
		filter.Limit = village.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	return ucase.repo.Fetch(ctx, server, filter)
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*models.Village, error) {
	villages, total, err := ucase.repo.Fetch(ctx, server, &models.VillageFilter{
		ID:    []int{id},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return nil, fmt.Errorf("Village with id: %d not found.", id)
	}
	return villages[0], nil
}
