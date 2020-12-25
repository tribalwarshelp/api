package resolvers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/tribalwarshelp/api/middleware"
	"github.com/tribalwarshelp/shared/models"
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

func getPlayer(ctx context.Context, id int) *models.Player {
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

func getVillage(ctx context.Context, id int) *models.Village {
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

func getTribe(ctx context.Context, id int) *models.Tribe {
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

func shouldCount(ctx context.Context) bool {
	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		if field.Name == countField {
			return true
		}
	}
	return false
}

func safeStrPointer(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

func safeIntPointer(s *int, def int) int {
	if s == nil {
		return def
	}
	return *s
}
