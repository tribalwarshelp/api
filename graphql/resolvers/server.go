package resolvers

import (
	"context"
	"github.com/tribalwarshelp/api/utils"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/server"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *serverResolver) Version(ctx context.Context, obj *models.Server) (*models.Version, error) {
	loaders := middleware.DataLoadersFromContext(ctx)
	if loaders != nil {
		lv, _ := loaders.VersionByCode.Load(obj.VersionCode.String())
		return lv, nil
	}
	return nil, nil
}

func (r *queryResolver) Servers(ctx context.Context,
	f *models.ServerFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.ServerList, error) {
	var err error
	list := &generated.ServerList{}
	list.Items, list.Total, err = r.ServerUcase.Fetch(ctx, server.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  utils.SafeIntPointer(limit, 0),
		Offset: utils.SafeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
		Select: shouldSelectItems(ctx),
	})
	return list, err
}

func (r *queryResolver) Server(ctx context.Context, key string) (*models.Server, error) {
	return r.ServerUcase.GetByKey(ctx, key)
}
