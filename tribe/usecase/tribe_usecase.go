package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/tribalwarshelp/api/middleware"
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

func (ucase *usecase) Fetch(ctx context.Context, cfg tribe.FetchConfig) ([]*models.Tribe, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.TribeFilter{}
	}

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > tribe.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = tribe.FetchLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*models.Tribe, error) {
	tribes, _, err := ucase.repo.Fetch(ctx, tribe.FetchConfig{
		Filter: &models.TribeFilter{
			ID: []int{id},
		},
		Limit:  1,
		Server: server,
		Count:  false,
		Select: true,
	})
	if err != nil {
		return nil, err
	}
	if len(tribes) == 0 {
		return nil, fmt.Errorf("Tribe (ID: %s) not found.", id)
	}
	return tribes[0], nil
}

func (ucase *usecase) SearchTribe(ctx context.Context, cfg tribe.SearchTribeConfig) ([]*models.FoundTribe, int, error) {
	if "" == strings.TrimSpace(cfg.Version) {
		return nil, 0, fmt.Errorf("Version is required.")
	}
	if "" == strings.TrimSpace(cfg.Query) {
		return nil, 0, fmt.Errorf("Your search is ambiguous. You must specify the variable 'query'.")
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > tribe.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = tribe.FetchLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
	return ucase.repo.SearchTribe(ctx, cfg)
}
