package dataloaders

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
	"time"

	"github.com/tribalwarshelp/api/player"
	"github.com/tribalwarshelp/api/tribe"
	"github.com/tribalwarshelp/api/version"
	"github.com/tribalwarshelp/api/village"
)

const (
	wait = 4 * time.Millisecond
)

type DataLoaders struct {
	VersionByCode *VersionLoader
}

type Config struct {
	PlayerRepo  player.Repository
	TribeRepo   tribe.Repository
	VillageRepo village.Repository
	VersionRepo version.Repository
}

func NewDataLoaders(cfg Config) *DataLoaders {
	return &DataLoaders{
		VersionByCode: &VersionLoader{
			wait:     wait,
			maxBatch: 0,
			fetch: func(keys []string) ([]*twmodel.Version, []error) {
				var codes []twmodel.VersionCode
				for _, code := range keys {
					codes = append(codes, twmodel.VersionCode(code))
				}
				versions, _, err := cfg.VersionRepo.Fetch(context.Background(), version.FetchConfig{
					Filter: &twmodel.VersionFilter{
						Code: codes,
					},
					Select: true,
				})
				if err != nil {
					return nil, []error{err}
				}

				versionByCode := make(map[twmodel.VersionCode]*twmodel.Version)
				for _, v := range versions {
					versionByCode[v.Code] = v
				}

				inOrder := make([]*twmodel.Version, len(keys))
				for i, code := range codes {
					inOrder[i] = versionByCode[code]
				}

				return inOrder, nil
			},
		},
	}
}
