package resolvers

import (
	"context"
	"github.com/tribalwarshelp/api/utils"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/serverstats"
	"github.com/tribalwarshelp/shared/models"
)

func (r *Resolver) ServerStats(ctx context.Context,
	server string,
	f *models.ServerStatsFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.ServerStats, error) {
	var err error
	list := &generated.ServerStats{}
	list.Items, list.Total, err = r.ServerStatsUcase.Fetch(ctx, serverstats.FetchConfig{
		Server: server,
		Filter: f,
		Sort:   sort,
		Limit:  utils.SafeIntPointer(limit, 0),
		Offset: utils.SafeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
		Select: shouldSelectItems(ctx),
	})
	return list, err
}
