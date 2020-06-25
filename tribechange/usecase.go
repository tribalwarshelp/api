package tribechange

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, server string, filter *models.TribeChangeFilter) ([]*models.TribeChange, int, error)
}
