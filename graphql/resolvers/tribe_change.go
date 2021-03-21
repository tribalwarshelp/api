package resolvers

import (
	"context"
	"github.com/tribalwarshelp/api/utils"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/tribechange"
	"github.com/tribalwarshelp/shared/models"
)

func (r *tribeChangeRecordResolver) Player(ctx context.Context, obj *models.TribeChange) (*models.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	return getPlayer(ctx, obj.PlayerID), nil
}

func (r *tribeChangeRecordResolver) NewTribe(ctx context.Context, obj *models.TribeChange) (*models.Tribe, error) {
	if obj.NewTribe != nil {
		return obj.NewTribe, nil
	}

	return getTribe(ctx, obj.NewTribeID), nil
}

func (r *tribeChangeRecordResolver) OldTribe(ctx context.Context, obj *models.TribeChange) (*models.Tribe, error) {
	if obj.OldTribe != nil {
		return obj.OldTribe, nil
	}

	return getTribe(ctx, obj.OldTribeID), nil
}

func (r *Resolver) TribeChanges(ctx context.Context,
	server string,
	f *models.TribeChangeFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.TribeChanges, error) {
	var err error
	list := &generated.TribeChanges{}
	list.Items, list.Total, err = r.TribeChangeUcase.Fetch(ctx, tribechange.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  utils.SafeIntPointer(limit, 0),
		Offset: utils.SafeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
		Server: server,
	})
	return list, err
}
