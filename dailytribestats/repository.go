package dailytribestats

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type FetchConfig struct {
	Server string
	Filter *models.DailyTribeStatsFilter
	Count  bool
}

type Repository interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.DailyTribeStats, int, error)
}
