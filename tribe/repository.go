package tribe

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type FetchConfig struct {
	Server string
	Filter *twmodel.TribeFilter
	Count  bool
	Select bool
	Sort   []string
	Limit  int
	Offset int
}

type SearchTribeConfig struct {
	Version string
	Query   string
	Count   bool
	Sort    []string
	Limit   int
	Offset  int
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*twmodel.Tribe, int, error)
	SearchTribe(ctx context.Context, cfg SearchTribeConfig) ([]*twmodel.FoundTribe, int, error)
}
