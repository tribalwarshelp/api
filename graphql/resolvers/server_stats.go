package resolvers

import (
	"context"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *Resolver) ServerStats(ctx context.Context, server string, filter *models.ServerStatsFilter) (*generated.ServerStats, error) {
	var err error
	list := &generated.ServerStats{}
	list.Items, list.Total, err = r.ServerStatsUcase.Fetch(ctx, server, filter)
	return list, err
}
