package playerhistory

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, server string, filter *models.PlayerHistoryFilter) ([]*models.PlayerHistory, int, error)
}
