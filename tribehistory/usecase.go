package tribehistory

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, server string, filter *models.TribeHistoryFilter) ([]*models.TribeHistory, int, error)
}
