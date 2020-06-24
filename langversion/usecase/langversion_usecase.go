package usecase

import (
	"context"
	"fmt"

	"github.com/tribalwarshelp/api/utils"

	"github.com/tribalwarshelp/api/langversion"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo langversion.Repository
}

func New(repo langversion.Repository) langversion.Usecase {
	return &usecase{
		repo,
	}
}

func (ucase *usecase) Fetch(ctx context.Context, filter *models.LangVersionFilter) ([]*models.LangVersion, int, error) {
	if filter == nil {
		filter = &models.LangVersionFilter{}
	}
	if filter.Limit > langversion.PaginationLimit || filter.Limit <= 0 {
		filter.Limit = langversion.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	return ucase.repo.Fetch(ctx, langversion.FetchConfig{
		Filter: filter,
		Count:  true,
	})
}

func (ucase *usecase) GetByTag(ctx context.Context, tag models.LanguageTag) (*models.LangVersion, error) {
	langversions, _, err := ucase.repo.Fetch(ctx, langversion.FetchConfig{
		Filter: &models.LangVersionFilter{
			Tag:   []models.LanguageTag{tag},
			Limit: 1,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(langversions) == 0 {
		return nil, fmt.Errorf("There is no lang version with tag: %s.", tag)
	}
	return langversions[0], nil
}
