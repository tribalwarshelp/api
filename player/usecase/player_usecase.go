package usecase

import (
	"context"
	"fmt"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

type usecase struct {
	repo player.Repository
}

func New(repo player.Repository) player.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg player.FetchConfig) ([]*models.Player, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &models.PlayerFilter{}
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

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > player.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = player.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSortExpressions(cfg.Sort)
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*models.Player, error) {
	players, _, err := ucase.repo.Fetch(ctx, player.FetchConfig{
		Server: server,
		Filter: &models.PlayerFilter{
			ID: []int{id},
		},
		Limit: 1,
		Count: false,
	})
	if err != nil {
		return nil, err
	}
	if len(players) == 0 {
		return nil, fmt.Errorf("Player (ID: %d) not found.", id)
	}
	return players[0], nil
}
