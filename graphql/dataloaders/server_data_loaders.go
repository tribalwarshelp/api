package dataloaders

import (
	"context"
	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/shared/models"
)

type ServerDataLoaders struct {
	PlayerByID  *PlayerLoader
	TribeByID   *TribeLoader
	VillageByID *VillageLoader
}

func NewServerDataLoaders(server string, cfg Config) *ServerDataLoaders {
	return &ServerDataLoaders{
		PlayerByID: &PlayerLoader{
			wait:     wait,
			maxBatch: 0,
			fetch: func(ids []int) ([]*models.Player, []error) {
				players, _, err := cfg.PlayerRepo.Fetch(context.Background(), player.FetchConfig{
					Filter: &models.PlayerFilter{
						ID: ids,
					},
					Select: true,
					Server: server,
				})
				if err != nil {
					return nil, []error{err}
				}

				playerByID := make(map[int]*models.Player)
				for _, player := range players {
					playerByID[player.ID] = player
				}

				inOrder := make([]*models.Player, len(ids))
				for i, id := range ids {
					inOrder[i] = playerByID[id]
				}

				return inOrder, nil
			},
		},
		TribeByID: &TribeLoader{
			wait:     wait,
			maxBatch: 0,
			fetch: func(ids []int) ([]*models.Tribe, []error) {
				tribes, _, err := cfg.TribeRepo.Fetch(context.Background(), tribe.FetchConfig{
					Server: server,
					Filter: &models.TribeFilter{
						ID: ids,
					},
					Select: true,
				})
				if err != nil {
					return nil, []error{err}
				}

				tribeByID := make(map[int]*models.Tribe)
				for _, tribe := range tribes {
					tribeByID[tribe.ID] = tribe
				}

				inOrder := make([]*models.Tribe, len(ids))
				for i, id := range ids {
					inOrder[i] = tribeByID[id]
				}

				return inOrder, nil
			},
		},
		VillageByID: &VillageLoader{
			wait:     wait,
			maxBatch: 0,
			fetch: func(ids []int) ([]*models.Village, []error) {
				villages, _, err := cfg.VillageRepo.Fetch(context.Background(), village.FetchConfig{
					Server: server,
					Count:  false,
					Filter: &models.VillageFilter{
						ID: ids,
					},
					Select: true,
				})
				if err != nil {
					return nil, []error{err}
				}

				villageByID := make(map[int]*models.Village)
				for _, village := range villages {
					villageByID[village.ID] = village
				}

				inOrder := make([]*models.Village, len(ids))
				for i, id := range ids {
					inOrder[i] = villageByID[id]
				}

				return inOrder, nil
			},
		},
	}
}
