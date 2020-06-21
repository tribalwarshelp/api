package resolvers

import (
	"context"
	"time"

	"github.com/tribalwarshelp/api/utils"

	"github.com/tribalwarshelp/api/graphql/generated"
	"github.com/tribalwarshelp/shared/models"
)

func (r *ennoblementResolver) NewOwner(ctx context.Context, obj *models.Ennoblement) (*models.Player, error) {
	if obj.NewOwner != nil {
		return obj.NewOwner, nil
	}

	return getPlayer(ctx, obj.NewOwnerID), nil
}

func (r *ennoblementResolver) NewOwnerTribe(ctx context.Context, obj *models.Ennoblement) (*models.Tribe, error) {
	if obj.NewOwnerTribe != nil {
		return obj.NewOwnerTribe, nil
	}

	return getTribe(ctx, obj.NewOwnerTribeID), nil
}

func (r *ennoblementResolver) OldOwner(ctx context.Context, obj *models.Ennoblement) (*models.Player, error) {
	if obj.OldOwner != nil {
		return obj.OldOwner, nil
	}

	return getPlayer(ctx, obj.OldOwnerID), nil
}

func (r *ennoblementResolver) OldOwnerTribe(ctx context.Context, obj *models.Ennoblement) (*models.Tribe, error) {
	if obj.OldOwnerTribe != nil {
		return obj.OldOwnerTribe, nil
	}

	return getTribe(ctx, obj.OldOwnerTribeID), nil
}

func (r *ennoblementResolver) Village(ctx context.Context, obj *models.Ennoblement) (*models.Village, error) {
	if obj.Village != nil {
		return obj.Village, nil
	}

	return getVillage(ctx, obj.VillageID), nil
}

func (r *ennoblementResolver) EnnobledAt(ctx context.Context, obj *models.Ennoblement) (*time.Time, error) {
	server, _ := getServer(ctx)
	t := formatDate(ctx, utils.LanguageTagFromServerKey(server), obj.EnnobledAt)
	return &t, nil
}

func (r *queryResolver) Ennoblements(ctx context.Context, server string, f *models.EnnoblementFilter) (*generated.EnnoblementsList, error) {
	var err error
	list := &generated.EnnoblementsList{}
	list.Items, list.Total, err = r.EnnoblementUcase.Fetch(ctx, server, f)
	return list, err
}