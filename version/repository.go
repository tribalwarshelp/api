package version

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type FetchConfig struct {
	Filter *twmodel.VersionFilter
	Select bool
	Count  bool
	Sort   []string
	Limit  int
	Offset int
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*twmodel.Version, int, error)
}
