package dataloaders

import (
	"context"
	"time"

	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/api/version"
	"github.com/tribalwarshelp/api/village"
	"github.com/tribalwarshelp/shared/models"
)

type DataLoaders struct {
	VersionByTag VersionLoader
}

type Config struct {
	PlayerRepo  player.Repository
	TribeRepo   tribe.Repository
	VillageRepo village.Repository
	VersionRepo version.Repository
}

func NewDataLoaders(cfg Config) *DataLoaders {
	return &DataLoaders{
		VersionByTag: VersionLoader{
			wait:     2 * time.Millisecond,
			maxBatch: 0,
			fetch: func(keys []string) ([]*models.Version, []error) {
				codes := []models.VersionCode{}
				for _, code := range keys {
					codes = append(codes, models.VersionCode(code))
				}
				versions, _, err := cfg.VersionRepo.Fetch(context.Background(), version.FetchConfig{
					Filter: &models.VersionFilter{
						Code: codes,
					},
				})
				if err != nil {
					return nil, []error{err}
				}

				versionByCode := make(map[models.VersionCode]*models.Version)
				for _, version := range versions {
					versionByCode[version.Code] = version
				}

				inOrder := make([]*models.Version, len(keys))
				for i, code := range codes {
					inOrder[i] = versionByCode[code]
				}

				return inOrder, nil
			},
		},
	}
}
