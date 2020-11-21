package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/version"
	"github.com/tribalwarshelp/shared/models"
)

func (r *versionResolver) Tag(ctx context.Context, obj *models.Version) (models.VersionCode, error) {
	return obj.Code, nil
}

func (r *queryResolver) Versions(ctx context.Context,
	f *models.VersionFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.VersionList, error) {
	defLimit := 0
	defOffset := 0
	if limit == nil {
		limit = &defLimit
	}
	if offset == nil {
		offset = &defOffset
	}

	var err error
	list := &generated.VersionList{}
	list.Items, list.Total, err = r.VersionUcase.Fetch(ctx, version.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  *limit,
		Offset: *offset,
		Count:  true,
	})
	return list, err
}

func (r *queryResolver) Version(ctx context.Context, code models.VersionCode) (*models.Version, error) {
	return r.VersionUcase.GetByCode(ctx, code)
}

func (r *queryResolver) LangVersions(ctx context.Context, filter *models.VersionFilter) (*generated.VersionList, error) {
	return r.Versions(ctx, filter, nil, nil, []string{})
}

func (r *queryResolver) LangVersion(ctx context.Context, tag models.VersionCode) (*models.Version, error) {
	return r.Version(ctx, tag)
}
