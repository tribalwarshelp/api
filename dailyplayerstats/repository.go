package dailyplayerstats

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"
)

type FetchConfig struct {
	Server string
	Filter *twmodel.DailyPlayerStatsFilter
	Select bool
	Count  bool
	Sort   []string
	Limit  int
	Offset int
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*twmodel.DailyPlayerStats, int, error)
}
