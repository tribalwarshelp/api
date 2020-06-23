package dataloaders

import (
	"context"
	"time"

	"github.com/tribalwarshelp/api/langversion"
	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/shared/models"
)

type DataLoaders struct {
	LangVersionByTag  LangVersionLoader
	PlayerServersByID PlayerServersLoader
}

type Config struct {
	PlayerRepo      player.Repository
	TribeRepo       tribe.Repository
	VillageRepo     village.Repository
	LangVersionRepo langversion.Repository
}

func New(cfg Config) *DataLoaders {
	return &DataLoaders{
		LangVersionByTag: LangVersionLoader{
			wait:     2 * time.Millisecond,
			maxBatch: 0,
			fetch: func(keys []string) ([]*models.LangVersion, []error) {
				tags := []models.LanguageTag{}
				for _, tag := range keys {
					tags = append(tags, models.LanguageTag(tag))
				}
				langVersions, _, err := cfg.LangVersionRepo.Fetch(context.Background(), &models.LangVersionFilter{
					Tag: tags,
				})
				if err != nil {
					return nil, []error{err}
				}

				langVersionByTag := make(map[models.LanguageTag]*models.LangVersion)
				for _, langVersion := range langVersions {
					langVersionByTag[langVersion.Tag] = langVersion
				}

				inOrder := make([]*models.LangVersion, len(keys))
				for i, tag := range tags {
					inOrder[i] = langVersionByTag[tag]
				}

				return inOrder, nil
			},
		},
		PlayerServersByID: PlayerServersLoader{
			wait:     2 * time.Millisecond,
			maxBatch: 0,
			fetch: func(keys []int) ([][]string, []error) {
				playerServersByID, err := cfg.PlayerRepo.FetchPlayerServers(context.Background(), keys...)
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
	}
}
