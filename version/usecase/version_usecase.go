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

func (ucase *usecase) Fetch(ctx context.Context, cfg version.FetchConfig) ([]*models.Version, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.VersionFilter{}
	}

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > version.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = version.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByCode(ctx context.Context, code models.VersionCode) (*models.Version, error) {
	versions, _, err := ucase.repo.Fetch(ctx, version.FetchConfig{
		Filter: &models.VersionFilter{
			Code: []models.VersionCode{code},
		},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("There is no version with code: %s.", code)
	}
	return versions[0], nil
}
