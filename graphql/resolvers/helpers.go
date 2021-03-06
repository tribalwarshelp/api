package resolvers

import (
	"context"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/99designs/gqlgen/graphql"

	"github.com/tribalwarshelp/api/middleware"
)

func getServer(ctx context.Context) (string, bool) {
	rctx := graphql.GetFieldContext(ctx)
	server := ""
	ok := false
	for rctx != nil {
		server, ok = rctx.Args["server"].(string)
		if ok {
			break
		}
		rctx = rctx.Parent
	}
	return server, ok
}

func getPlayer(ctx context.Context, id int) *twmodel.Player {
	if server, ok := getServer(ctx); ok {
		dataloaders := middleware.ServerDataLoadersFromContext(ctx)
		if dataloaders != nil {
			if dataloader, ok := dataloaders[server]; ok {
				player, _ := dataloader.PlayerByID.Load(id)
				if player != nil {
					return player
				}
			}
		}
	}
	return nil
}

func getVillage(ctx context.Context, id int) *twmodel.Village {
	if server, ok := getServer(ctx); ok {
		dataloaders := middleware.ServerDataLoadersFromContext(ctx)
		if dataloaders != nil {
			if dataloader, ok := dataloaders[server]; ok {
				player, _ := dataloader.VillageByID.Load(id)
				if player != nil {
					return player
				}
			}
		}
	}
	return nil
}

func getTribe(ctx context.Context, id int) *twmodel.Tribe {
	if server, ok := getServer(ctx); ok {
		dataloaders := middleware.ServerDataLoadersFromContext(ctx)
		if dataloaders != nil {
			if dataloader, ok := dataloaders[server]; ok {
				player, _ := dataloader.TribeByID.Load(id)
				if player != nil {
					return player
				}
			}
		}
	}
	return nil
}

func findField(ctx context.Context, name string) bool {
	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		if field.Name == name {
			return true
		}
	}
	return false
}

func shouldCount(ctx context.Context) bool {
	return findField(ctx, countField)
}

func shouldSelectItems(ctx context.Context) bool {
	return findField(ctx, itemsField)
}
