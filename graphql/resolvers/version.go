package resolvers

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/version"
)

func (r *queryResolver) Versions(ctx context.Context,
	f *twmodel.VersionFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.VersionList, error) {
	var err error
	list := &generated.VersionList{}
	list.Items, list.Total, err = r.VersionUcase.Fetch(ctx, version.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  safeptr.SafeIntPointer(limit, 0),
		Offset: safeptr.SafeIntPointer(offset, 0),
		Select: shouldSelectItems(ctx),
		Count:  shouldCount(ctx),
	})
	return list, err
}

func (r *queryResolver) Version(ctx context.Context, code twmodel.VersionCode) (*twmodel.Version, error) {
	return r.VersionUcase.GetByCode(ctx, code)
}
