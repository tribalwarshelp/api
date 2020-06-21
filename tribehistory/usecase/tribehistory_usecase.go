package usecase

import (
	"context"

	"github.com/tribalwarshelp/api/tribehistory"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo tribehistory.Repository
}

func New(repo tribehistory.Repository) tribehistory.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, server string, filter *models.TribeHistoryFilter) ([]*models.TribeHistory, int, error) {
	if filter == nil {
		filter = &models.TribeHistoryFilter{}
	}
	if filter.Limit > tribehistory.PaginationLimit || filter.Limit <= 0 {
		filter.Limit = tribehistory.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	return ucase.repo.Fetch(ctx, server, filter)
}
