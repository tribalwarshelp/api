package dataloaders

import (
	"context"
	"time"

	"github.com/tribalwarshelp/shared/models"
)

type VersionDataLoaders struct {
	PlayerServersByID     PlayerServersLoader
	PlayerNameChangesByID PlayerNameChangesLoader
}

func NewVersionDataLoaders(versionCode models.VersionCode, cfg Config) *VersionDataLoaders {
	return &VersionDataLoaders{
		PlayerServersByID: PlayerServersLoader{
			wait:     2 * time.Millisecond,
			maxBatch: 0,
			fetch: func(keys []int) ([][]string, []error) {
				playerServersByID, err := cfg.PlayerRepo.FetchPlayerServers(context.Background(), versionCode, keys...)
				if err != nil {
					return nil, []error{err}
				}
				inOrder := make([][]string, len(keys))
				for i, id := range keys {
					inOrder[i] = playerServersByID[id]
				}
				return inOrder, nil
			},
		},
		PlayerNameChangesByID: PlayerNameChangesLoader{
			wait:     2 * time.Millisecond,
			maxBatch: 0,
			fetch: func(keys []int) ([][]*models.PlayerNameChange, []error) {
				playerNameChangesByID, err := cfg.PlayerRepo.FetchNameChanges(context.Background(), versionCode, keys...)
				if err != nil {
					return nil, []error{err}
				}
				inOrder := make([][]*models.PlayerNameChange, len(keys))
				for i, id := range keys {
					inOrder[i] = playerNameChangesByID[id]
				}
				return inOrder, nil
			},
		},
	}
}
