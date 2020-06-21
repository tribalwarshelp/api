package resolvers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/shared/models"
)

func (r *liveEnnoblementResolver) NewOwner(ctx context.Context, obj *models.LiveEnnoblement) (*models.Player, error) {
	if obj.NewOwner != nil {
		return obj.NewOwner, nil
	}

	if server, ok := getServer(graphql.GetFieldContext(ctx)); ok {
		dataloaders := middleware.ServerDataLoadersFromContext(ctx)
		if dataloaders != nil {
			if dataloader, ok := dataloaders[server]; ok {
				player, _ := dataloader.PlayerByID.Load(obj.NewOwnerID)
				if player != nil {
					return player, nil
				}
			}
		}
	}

	return nil, nil
}

func (r *liveEnnoblementResolver) OldOwner(ctx context.Context, obj *models.LiveEnnoblement) (*models.Player, error) {
	if obj.OldOwner != nil {
		return obj.OldOwner, nil
	}

	if server, ok := getServer(graphql.GetFieldContext(ctx)); ok {
		dataloaders := middleware.ServerDataLoadersFromContext(ctx)
		if dataloaders != nil {
			if dataloader, ok := dataloaders[server]; ok {
				player, _ := dataloader.PlayerByID.Load(obj.OldOwnerID)
				if player != nil {
					return player, nil
				}
			}
		}
	}

	return nil, nil
}

func (r *liveEnnoblementResolver) Village(ctx context.Context, obj *models.LiveEnnoblement) (*models.Village, error) {
	if obj.Village != nil {
		return obj.Village, nil
	}

	if server, ok := getServer(graphql.GetFieldContext(ctx)); ok {
		dataloaders := middleware.ServerDataLoadersFromContext(ctx)
		if dataloaders != nil {
			if dataloader, ok := dataloaders[server]; ok {
				village, _ := dataloader.VillageByID.Load(obj.VillageID)
				if village != nil {
					return village, nil
				}
			}
		}
	}

	return nil, nil
}

func (r *queryResolver) LiveEnnoblements(ctx context.Context, server string) ([]*models.LiveEnnoblement, error) {
	return r.LiveEnnoblementUcase.Fetch(ctx, server)
}
