package dataloaders

import (
	"context"
	"log"
	"time"

	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/shared/models"
)

type DataLoaders struct {
	PlayerByID PlayerLoader
	TribeByID  TribeLoader
}

type Config struct {
	PlayerRepo player.Repository
	TribeRepo  tribe.Repository
}

func New(server string, cfg Config) *DataLoaders {
	return &DataLoaders{
		PlayerByID: PlayerLoader{
			wait:     2 * time.Millisecond,
			maxBatch: 0,
			fetch: func(ids []int) ([]*models.Player, []error) {
				log.Println("playerbyid", ids)
				players, _, err := cfg.PlayerRepo.Fetch(context.Background(), server, &models.PlayerFilter{
					ID: ids,
				})
				if err != nil {
					return nil, []error{err}
				}

				playerByID := make(map[int]*models.Player)
				for _, player := range players {
					playerByID[player.ID] = player
				}

				sorted := make([]*models.Player, len(ids))
				for i, id := range ids {
					sorted[i] = playerByID[id]
				}

				return sorted, nil
			},
		},
		TribeByID: TribeLoader{
			wait:     2 * time.Millisecond,
			maxBatch: 0,
			fetch: func(ids []int) ([]*models.Tribe, []error) {
				log.Println("tribebyid", ids)
				tribes, _, err := cfg.TribeRepo.Fetch(context.Background(), server, &models.TribeFilter{
					ID: ids,
				})
				if err != nil {
					return nil, []error{err}
				}

				tribeByID := make(map[int]*models.Tribe)
				for _, tribe := range tribes {
					tribeByID[tribe.ID] = tribe
				}

				sorted := make([]*models.Tribe, len(ids))
				for i, id := range ids {
					sorted[i] = tribeByID[id]
				}

				return sorted, nil
			},
		},
	}
}
