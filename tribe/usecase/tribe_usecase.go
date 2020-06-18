package usecase

import (
	"context"
	"fmt"

	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo tribe.Repository
}

func New(repo tribe.Repository) tribe.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, server string, filter *models.TribeFilter) ([]*models.Tribe, int, error) {
	if filter == nil {
		filter = &models.TribeFilter{}
	}
	if filter.Limit > tribe.PaginationLimit || filter.Limit <= 0 {
		filter.Limit = tribe.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	return ucase.repo.Fetch(ctx, server, filter)
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*models.Tribe, error) {
	tribes, total, err := ucase.repo.Fetch(ctx, server, &models.TribeFilter{
		ID:    []int{id},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return nil, fmt.Errorf("Tribe (ID: %s) not found.", id)
	}
	return tribes[0], nil
}
