package usecase

import (
	"context"
	"fmt"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/utils"

	"github.com/tribalwarshelp/api/version"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo version.Repository
}

func New(repo version.Repository) version.Usecase {
	return &usecase{
		repo,
	}
}

func (ucase *usecase) Fetch(ctx context.Context, filter *models.VersionFilter) ([]*models.Version, int, error) {
	if filter == nil {
		filter = &models.VersionFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (filter.Limit > version.PaginationLimit || filter.Limit <= 0) {
		filter.Limit = version.PaginationLimit
	}
	if len(filter.Tag) > 0 {
		filter.Code = append(filter.Code, filter.Tag...)
	}
	if len(filter.TagNEQ) > 0 {
		filter.CodeNEQ = append(filter.Code, filter.TagNEQ...)
	}
	filter.Sort = utils.SanitizeSortExpression(filter.Sort)
	return ucase.repo.Fetch(ctx, version.FetchConfig{
		Filter: filter,
		Count:  true,
	})
}

func (ucase *usecase) GetByCode(ctx context.Context, code models.VersionCode) (*models.Version, error) {
	versions, _, err := ucase.repo.Fetch(ctx, version.FetchConfig{
		Filter: &models.VersionFilter{
			Code:  []models.VersionCode{code},
			Limit: 1,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("There is no version with code: %s.", code)
	}
	return versions[0], nil
}
