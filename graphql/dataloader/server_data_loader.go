package dataloader

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/api/village"
)

type ServerDataLoader struct {
	PlayerByID  *PlayerLoader
	TribeByID   *TribeLoader
	VillageByID *VillageLoader
}

func NewServerDataLoader(server string, cfg Config) *ServerDataLoader {
	return &ServerDataLoader{
		PlayerByID: &PlayerLoader{
			wait:     wait,
			maxBatch: maxBatch,
			fetch: func(ids []int) ([]*twmodel.Player, []error) {
				players, _, err := cfg.PlayerRepo.Fetch(context.Background(), player.FetchConfig{
					Filter: &twmodel.PlayerFilter{
						ID: ids,
					},
					Select: true,
					Server: server,
				})
				if err != nil {
					return nil, []error{err}
				}

				playerByID := make(map[int]*twmodel.Player)
				for _, p := range players {
					playerByID[p.ID] = p
				}

				inOrder := make([]*twmodel.Player, len(ids))
				for i, id := range ids {
					inOrder[i] = playerByID[id]
				}

				return inOrder, nil
			},
		},
		TribeByID: &TribeLoader{
			wait:     wait,
			maxBatch: maxBatch,
			fetch: func(ids []int) ([]*twmodel.Tribe, []error) {
				tribes, _, err := cfg.TribeRepo.Fetch(context.Background(), tribe.FetchConfig{
					Server: server,
					Filter: &twmodel.TribeFilter{
						ID: ids,
					},
					Select: true,
				})
				if err != nil {
					return nil, []error{err}
				}

				tribeByID := make(map[int]*twmodel.Tribe)
				for _, t := range tribes {
					tribeByID[t.ID] = t
				}

				inOrder := make([]*twmodel.Tribe, len(ids))
				for i, id := range ids {
					inOrder[i] = tribeByID[id]
				}

				return inOrder, nil
			},
		},
		VillageByID: &VillageLoader{
			wait:     wait,
			maxBatch: maxBatch,
			fetch: func(ids []int) ([]*twmodel.Village, []error) {
				villages, _, err := cfg.VillageRepo.Fetch(context.Background(), village.FetchConfig{
					Server: server,
					Count:  false,
					Filter: &twmodel.VillageFilter{
						ID: ids,
					},
					Select: true,
				})
				if err != nil {
					return nil, []error{err}
				}

				villageByID := make(map[int]*twmodel.Village)
				for _, v := range villages {
					villageByID[v.ID] = v
				}

				inOrder := make([]*twmodel.Village, len(ids))
				for i, id := range ids {
					inOrder[i] = villageByID[id]
				}

				return inOrder, nil
			},
		},
	}
}
