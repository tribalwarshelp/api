package resolvers

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/shared/models"
)

func (r *ennoblementResolver) NewOwner(ctx context.Context, obj *models.Ennoblement) (*models.Player, error) {
	if obj.NewOwner != nil {
		return obj.NewOwner, nil
	}

	rctx := graphql.GetFieldContext(ctx)
	server, ok := rctx.Parent.Parent.Args["server"].(string)
	if ok {
		dataloaders := middleware.DataLoadersFromContext(ctx)
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

func (r *ennoblementResolver) OldOwner(ctx context.Context, obj *models.Ennoblement) (*models.Player, error) {
	if obj.OldOwner != nil {
		return obj.OldOwner, nil
	}

	rctx := graphql.GetFieldContext(ctx)
	server, ok := rctx.Parent.Parent.Args["server"].(string)
	if ok {
		dataloaders := middleware.DataLoadersFromContext(ctx)
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

func (r *ennoblementResolver) Village(ctx context.Context, obj *models.Ennoblement) (*models.Village, error) {
	if obj.Village != nil {
		return obj.Village, nil
	}

	rctx := graphql.GetFieldContext(ctx)
	server, ok := rctx.Parent.Parent.Args["server"].(string)
	if ok {
		dataloaders := middleware.DataLoadersFromContext(ctx)
		if dataloaders != nil {
			if dataloader, ok := dataloaders[server]; ok {
				log.Print("?")
				village, _ := dataloader.VillageByID.Load(obj.VillageID)
				if village != nil {
					return village, nil
				}
			}
		}
	}

	return nil, nil
}

func (r *queryResolver) Ennoblements(ctx context.Context, server string) ([]*models.Ennoblement, error) {
	return r.EnnoblementUcase.Fetch(ctx, server)
}
