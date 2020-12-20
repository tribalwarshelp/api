package version

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, cfg FetchConfig) ([]*models.Version, int, error)
	GetByCode(ctx context.Context, code models.VersionCode) (*models.Version, error)
}