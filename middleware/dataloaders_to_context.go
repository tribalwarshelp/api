package middleware

import (
	"context"
	"net/http"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/tribalwarshelp/api/graphql/dataloaders"
	"github.com/tribalwarshelp/api/server"

	"github.com/gin-gonic/gin"
)

var serverDataLoadersContextKey ContextKey = "serverDataLoaders"
var dataloadersContextKey ContextKey = "dataloaders"

func DataLoadersToContext(serverRepo server.Repository, cfg dataloaders.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		loaders := make(map[string]*dataloaders.ServerDataLoaders)
		servers, _, err := serverRepo.Fetch(c.Request.Context(), nil)
		if err != nil {
			c.JSON(http.StatusOK, &gqlerror.Error{
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		for _, server := range servers {
			loaders[server.Key] = dataloaders.NewServerDataLoaders(server.Key, cfg)
		}
		ctx = StoreServerDataLoadersInContext(ctx, loaders)
		ctx = StoreDataLoadersInContext(ctx, dataloaders.New(cfg))
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
