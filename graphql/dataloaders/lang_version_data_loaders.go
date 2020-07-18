package dataloaders

import (
	"context"
	"time"

	"github.com/tribalwarshelp/shared/models"
)

type LangVersionDataLoaders struct {
	PlayerServersByID     PlayerServersLoader
	PlayerNameChangesByID PlayerNameChangesLoader
}

func NewLangVersionDataLoaders(langTag models.LanguageTag, cfg Config) *LangVersionDataLoaders {
	return &LangVersionDataLoaders{
		PlayerServersByID: PlayerServersLoader{
			wait:     2 * time.Millisecond,
			maxBatch: 0,
			fetch: func(keys []int) ([][]string, []error) {
				playerServersByID, err := cfg.PlayerRepo.FetchPlayerServers(context.Background(), langTag, keys...)
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
				playerNameChangesByID, err := cfg.PlayerRepo.FetchNameChanges(context.Background(), langTag, keys...)
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