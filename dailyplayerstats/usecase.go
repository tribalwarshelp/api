package dailyplayerstats

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, server string, filter *models.DailyPlayerStatsFilter) ([]*models.DailyPlayerStats, int, error)
}
