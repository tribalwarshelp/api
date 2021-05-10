package resolvers

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/server"

	"github.com/tribalwarshelp/api/graphql/generated"
)

func (r *serverResolver) Version(ctx context.Context, obj *twmodel.Server) (*twmodel.Version, error) {
	loaders := middleware.DataLoaderFromContext(ctx)
	if loaders != nil {
		lv, _ := loaders.VersionByCode.Load(obj.VersionCode.String())
		return lv, nil
	}
	return nil, nil
}

func (r *queryResolver) Servers(ctx context.Context,
	f *twmodel.ServerFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.ServerList, error) {
	var err error
	list := &generated.ServerList{}
	list.Items, list.Total, err = r.ServerUcase.Fetch(ctx, server.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  safeptr.SafeIntPointer(limit, 0),
		Offset: safeptr.SafeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
		Select: shouldSelectItems(ctx),
	})
	return list, err
}

func (r *queryResolver) Server(ctx context.Context, key string) (*twmodel.Server, error) {
	return r.ServerUcase.GetByKey(ctx, key)
}
