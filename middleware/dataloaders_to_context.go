package middleware

import (
	"context"
	"net/http"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/tribalwarshelp/api/graphql/dataloaders"
	"github.com/tribalwarshelp/api/server"
	"github.com/tribalwarshelp/api/utils"
	"github.com/tribalwarshelp/shared/models"

	"github.com/gin-gonic/gin"
)

var serverDataLoadersContextKey ContextKey = "serverDataLoaders"
var versionLoadersContextKey ContextKey = "versionLoaders"
var dataloadersContextKey ContextKey = "dataloaders"

type DataLoadersToContextConfig struct {
	ServerRepo server.Repository
}

func DataLoadersToContext(dltcc DataLoadersToContextConfig, cfg dataloaders.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		serverDataLoaders := make(map[string]*dataloaders.ServerDataLoaders)
		versionDataLoaders := make(map[models.VersionCode]*dataloaders.VersionDataLoaders)
		servers, _, err := dltcc.ServerRepo.Fetch(c.Request.Context(), server.FetchConfig{
			Columns: []string{utils.Underscore("versionCode"), "key"},
		})
		if err != nil {
			c.JSON(http.StatusOK, &gqlerror.Error{
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		for _, server := range servers {
			serverDataLoaders[server.Key] = dataloaders.NewServerDataLoaders(server.Key, cfg)
			if _, ok := versionDataLoaders[server.VersionCode]; !ok {
				versionDataLoaders[server.VersionCode] = dataloaders.NewVersionDataLoaders(server.VersionCode, cfg)
			}
		}
		ctx = StoreServerDataLoadersInContext(ctx, serverDataLoaders)
		ctx = StoreVersionDataLoadersInContext(ctx, versionDataLoaders)
		ctx = StoreDataLoadersInContext(ctx, dataloaders.NewDataLoaders(cfg))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func StoreServerDataLoadersInContext(ctx context.Context, loaders map[string]*dataloaders.ServerDataLoaders) context.Context {
	return context.WithValue(ctx, serverDataLoadersContextKey, loaders)
}

func ServerDataLoadersFromContext(ctx context.Context) map[string]*dataloaders.ServerDataLoaders {
	dl := ctx.Value(serverDataLoadersContextKey)
	if dl == nil {
		return nil
	}

	return dl.(map[string]*dataloaders.ServerDataLoaders)
}

func StoreVersionDataLoadersInContext(ctx context.Context, loaders map[models.VersionCode]*dataloaders.VersionDataLoaders) context.Context {
	return context.WithValue(ctx, versionLoadersContextKey, loaders)
}

func VersionDataLoadersFromContext(ctx context.Context) map[models.VersionCode]*dataloaders.VersionDataLoaders {
	dl := ctx.Value(versionLoadersContextKey)
	if dl == nil {
		return nil
	}

	return dl.(map[models.VersionCode]*dataloaders.VersionDataLoaders)
}

func StoreDataLoadersInContext(ctx context.Context, loaders *dataloaders.DataLoaders) context.Context {
	return context.WithValue(ctx, dataloadersContextKey, loaders)
}

func DataLoadersFromContext(ctx context.Context) *dataloaders.DataLoaders {
	dl := ctx.Value(dataloadersContextKey)
	if dl == nil {
		return nil
	}

	return dl.(*dataloaders.DataLoaders)
}
