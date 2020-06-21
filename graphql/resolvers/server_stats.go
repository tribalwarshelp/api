package resolvers

import (
	"context"
	"time"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

func (r *serverStatsRecordResolver) CreatedAt(ctx context.Context, obj *models.ServerStats) (*time.Time, error) {
	server, _ := getServer(ctx)
	t := formatDate(ctx, utils.LanguageTagFromServerKey(server), obj.CreatedAt)
	return &t, nil
}

func (r *Resolver) ServerStats(ctx context.Context, server string, filter *models.ServerStatsFilter) (*generated.ServerStatsList, error) {
	var err error
	list := &generated.ServerStatsList{}
	list.Items, list.Total, err = r.ServerStatsUcase.Fetch(ctx, server, filter)
	return list, err
}
