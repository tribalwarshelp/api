package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/api/server"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *serverResolver) Version(ctx context.Context, obj *models.Server) (*models.Version, error) {
	loaders := middleware.DataLoadersFromContext(ctx)
	if loaders != nil {
		lv, _ := loaders.VersionByTag.Load(obj.VersionCode.String())
		return lv, nil
	}
	return nil, nil
}

func (r *serverResolver) LangVersion(ctx context.Context, obj *models.Server) (*models.Version, error) {
	return r.Version(ctx, obj)
}

func (r *queryResolver) Servers(ctx context.Context,
	f *models.ServerFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.ServerList, error) {
	defLimit := 0
	defOffset := 0
	if limit == nil {
		limit = &defLimit
	}
	if offset == nil {
		offset = &defOffset
	}

	var err error
	list := &generated.ServerList{}
	list.Items, list.Total, err = r.ServerUcase.Fetch(ctx, server.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  *limit,
		Offset: *offset,
		Count:  true,
	})
	return list, err
}

func (r *queryResolver) Server(ctx context.Context, key string) (*models.Server, error) {
	return r.ServerUcase.GetByKey(ctx, key)
}
