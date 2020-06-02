package langversion

import (
	"context"

	"github.com/tribalwarshelp/shared/models"
)

type Usecase interface {
	Fetch(ctx context.Context, filter *models.LangVersionFilter) ([]*models.LangVersion, int, error)
	GetByTag(ctx context.Context, tag models.LanguageTag) (*models.LangVersion, error)
}
