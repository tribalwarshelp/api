package usecase

import (
	"context"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/version"
)

type usecase struct {
	repo version.Repository
}

func New(repo version.Repository) version.Usecase {
	return &usecase{
		repo,
	}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg version.FetchConfig) ([]*twmodel.Version, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.VersionFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > version.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = version.FetchLimit
	}
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByCode(ctx context.Context, code twmodel.VersionCode) (*twmodel.Version, error) {
	versions, _, err := ucase.repo.Fetch(ctx, version.FetchConfig{
		Filter: &twmodel.VersionFilter{
			Code: []twmodel.VersionCode{code},
		},
		Limit:  1,
		Select: true,
		Count:  false,
	})
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, errors.Errorf("version (Code: %s) not found", code)
	}
	return versions[0], nil
}
