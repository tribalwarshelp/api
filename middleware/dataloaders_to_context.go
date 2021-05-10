package middleware

import (
	"context"
	"github.com/Kichiyaki/goutil/strutil"
	"github.com/tribalwarshelp/shared/tw/twmodel"
	"net/http"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/gin-gonic/gin"

	"github.com/tribalwarshelp/api/graphql/dataloader"
	"github.com/tribalwarshelp/api/server"
)

var serverDataLoadersContextKey ContextKey = "serverDataLoaders"
var versionDataLoadersContextKey ContextKey = "versionLoaders"
var dataLoaderContextKey ContextKey = "dataloader"

type DataLoadersToContextConfig struct {
	ServerRepo server.Repository
}

func DataLoadersToContext(dltcc DataLoadersToContextConfig, cfg dataloader.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		serverDataLoaders := make(map[string]*dataloader.ServerDataLoader)
		versionDataLoaders := make(map[twmodel.VersionCode]*dataloader.VersionDataLoader)
		servers, _, err := dltcc.ServerRepo.Fetch(c.Request.Context(), server.FetchConfig{
			Columns: []string{strutil.Underscore("versionCode"), "key"},
			Select:  true,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, &gqlerror.Error{
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		for _, serv := range servers {
			serverDataLoaders[serv.Key] = dataloader.NewServerDataLoader(serv.Key, cfg)
			if _, ok := versionDataLoaders[serv.VersionCode]; !ok {
				versionDataLoaders[serv.VersionCode] = dataloader.NewVersionDataLoader(serv.VersionCode, cfg)
			}
		}
		ctx = StoreServerDataLoadersInContext(ctx, serverDataLoaders)
		ctx = StoreVersionDataLoadersInContext(ctx, versionDataLoaders)
		ctx = StoreDataLoaderInContext(ctx, dataloader.NewDataLoader(cfg))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func StoreServerDataLoadersInContext(ctx context.Context, loaders map[string]*dataloader.ServerDataLoader) context.Context {
	return context.WithValue(ctx, serverDataLoadersContextKey, loaders)
}

func ServerDataLoadersFromContext(ctx context.Context) map[string]*dataloader.ServerDataLoader {
	dl := ctx.Value(serverDataLoadersContextKey)
	if dl == nil {
		return nil
	}

	return dl.(map[string]*dataloader.ServerDataLoader)
}

func StoreVersionDataLoadersInContext(ctx context.Context, loaders map[twmodel.VersionCode]*dataloader.VersionDataLoader) context.Context {
	return context.WithValue(ctx, versionDataLoadersContextKey, loaders)
}

func VersionDataLoadersFromContext(ctx context.Context) map[twmodel.VersionCode]*dataloader.VersionDataLoader {
	dl := ctx.Value(versionDataLoadersContextKey)
	if dl == nil {
		return nil
	}

	return dl.(map[twmodel.VersionCode]*dataloader.VersionDataLoader)
}

func StoreDataLoaderInContext(ctx context.Context, loaders *dataloader.DataLoader) context.Context {
	return context.WithValue(ctx, dataLoaderContextKey, loaders)
}

func DataLoaderFromContext(ctx context.Context) *dataloader.DataLoader {
	dl := ctx.Value(dataLoaderContextKey)
	if dl == nil {
		return nil
	}

	return dl.(*dataloader.DataLoader)
}
