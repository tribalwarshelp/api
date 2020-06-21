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
	LangVersionByTag LangVersionLoader
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
			fetch: func(tagsS []string) ([]*models.LangVersion, []error) {
				tags := []models.LanguageTag{}
				for _, tag := range tagsS {
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

				inOrder := make([]*models.LangVersion, len(tagsS))
				for i, tag := range tags {
					inOrder[i] = langVersionByTag[tag]
				}

				return inOrder, nil
			},
		},
	}
}