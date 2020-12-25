package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/version"
	"github.com/tribalwarshelp/shared/models"
)

func (r *queryResolver) Versions(ctx context.Context,
	f *models.VersionFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.VersionList, error) {
	var err error
	list := &generated.VersionList{}
	list.Items, list.Total, err = r.VersionUcase.Fetch(ctx, version.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  safeIntPointer(limit, 0),
		Offset: safeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
	})
	return list, err
}

func (r *queryResolver) Version(ctx context.Context, code models.VersionCode) (*models.Version, error) {
	return r.VersionUcase.GetByCode(ctx, code)
}
