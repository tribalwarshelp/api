package tribe

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.Tribe, int, error)
	GetByID(ctx context.Context, server string, id int) (*models.Tribe, error)
}
