package dataloader

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type VersionDataLoader struct {
	PlayerServersByID     *PlayerServersLoader
	PlayerNameChangesByID *PlayerNameChangesLoader
}

func NewVersionDataLoader(versionCode twmodel.VersionCode, cfg Config) *VersionDataLoader {
	return &VersionDataLoader{
		PlayerServersByID: &PlayerServersLoader{
			wait:     wait,
			maxBatch: maxBatch,
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
			maxBatch: maxBatch,
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
