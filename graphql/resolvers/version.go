package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *versionResolver) Tag(ctx context.Context, obj *models.Version) (models.VersionCode, error) {
	return obj.Code, nil
}

func (r *queryResolver) Versions(ctx context.Context, filter *models.VersionFilter) (*generated.VersionList, error) {
	var err error
	list := &generated.VersionList{}
	list.Items, list.Total, err = r.VersionUcase.Fetch(ctx, filter)
	return list, err
}

func (r *queryResolver) Version(ctx context.Context, code models.VersionCode) (*models.Version, error) {
	return r.VersionUcase.GetByCode(ctx, code)
}

func (r *queryResolver) LangVersions(ctx context.Context, filter *models.VersionFilter) (*generated.VersionList, error) {
	return r.Versions(ctx, filter)
}

func (r *queryResolver) LangVersion(ctx context.Context, tag models.VersionCode) (*models.Version, error) {
	return r.Version(ctx, tag)
}
