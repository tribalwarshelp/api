package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/middleware"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *serverResolver) LangVersion(ctx context.Context, obj *models.Server) (*models.LangVersion, error) {
	loaders := middleware.DataLoadersFromContext(ctx)
	if loaders != nil {
		lv, _ := loaders.LangVersionByTag.Load(obj.LangVersionTag.String())
		return lv, nil
	}
	return nil, nil
}

func (r *queryResolver) Servers(ctx context.Context, filter *models.ServerFilter) (*generated.ServerList, error) {
	var err error
	list := &generated.ServerList{}
	list.Items, list.Total, err = r.ServerUcase.Fetch(ctx, filter)
	return list, err
}

func (r *queryResolver) Server(ctx context.Context, key string) (*models.Server, error) {
	return r.ServerUcase.GetByKey(ctx, key)
}
