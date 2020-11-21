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
	if cfg.Filter.Limit > 0 {
		cfg.Limit = cfg.Filter.Limit
	}
	if cfg.Filter.Offset > 0 {
		cfg.Offset = cfg.Filter.Offset
	}
	if cfg.Filter.Sort != "" {
		cfg.Sort = append(cfg.Sort, cfg.Filter.Sort)
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > version.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = version.PaginationLimit
	}
	if len(cfg.Filter.Tag) > 0 {
		cfg.Filter.Code = append(cfg.Filter.Code, cfg.Filter.Tag...)
	}
	if len(cfg.Filter.TagNEQ) > 0 {
		cfg.Filter.CodeNEQ = append(cfg.Filter.Code, cfg.Filter.TagNEQ...)
	}
	cfg.Sort = utils.SanitizeSortExpressions(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
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
