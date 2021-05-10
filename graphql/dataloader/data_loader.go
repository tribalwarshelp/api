package dataloader

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
	wait     = 2 * time.Millisecond
	maxBatch = 200
)

type DataLoader struct {
	VersionByCode *VersionLoader
}

type Config struct {
	PlayerRepo  player.Repository
	TribeRepo   tribe.Repository
	VillageRepo village.Repository
	VersionRepo version.Repository
}

func NewDataLoader(cfg Config) *DataLoader {
	return &DataLoader{
		VersionByCode: &VersionLoader{
			wait:     wait,
			maxBatch: maxBatch,
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
