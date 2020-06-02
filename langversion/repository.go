package langversion

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Repository interface {
	Fetch(ctx context.Context, filter *models.LangVersionFilter) ([]*models.LangVersion, int, error)
}
