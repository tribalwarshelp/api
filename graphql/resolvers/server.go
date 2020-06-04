package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *queryResolver) Servers(ctx context.Context, filter *models.ServerFilter) (*generated.ServersList, error) {
	var err error
	list := &generated.ServersList{}
	list.Items, list.Total, err = r.ServerUcase.Fetch(ctx, filter)
	return list, err
}

func (r *queryResolver) Server(ctx context.Context, key string) (*models.Server, error) {
	return r.ServerUcase.GetByKey(ctx, key)
}
