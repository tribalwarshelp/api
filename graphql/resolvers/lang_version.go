package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *queryResolver) LangVersions(ctx context.Context, filter *models.LangVersionFilter) (*generated.LangVersionList, error) {
	var err error
	list := &generated.LangVersionList{}
	list.Items, list.Total, err = r.LangVersionUcase.Fetch(ctx, filter)
	return list, err
}

func (r *queryResolver) LangVersion(ctx context.Context, tag models.LanguageTag) (*models.LangVersion, error) {
	return r.LangVersionUcase.GetByTag(ctx, tag)
}
