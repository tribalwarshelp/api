package resolvers

import (
	"context"
	"time"

	"github.com/tribalwarshelp/api/utils"

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

func (r *serverResolver) DataUpdatedAt(ctx context.Context, obj *models.Server) (*time.Time, error) {
	loaders := middleware.DataLoadersFromContext(ctx)
	if loaders != nil {
		lv, err := loaders.LangVersionByTag.Load(obj.LangVersionTag.String())
		if err == nil {
			dataUpdatedAt := obj.DataUpdatedAt.In(utils.GetLocation(lv.Timezone))
			return &dataUpdatedAt, nil
		}
	}
	return &obj.DataUpdatedAt, nil
}

func (r *queryResolver) Servers(ctx context.Context, filter *models.ServerFilter) (*generated.ServersList, error) {
	var err error
	list := &generated.ServersList{}
	list.Items, list.Total, err = r.ServerUcase.Fetch(ctx, filter)
	return list, err
}

func (r *queryResolver) Server(ctx context.Context, key string) (*models.Server, error) {
	return r.ServerUcase.GetByKey(ctx, key)
}
