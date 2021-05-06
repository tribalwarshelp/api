package resolvers

import (
	"context"
	"github.com/Kichiyaki/goutil/safeptr"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/tribalwarshelp/api/ennoblement"
	"github.com/tribalwarshelp/api/graphql/generated"
)

func (r *ennoblementResolver) NewOwner(ctx context.Context, obj *twmodel.Ennoblement) (*twmodel.Player, error) {
	if obj.NewOwner != nil {
		return obj.NewOwner, nil
	}

	return getPlayer(ctx, obj.NewOwnerID), nil
}

func (r *ennoblementResolver) NewOwnerTribe(ctx context.Context, obj *twmodel.Ennoblement) (*twmodel.Tribe, error) {
	if obj.NewOwnerTribe != nil {
		return obj.NewOwnerTribe, nil
	}

	return getTribe(ctx, obj.NewOwnerTribeID), nil
}

func (r *ennoblementResolver) OldOwner(ctx context.Context, obj *twmodel.Ennoblement) (*twmodel.Player, error) {
	if obj.OldOwner != nil {
		return obj.OldOwner, nil
	}

	return getPlayer(ctx, obj.OldOwnerID), nil
}

func (r *ennoblementResolver) OldOwnerTribe(ctx context.Context, obj *twmodel.Ennoblement) (*twmodel.Tribe, error) {
	if obj.OldOwnerTribe != nil {
		return obj.OldOwnerTribe, nil
	}

	return getTribe(ctx, obj.OldOwnerTribeID), nil
}

func (r *ennoblementResolver) Village(ctx context.Context, obj *twmodel.Ennoblement) (*twmodel.Village, error) {
	if obj.Village != nil {
		return obj.Village, nil
	}

	return getVillage(ctx, obj.VillageID), nil
}

func (r *queryResolver) Ennoblements(ctx context.Context, server string,
	f *twmodel.EnnoblementFilter,
	limit *int,
	offset *int,
	sort []string) (*generated.EnnoblementList, error) {
	var err error
	list := &generated.EnnoblementList{}
	list.Items, list.Total, err = r.EnnoblementUcase.Fetch(ctx, ennoblement.FetchConfig{
		Server: server,
		Filter: f,
		Sort:   sort,
		Limit:  safeptr.SafeIntPointer(limit, 0),
		Offset: safeptr.SafeIntPointer(offset, 0),
		Count:  shouldCount(ctx),
		Select: shouldSelectItems(ctx),
	})
	return list, err
}
