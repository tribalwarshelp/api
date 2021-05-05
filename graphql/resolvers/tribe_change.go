package resolvers

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/api/tribechange"
)

func (r *tribeChangeRecordResolver) Player(ctx context.Context, obj *twmodel.TribeChange) (*twmodel.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}

	return getPlayer(ctx, obj.PlayerID), nil
}

func (r *tribeChangeRecordResolver) NewTribe(ctx context.Context, obj *twmodel.TribeChange) (*twmodel.Tribe, error) {
	if obj.NewTribe != nil {
		return obj.NewTribe, nil
	}

	return getTribe(ctx, obj.NewTribeID), nil
}

func (r *tribeChangeRecordResolver) OldTribe(ctx context.Context, obj *twmodel.TribeChange) (*twmodel.Tribe, error) {
	if obj.OldTribe != nil {
		return obj.OldTribe, nil
	}

	return getTribe(ctx, obj.OldTribeID), nil
}

func (r *Resolver) TribeChanges(ctx context.Context,
	server string,
	f *twmodel.TribeChangeFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.TribeChanges, error) {
	var err error
	list := &generated.TribeChanges{}
	list.Items, list.Total, err = r.TribeChangeUcase.Fetch(ctx, tribechange.FetchConfig{
		Filter: f,
		Sort:   sort,
		Limit:  safeptr.SafeIntPointer(limit, 0),
		Offset: safeptr.SafeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
		Select: shouldSelectItems(ctx),
		Server: server,
	})
	return list, err
}
