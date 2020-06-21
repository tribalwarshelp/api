package serverstats

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, server string, filter *models.ServerStatsFilter) ([]*models.ServerStats, int, error)
}
