package dailytribestats

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, server string, filter *models.DailyTribeStatsFilter) ([]*models.DailyTribeStats, int, error)
}
