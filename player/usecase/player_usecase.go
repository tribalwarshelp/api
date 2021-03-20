package usecase

import (
	"context"
	"fmt"
	"strings"

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

	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > player.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = player.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
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

func (ucase *usecase) SearchPlayer(ctx context.Context, cfg player.SearchPlayerConfig) ([]*models.FoundPlayer, int, error) {
	if "" == strings.TrimSpace(cfg.Version) {
		return nil, 0, fmt.Errorf("Version is required.")
	}
	if "" == strings.TrimSpace(cfg.Name) && cfg.ID <= 0 {
		return nil, 0, fmt.Errorf("Your search is ambiguous. You must specify the variable 'name' or 'id'.")
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > player.PaginationLimit || cfg.Limit <= 0) {
		cfg.Limit = player.PaginationLimit
	}
	cfg.Sort = utils.SanitizeSorts(cfg.Sort)
	return ucase.repo.SearchPlayer(ctx, cfg)
}
