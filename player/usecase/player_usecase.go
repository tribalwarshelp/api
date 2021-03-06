package usecase

import (
	"context"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/tw/twmodel"
	"strings"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/player"
)

type usecase struct {
	repo player.Repository
}

func New(repo player.Repository) player.Usecase {
	return &usecase{repo}
}

func (ucase *usecase) Fetch(ctx context.Context, cfg player.FetchConfig) ([]*twmodel.Player, int, error) {
	if cfg.Filter == nil {
		cfg.Filter = &twmodel.PlayerFilter{}
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > player.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = player.FetchLimit
	}
	if len(cfg.Sort) > player.MaxOrders {
		cfg.Sort = cfg.Sort[0:player.MaxOrders]
	}
	return ucase.repo.Fetch(ctx, cfg)
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*twmodel.Player, error) {
	players, _, err := ucase.repo.Fetch(ctx, player.FetchConfig{
		Server: server,
		Filter: &twmodel.PlayerFilter{
			ID: []int{id},
		},
		Limit:  1,
		Count:  false,
		Select: true,
	})
	if err != nil {
		return nil, err
	}
	if len(players) == 0 {
		return nil, errors.Errorf("Player (ID: %d) not found.", id)
	}
	return players[0], nil
}

func (ucase *usecase) SearchPlayer(ctx context.Context, cfg player.SearchPlayerConfig) ([]*twmodel.FoundPlayer, int, error) {
	if "" == strings.TrimSpace(cfg.Version) {
		return nil, 0, errors.New("Version is required.")
	}
	if "" == strings.TrimSpace(cfg.Name) && cfg.ID <= 0 {
		return nil, 0, errors.New("Query is too ambiguous. You must specify the variable 'name' or 'id'.")
	}
	if !middleware.CanExceedLimit(ctx) && (cfg.Limit > player.FetchLimit || cfg.Limit <= 0) {
		cfg.Limit = player.FetchLimit
	}
	if len(cfg.Sort) > player.MaxOrders {
		cfg.Sort = cfg.Sort[0:player.MaxOrders]
	}
	return ucase.repo.SearchPlayer(ctx, cfg)
}
