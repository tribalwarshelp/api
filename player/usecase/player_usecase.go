package usecase

import (
	"context"
	"fmt"

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

func (ucase *usecase) Fetch(ctx context.Context, server string, filter *models.PlayerFilter) ([]*models.Player, int, error) {
	if filter == nil {
		filter = &models.PlayerFilter{}
	}
	if filter.Limit > player.PaginationLimit || filter.Limit <= 0 {
		filter.Limit = player.PaginationLimit
	}
	filter.Sort = utils.SanitizeSort(filter.Sort)
	if filter.TribeFilter != nil {
		filter.TribeFilter.Sort = utils.SanitizeSort(filter.TribeFilter.Sort)
	}
	return ucase.repo.Fetch(ctx, player.FetchConfig{
		Server: server,
		Filter: filter,
		Count:  true,
	})
}

func (ucase *usecase) GetByID(ctx context.Context, server string, id int) (*models.Player, error) {
	players, _, err := ucase.repo.Fetch(ctx, player.FetchConfig{
		Server: server,
		Filter: &models.PlayerFilter{
			ID:    []int{id},
			Limit: 1,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(players) == 0 {
		return nil, fmt.Errorf("Player (ID: %d) not found.", id)
	}
	return players[0], nil
}
