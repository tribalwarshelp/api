package resolvers

import (
	"context"
	"time"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"
)

func (r *playerHistoryRecordResolver) Player(ctx context.Context, obj *models.PlayerHistory) (*models.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	return getPlayer(ctx, obj.PlayerID), nil
}

func (r *playerHistoryRecordResolver) Tribe(ctx context.Context, obj *models.PlayerHistory) (*models.Tribe, error) {
	if obj.Tribe != nil {
		return obj.Tribe, nil
	}

	return getTribe(ctx, obj.TribeID), nil
}

func (r *playerHistoryRecordResolver) CreatedAt(ctx context.Context, obj *models.PlayerHistory) (*time.Time, error) {
	server, _ := getServer(ctx)
	t := formatDate(ctx, utils.LanguageTagFromServerKey(server), obj.CreatedAt)
	return &t, nil
}

func (r *Resolver) PlayerHistory(ctx context.Context, server string, filter *models.PlayerHistoryFilter) (*generated.PlayerHistory, error) {
	var err error
	list := &generated.PlayerHistory{}
	list.Items, list.Total, err = r.PlayerHistoryUcase.Fetch(ctx, server, filter)
	return list, err
}
