package dataloaders

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type VersionDataLoaders struct {
	PlayerServersByID     *PlayerServersLoader
	PlayerNameChangesByID *PlayerNameChangesLoader
}

func NewVersionDataLoaders(versionCode twmodel.VersionCode, cfg Config) *VersionDataLoaders {
	return &VersionDataLoaders{
		PlayerServersByID: &PlayerServersLoader{
			wait:     wait,
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
		PlayerNameChangesByID: &PlayerNameChangesLoader{
			wait:     wait,
			maxBatch: 0,
			fetch: func(keys []int) ([][]*twmodel.PlayerNameChange, []error) {
				playerNameChangesByID, err := cfg.PlayerRepo.FetchNameChanges(context.Background(), versionCode, keys...)
				if err != nil {
					return nil, []error{err}
				}
				inOrder := make([][]*twmodel.PlayerNameChange, len(keys))
				for i, id := range keys {
					inOrder[i] = playerNameChangesByID[id]
				}
				return inOrder, nil
			},
		},
	}
}
