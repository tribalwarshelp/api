package village

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, server string, filter *models.VillageFilter) ([]*models.Village, int, error)
	GetByID(ctx context.Context, server string, id int) (*models.Village, error)
}
