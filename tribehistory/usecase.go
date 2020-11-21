package tribehistory

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.TribeHistory, int, error)
}
